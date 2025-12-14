package handlers

import (
	"context"
	"net/http"
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
	Name     string  `json:"name" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
	Category string  `json:"category" binding:"required"`
	ImageURL string  `json:"imageUrl" binding:"required"`
	IsActive *bool   `json:"isActive"`
}

// ProductUpdateRequest represents the payload for updating a product.
type ProductUpdateRequest struct {
	Name     *string  `json:"name"`
	Price    *float64 `json:"price"`
	Category *string  `json:"category"`
	ImageURL *string  `json:"imageUrl"`
	IsActive *bool    `json:"isActive"`
}

// GetAllProducts returns all products for admin users.
func GetAllProducts(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		cursor, err := db.Collection("products").Find(context.Background(), bson.M{})
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

		c.JSON(http.StatusOK, products)
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

		isActive := true
		if req.IsActive != nil {
			isActive = *req.IsActive
		}

		product := models.Product{
			Name:      req.Name,
			Price:     req.Price,
			Category:  req.Category,
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
			update["category"] = *req.Category
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
