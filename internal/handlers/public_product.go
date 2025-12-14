package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"backend/internal/models"
)

func GetProducts(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		cursor, err := db.Collection("products").Find(
			context.Background(),
			bson.M{"isActive": true},
		)
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
