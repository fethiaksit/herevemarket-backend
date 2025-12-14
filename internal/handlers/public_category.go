package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"backend/internal/models"
)

func GetCategories(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		cursor, err := db.Collection("categories").Find(
			context.Background(),
			bson.M{"isActive": true},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		defer cursor.Close(context.Background())

		var categories []models.Category
		if err := cursor.All(context.Background(), &categories); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "decode error"})
			return
		}

		c.JSON(http.StatusOK, categories)
	}
}
