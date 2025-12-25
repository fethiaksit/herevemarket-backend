package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
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
	Name        string   `json:"name" binding:"required"`
	Price       float64  `json:"price" binding:"required"`
	Category    []string `json:"category" binding:"required"`
	ImageURL    string   `json:"imageUrl" binding:"required"`
	Description string   `json:"description"`
	Barcode     string   `json:"barcode"`
	Brand       string   `json:"brand"`
	Stock       *int     `json:"stock"`
	IsActive    *bool    `json:"isActive"`
	IsCampaign  *bool    `json:"isCampaign"`
}

type ProductUpdateRequest struct {
	Name        *string   `json:"name"`
	Price       *float64  `json:"price"`
	Category    *[]string `json:"category"`
	ImageURL    *string   `json:"imageUrl"`
	Description *string   `json:"description"`
	Barcode     *string   `json:"barcode"`
	Brand       *string   `json:"brand"`
	Stock       *int      `json:"stock"`
	IsActive    *bool     `json:"isActive"`
	IsCampaign  *bool     `json:"isCampaign"`
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

		filter := bson.M{
			"isDeleted": bson.M{"$ne": true},
		}

		if category := strings.TrimSpace(c.Query("category")); category != "" {
			filter["category"] = bson.M{"$in": []string{category}}
		}

		if search := strings.TrimSpace(c.Query("search")); search != "" {
			filter["name"] = bson.M{"$regex": search, "$options": "i"}
		}

		if isActive := strings.TrimSpace(c.Query("isActive")); isActive != "" {
			filter["isActive"] = strings.EqualFold(isActive, "true")
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
		log.Println("CreateProduct: request received")
		body, err := c.GetRawData()
		if err != nil {
			log.Println("CreateProduct RETURN 400:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		log.Println("CreateProduct raw body:", string(body))

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		var req ProductCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Println("CreateProduct bind error:", err)
			log.Println("CreateProduct RETURN 400:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("CreateProduct parsed request: %+v", req)

		if strings.TrimSpace(req.Name) == "" {
			log.Println("CreateProduct RETURN 400:", "name required")
			c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
			return
		}

		if req.Price <= 0 {
			log.Println("CreateProduct RETURN 400:", "invalid price")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
			return
		}

		categories := normalizeCategories(req.Category)
		if len(categories) == 0 {
			log.Println("CreateProduct RETURN 400:", "category required")
			c.JSON(http.StatusBadRequest, gin.H{"error": "category required"})
			return
		}

		if req.Stock == nil {
			log.Println("CreateProduct RETURN 400:", "stock required")
			c.JSON(http.StatusBadRequest, gin.H{"error": "stock required"})
			return
		}

		if *req.Stock < 0 {
			log.Println("CreateProduct RETURN 400:", "stock must be zero or greater")
			c.JSON(http.StatusBadRequest, gin.H{"error": "stock must be zero or greater"})
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

		barcode := strings.TrimSpace(req.Barcode)
		brand := strings.TrimSpace(req.Brand)
		description := strings.TrimSpace(req.Description)

		product := models.Product{
			Name:        req.Name,
			Price:       req.Price,
			Category:    models.StringList(categories),
			ImageURL:    req.ImageURL,
			Description: description,
			Barcode:     barcode,
			Brand:       brand,
			Stock:       *req.Stock,
			InStock:     *req.Stock > 0,
			IsActive:    isActive,
			IsCampaign:  isCampaign,
			IsDeleted:   false,
			CreatedAt:   now,
		}

		log.Printf("CreateProduct inserting product: %+v", product)
		res, err := db.Collection("products").InsertOne(context.Background(), product)
		if err != nil {
			log.Println("CreateProduct insert error:", err)
			if mongo.IsDuplicateKeyError(err) {
				log.Println("CreateProduct RETURN 409:", err)
				c.JSON(http.StatusConflict, gin.H{"error": "barcode already exists"})
				return
			}
			log.Println("CreateProduct RETURN 500:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		product.ID = res.InsertedID.(primitive.ObjectID)
		log.Println("CreateProduct insert success:", res.InsertedID)
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
			log.Println("UpdateProduct RETURN 400:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		log.Println("UpdateProduct request received for id:", id.Hex())

		body, err := c.GetRawData()
		if err != nil {
			log.Println("UpdateProduct read body error:", err)
			log.Println("UpdateProduct RETURN 400:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		log.Println("UpdateProduct raw body:", string(body))

		var raw map[string]interface{}
		if err := json.Unmarshal(body, &raw); err != nil {
			log.Println("UpdateProduct raw json error:", err)
			log.Println("UpdateProduct RETURN 400:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		if val, ok := raw["isCampaign"]; ok {
			if _, ok := val.(bool); !ok {
				log.Println("UpdateProduct RETURN 400:", "isCampaign must be boolean")
				c.JSON(http.StatusBadRequest, gin.H{"error": "isCampaign must be boolean"})
				return
			}
		}

		var req ProductUpdateRequest
		if err := json.Unmarshal(body, &req); err != nil {
			log.Println("UpdateProduct bind error:", err)
			log.Println("UpdateProduct RETURN 400:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		log.Printf("UpdateProduct parsed request: %+v", req)

		updateSet := bson.M{}
		updateUnset := bson.M{}

		if req.Name != nil {
			updateSet["name"] = *req.Name
		}
		if req.Price != nil {
			if *req.Price <= 0 {
				log.Println("UpdateProduct RETURN 400:", "invalid price")
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
				return
			}
			updateSet["price"] = *req.Price
		}
		if req.Category != nil {
			cats := normalizeCategories(*req.Category)
			if len(cats) == 0 {
				log.Println("UpdateProduct RETURN 400:", "category required")
				c.JSON(http.StatusBadRequest, gin.H{"error": "category required"})
				return
			}

			updateSet["category"] = models.StringList(cats)

		}
		if req.ImageURL != nil {
			updateSet["imageUrl"] = *req.ImageURL
		}
		if req.Description != nil {
			updateSet["description"] = strings.TrimSpace(*req.Description)
		}
		if req.Barcode != nil {
			barcode := strings.TrimSpace(*req.Barcode)
			if barcode == "" {
				updateUnset["barcode"] = ""
			} else {
				updateSet["barcode"] = barcode
			}
		}
		if req.Brand != nil {
			updateSet["brand"] = strings.TrimSpace(*req.Brand)
		}
		if req.Stock != nil {
			if *req.Stock < 0 {
				log.Println("UpdateProduct RETURN 400:", "stock must be zero or greater")
				c.JSON(http.StatusBadRequest, gin.H{"error": "stock must be zero or greater"})
				return
			}
			updateSet["stock"] = *req.Stock
		}
		if req.IsActive != nil {
			updateSet["isActive"] = *req.IsActive
		}
		if req.IsCampaign != nil {
			updateSet["isCampaign"] = *req.IsCampaign
		}

		if len(updateSet) == 0 && len(updateUnset) == 0 {
			log.Println("UpdateProduct RETURN 400:", "no fields to update")
			c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
			return
		}

		update := bson.M{}
		if len(updateSet) > 0 {
			update["$set"] = updateSet
		}
		if len(updateUnset) > 0 {
			update["$unset"] = updateUnset
		}
		log.Printf("UpdateProduct update document: %+v", update)

		result, err := db.Collection("products").UpdateOne(
			context.Background(),
			bson.M{
				"_id":       id,
				"isDeleted": bson.M{"$ne": true},
			},
			update,
		)

		if err != nil {
			log.Println("UpdateProduct update error:", err)
			if mongo.IsDuplicateKeyError(err) {
				log.Println("UpdateProduct RETURN 409:", err)
				c.JSON(http.StatusConflict, gin.H{"error": "barcode already exists"})
				return
			}
			log.Println("UpdateProduct RETURN 500:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		log.Printf("UpdateProduct update result: matched=%d modified=%d", result.MatchedCount, result.ModifiedCount)

		if result.MatchedCount == 0 {
			log.Println("UpdateProduct RETURN 404:", "product not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}

		var updated models.Product
		err = db.Collection("products").FindOne(
			context.Background(),
			bson.M{
				"_id":       id,
				"isDeleted": bson.M{"$ne": true},
			},
		).Decode(&updated)

		if err == mongo.ErrNoDocuments {
			log.Println("UpdateProduct RETURN 404:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		if err != nil {
			log.Println("UpdateProduct find error:", err)
			log.Println("UpdateProduct RETURN 500:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		updated.InStock = updated.Stock > 0
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

		now := time.Now()

		res, err := db.Collection("products").UpdateOne(
			context.Background(),
			bson.M{
				"_id":       id,
				"isDeleted": bson.M{"$ne": true},
			},
			bson.M{"$set": bson.M{
				"isDeleted": true,
				"deletedAt": now,
				"isActive":  false,
			}},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		if res.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
	}
}
