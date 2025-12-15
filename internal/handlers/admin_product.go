package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"backend/internal/models"
)

// ProductCreateRequest represents the payload for creating a product.
type ProductCreateRequest struct {
	Name     string   `json:"name" binding:"required"`
	Price    float64  `json:"price" binding:"required"`
	Category []string `json:"category" binding:"required"`
	ImageURL string   `json:"imageUrl" binding:"required"`
	IsActive *bool    `json:"isActive"`
}

// ProductUpdateRequest represents the payload for updating a product.
type ProductUpdateRequest struct {
	Name     *string   `json:"name"`
	Price    *float64  `json:"price"`
	Category *[]string `json:"category"`
	ImageURL *string   `json:"imageUrl"`
	IsActive *bool     `json:"isActive"`
}

func normalizeCategories(values []string) []string {
	cleaned := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))

	for _, v := range values {
		name := strings.TrimSpace(v)
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		cleaned = append(cleaned, name)
	}

	return cleaned
}

// GetAllProducts returns all products for admin users.
func GetAllProducts(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, limit, err := parsePaginationParams(c.Query("page"), c.Query("limit"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{}

		if category := strings.TrimSpace(c.Query("category")); category != "" {
			filter["category"] = bson.M{"$in": []string{category}}
		}

		if search := strings.TrimSpace(c.Query("search")); search != "" {
			filter["name"] = bson.M{"$regex": search, "$options": "i"}
		}

		if isActive := strings.TrimSpace(c.Query("isActive")); isActive != "" {
			active, parseErr := strconv.ParseBool(isActive)
			if parseErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid isActive parameter"})
				return
			}
			filter["isActive"] = active
		}

		findOptions := options.Find().
			SetSkip((page - 1) * limit).
			SetLimit(limit).
			SetSort(bson.D{{Key: "createdAt", Value: -1}})

		total, err := db.Collection("products").CountDocuments(context.Background(), filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		cursor, err := db.Collection("products").Find(context.Background(), filter, findOptions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		defer cursor.Close(context.Background())

		var products []models.Product
		if err := cursor.All(context.Background(), &products); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "decode error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": products,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		})
	}
}

// CreateProduct handles creation of a new product.
func CreateProduct(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ProductCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		categories := normalizeCategories(req.Category)
		if len(categories) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "category required"})
			return
		}

		isActive := true
		if req.IsActive != nil {
			isActive = *req.IsActive
		}

		product := models.Product{
			Name:      req.Name,
			Price:     req.Price,
			Category:  categories,
			ImageURL:  req.ImageURL,
			IsActive:  isActive,
			CreatedAt: time.Now(),
		}

		result, err := db.Collection("products").InsertOne(context.Background(), product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		product.ID = result.InsertedID.(primitive.ObjectID)

		c.JSON(http.StatusCreated, product)
	}
}

// UpdateProduct updates fields of an existing product.
func UpdateProduct(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var req ProductUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		update := bson.M{}
		if req.Name != nil {
			update["name"] = *req.Name
		}
		if req.Price != nil {
			update["price"] = *req.Price
		}
		if req.Category != nil {
			categories := normalizeCategories(*req.Category)
			if len(categories) == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "category required"})
				return
			}
			update["category"] = categories
		}
		if req.ImageURL != nil {
			update["imageUrl"] = *req.ImageURL
		}
		if req.IsActive != nil {
			update["isActive"] = *req.IsActive
		}

		if len(update) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
			return
		}

		var updated models.Product
		err = db.Collection("products").FindOneAndUpdate(
			context.Background(),
			bson.M{"_id": id},
			bson.M{"$set": update},
			options.FindOneAndUpdate().SetReturnDocument(options.After),
		).Decode(&updated)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		c.JSON(http.StatusOK, updated)
	}
}

// DeleteProduct marks a product as inactive.
func DeleteProduct(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		result, err := db.Collection("products").UpdateOne(
			context.Background(),
			bson.M{"_id": id},
			bson.M{"$set": bson.M{"isActive": false}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
