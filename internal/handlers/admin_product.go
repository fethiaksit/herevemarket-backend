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
	Name       string   `json:"name" binding:"required"`
	Price      float64  `json:"price" binding:"required"`
	Category   []string `json:"category" binding:"required"`
	ImageURL   string   `json:"imageUrl" binding:"required"`
	IsActive   *bool    `json:"isActive"`
	IsCampaign *bool    `json:"isCampaign"`
}

type ProductUpdateRequest struct {
	Name       *string   `json:"name"`
	Price      *float64  `json:"price"`
	Category   *[]string `json:"category"`
	ImageURL   *string   `json:"imageUrl"`
	IsActive   *bool     `json:"isActive"`
	IsCampaign *bool     `json:"isCampaign"`
}

type ProductBulkDeleteRequest struct {
	IDs []string `json:"ids" binding:"required"`
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

		includeDeleted := strings.EqualFold(strings.TrimSpace(c.Query("includeDeleted")), "true")
		if !includeDeleted {
			filter["isDeleted"] = bson.M{"$ne": true}
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

		products, err := decodeProducts(ctx, cursor)
		if err != nil {
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

		isCampaign := false
		if req.IsCampaign != nil {
			isCampaign = *req.IsCampaign
		}

		now := time.Now()
		product := models.Product{
			Name:       req.Name,
			Price:      req.Price,
			Category:   models.StringList(categories),
			ImageURL:   req.ImageURL,
			IsActive:   isActive,
			IsCampaign: isCampaign,
			CreatedAt:  now,
			UpdatedAt:  now,
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
		if req.IsCampaign != nil {
			update["isCampaign"] = *req.IsCampaign
		}

		if len(update) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
			return
		}

		var existing models.Product
		if err := db.Collection("products").FindOne(context.Background(), bson.M{"_id": id}).Decode(&existing); err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		if existing.IsDeleted {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product is deleted"})
			return
		}

		update["updatedAt"] = time.Now()

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

		ctx := context.Background()
		var product models.Product
		if err := db.Collection("products").FindOne(ctx, bson.M{"_id": id}).Decode(&product); err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		if product.IsDeleted {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product already deleted"})
			return
		}

		now := time.Now()
		update := bson.M{
			"isDeleted":  true,
			"deletedAt":  now,
			"updatedAt":  now,
			"isActive":   false,
			"isCampaign": false,
		}

		if _, err := db.Collection("products").UpdateByID(ctx, id, bson.M{"$set": update}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	}
}

/* =======================
   BULK DELETE (SOFT)
======================= */

func BulkDeleteProducts(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ProductBulkDeleteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		if len(req.IDs) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ids required"})
			return
		}

		objectIDs := make([]primitive.ObjectID, 0, len(req.IDs))
		for _, raw := range req.IDs {
			id, err := primitive.ObjectIDFromHex(strings.TrimSpace(raw))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id: " + raw})
				return
			}
			objectIDs = append(objectIDs, id)
		}

		now := time.Now()
		filter := bson.M{
			"_id":       bson.M{"$in": objectIDs},
			"isDeleted": bson.M{"$ne": true},
		}
		update := bson.M{
			"$set": bson.M{
				"isDeleted":  true,
				"deletedAt":  now,
				"updatedAt":  now,
				"isActive":   false,
				"isCampaign": false,
			},
		}

		res, err := db.Collection("products").UpdateMany(context.Background(), filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		if res.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "products not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":        "Products deleted successfully",
			"matchedCount":   res.MatchedCount,
			"modifiedCount":  res.ModifiedCount,
			"alreadyDeleted": res.MatchedCount - res.ModifiedCount,
		})
	}
}
