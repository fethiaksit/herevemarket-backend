package handlers

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"backend/internal/models"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type LoginResponseUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AuthTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

func Register(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		name := strings.TrimSpace(req.Name)
		email := strings.ToLower(strings.TrimSpace(req.Email))

		if name == "" || email == "" || strings.TrimSpace(req.Password) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name, email and password are required"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		count, err := db.Collection("users").CountDocuments(ctx, bson.M{"email": email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "password hash failed"})
			return
		}

		now := time.Now()
		user := models.User{
			Name:         name,
			Email:        email,
			PasswordHash: string(hash),
			Role:         "user",
			IsActive:     true,
			CreatedAt:    now,
		}

		if _, err := db.Collection("users").InsertOne(ctx, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
	}
}

func Login(db *mongo.Database, jwtSecret string, accessTTL, refreshTTL time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		email := strings.ToLower(strings.TrimSpace(req.Email))
		if email == "" || strings.TrimSpace(req.Password) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var user models.User
		if err := db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		if !user.IsActive {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		tokens, err := issueTokens(c, db, user.ID, user.Email, user.Role, jwtSecret, accessTTL, refreshTTL)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"accessToken":  tokens.AccessToken,
			"refreshToken": tokens.RefreshToken,
			"expiresIn":    tokens.ExpiresIn,
			"user": LoginResponseUser{
				ID:    user.ID.Hex(),
				Name:  user.Name,
				Email: user.Email,
			},
		})
	}
}

func Refresh(db *mongo.Database, jwtSecret string, accessTTL, refreshTTL time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RefreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		plain := strings.TrimSpace(req.RefreshToken)
		if plain == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refreshToken is required"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		hash := hashToken(plain)
		var token models.RefreshToken
		if err := db.Collection("refresh_tokens").FindOne(ctx, bson.M{
			"tokenHash": hash,
			"revoked":   false,
		}).Decode(&token); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}

		if time.Now().After(token.ExpiresAt) {
			_, _ = db.Collection("refresh_tokens").UpdateByID(ctx, token.ID, bson.M{"$set": bson.M{"revoked": true}})
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token_expired"})
			return
		}

		var user models.User
		if err := db.Collection("users").FindOne(ctx, bson.M{"_id": token.UserID}).Decode(&user); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		if !user.IsActive {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		newTokens, err := issueTokens(c, db, user.ID, user.Email, user.Role, jwtSecret, accessTTL, refreshTTL)
		if err != nil {
			return
		}

		_, _ = db.Collection("refresh_tokens").UpdateByID(ctx, token.ID, bson.M{
			"$set": bson.M{
				"revoked":         true,
				"replacedByToken": newTokens.RefreshTokenID,
			},
		})

		c.JSON(http.StatusOK, gin.H{
			"accessToken":  newTokens.AccessToken,
			"refreshToken": newTokens.RefreshToken,
			"expiresIn":    newTokens.ExpiresIn,
			"user": LoginResponseUser{
				ID:    user.ID.Hex(),
				Name:  user.Name,
				Email: user.Email,
			},
		})
	}
}

func Logout(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RefreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		plain := strings.TrimSpace(req.RefreshToken)
		if plain == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refreshToken is required"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		hash := hashToken(plain)
		res, err := db.Collection("refresh_tokens").UpdateOne(ctx, bson.M{
			"tokenHash": hash,
			"revoked":   false,
		}, bson.M{"$set": bson.M{"revoked": true}})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		if res.MatchedCount == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "logged out"})
	}
}

type issuedTokens struct {
	AccessToken    string
	RefreshToken   string
	RefreshTokenID primitive.ObjectID
	ExpiresIn      int64
}

func issueTokens(c *gin.Context, db *mongo.Database, userID primitive.ObjectID, email, role, secret string, accessTTL, refreshTTL time.Duration) (*issuedTokens, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   userID.Hex(),
		"role":  role,
		"email": email,
		"exp":   now.Add(accessTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return nil, err
	}

	plainRefresh := generateRefreshString()
	if plainRefresh == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return nil, errors.New("could not generate refresh token")
	}
	hashed := hashToken(plainRefresh)

	refresh := models.RefreshToken{
		UserID:    userID,
		TokenHash: hashed,
		ExpiresAt: now.Add(refreshTTL),
		Revoked:   false,
		CreatedAt: now,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := db.Collection("refresh_tokens").InsertOne(ctx, refresh)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return nil, err
	}

	refreshID := res.InsertedID.(primitive.ObjectID)
	return &issuedTokens{
		AccessToken:    accessToken,
		RefreshToken:   plainRefresh,
		RefreshTokenID: refreshID,
		ExpiresIn:      int64(accessTTL.Seconds()),
	}, nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func generateRefreshString() string {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return ""
	}
	return hex.EncodeToString(buf)
}
