package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetProducts(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		const route = "GET /products"
		defer handlePanic(c, route)

		log.Printf("[%s] hit with query: page=%s limit=%s category=%s search=%s", route, c.Query("page"), c.Query("limit"), c.Query("category"), c.Query("search"))

		if err := ensureDBConnection(c.Request.Context(), db); err != nil {
			respondWithError(c, http.StatusServiceUnavailable, route, "database unavailable")
			return
		}

		page, limit, err := parsePaginationParams(c.Query("page"), c.Query("limit"))
		if err != nil {
			respondWithError(c, http.StatusBadRequest, route, "invalid pagination params")
			return
		}

		// Include products that are explicitly active as well as legacy entries
		// where the isActive flag might be missing.
		filter := bson.M{
			"isActive":  bson.M{"$ne": false},
			"isDeleted": bson.M{"$ne": true},
		}

		if category := strings.TrimSpace(c.Query("category")); category != "" {
			filter["category"] = bson.M{"$in": []string{category}}
		}

		if search := strings.TrimSpace(c.Query("search")); search != "" {
			filter["name"] = bson.M{"$regex": search, "$options": "i"}
		}

		findOptions := options.Find().
			SetSkip((page - 1) * limit).
			SetLimit(limit).
			SetSort(bson.D{{Key: "createdAt", Value: -1}})

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		cursor, err := db.Collection("products").Find(ctx, filter, findOptions)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, route, "db error")
			return
		}
		defer cursor.Close(ctx)

		products, err := decodeProducts(ctx, cursor)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, route, "decode error")
			return
		}

		log.Printf("[%s] returning %d products", route, len(products))
		c.JSON(http.StatusOK, products)
	}
}
func GetCampaignProducts(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		const route = "GET /products/campaign"
		defer handlePanic(c, route)

		if err := ensureDBConnection(c.Request.Context(), db); err != nil {
			respondWithError(c, http.StatusServiceUnavailable, route, "database unavailable")
			return
		}

		limit := parseCampaignLimit(c.Query("limit"))
		sortBy, err := parseCampaignSort(c.Query("sort"))
		if err != nil {
			respondWithError(c, http.StatusBadRequest, route, "invalid sort value")
			return
		}

		filter := bson.M{
			"isActive":   true,
			"isCampaign": true,
		}

		findOptions := options.Find().
			SetLimit(limit).
			SetSort(sortBy)

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		cursor, err := db.Collection("products").Find(ctx, filter, findOptions)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, route, "db error")
			return
		}
		defer cursor.Close(ctx)

		products, err := decodeProducts(ctx, cursor)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, route, "decode error")
			return
		}

		c.JSON(http.StatusOK, products)
	}
}

func parseCampaignLimit(rawLimit string) int64 {
	const (
		defaultLimit = int64(12)
		maxLimit     = int64(30)
	)

	if rawLimit == "" {
		return defaultLimit
	}

	limit, err := strconv.ParseInt(rawLimit, 10, 64)
	if err != nil || limit <= 0 {
		return defaultLimit
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
}

func parseCampaignSort(rawSort string) (bson.D, error) {
	sortValue := strings.TrimSpace(rawSort)
	if sortValue == "" || sortValue == "newest" {
		return bson.D{{Key: "createdAt", Value: -1}}, nil
	}
	switch sortValue {
	case "price_asc":
		return bson.D{{Key: "price", Value: 1}}, nil
	case "price_desc":
		return bson.D{{Key: "price", Value: -1}}, nil
	default:
		return bson.D{}, fmt.Errorf("invalid sort: %s", sortValue)
	}
}
