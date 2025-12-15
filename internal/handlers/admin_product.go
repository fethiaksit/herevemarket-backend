package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"backend/internal/models"
)

/* =======================
   REQUEST MODELLERİ
======================= */

type ProductCreateRequest struct {
	Name     string   `json:"name" binding:"required"`
	Price    float64  `json:"price" binding:"required"`
	Category []string `json:"category" binding:"required"`
	ImageURL string   `json:"imageUrl" binding:"required"`
	IsActive *bool    `json:"isActive"`
}

type ProductUpdateRequest struct {
	Name     *string   `json:"name"`
	Price    *float64  `json:"price"`
	Category *[]string `json:"category"`
	ImageURL *string   `json:"imageUrl"`
	IsActive *bool     `json:"isActive"`
}

/* =======================
   HELPERS
======================= */

func normalizeCategories(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0)

	for _, v := range values {
		name := strings.TrimSpace(v)
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
	}
	return out
}

/* =======================
   GET (ADMIN) – LIST
======================= */

func GetAllProducts(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, limit, err := parsePaginationParams(
			c.Query("page"),
			c.Query("limit"),
		)
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
			filter["isActive"] = isActive == "true"
		}

		ctx := context.Background()

		total, err := db.Collection("products").CountDocuments(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		opts := options.Find().
			SetSkip((page - 1) * limit).
			SetLimit(limit).
			SetSort(bson.D{{Key: "createdAt", Value: -1}})

		cursor, err := db.Collection("products").Find(ctx, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		defer cursor.Close(ctx)

		products := make([]models.Product, 0)

		for cursor.Next(ctx) {
			var p models.Product
			if err := cursor.Decode(&p); err != nil {
				// ⛑ eski string category kayıtları için fallback
				var raw bson.M
				_ = cursor.Decode(&raw)

				if cat, ok := raw["category"].(string); ok {
					raw["category"] = []string{cat}
					bytes, _ := bson.Marshal(raw)
					_ = bson.Unmarshal(bytes, &p)
					products = append(products, p)
					continue
				}

				c.JSON(http.StatusInternalServerError, gin.H{"error": "decode error"})
				return
			}
			products = append(products, p)
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

/* =======================
   CREATE
======================= */

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

		res, err := db.Collection("products").InsertOne(context.Background(), product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		product.ID = res.InsertedID.(primitive.ObjectID)
		c.JSON(http.StatusCreated, product)
	}
}

/* =======================
   UPDATE
======================= */

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
			cats := normalizeCategories(*req.Category)
			if len(cats) == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "category required"})
				return
			}
			update["category"] = cats
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

/* =======================
   DELETE (SOFT)
======================= */

func DeleteProduct(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		res, err := db.Collection("products").UpdateOne(
			context.Background(),
			bson.M{"_id": id},
			bson.M{"$set": bson.M{"isActive": false}},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		if res.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
