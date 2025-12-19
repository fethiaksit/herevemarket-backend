package handlers

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"backend/internal/models"
)

func parseOptionalBool(value interface{}, defaultValue bool) bool {
	switch typed := value.(type) {
	case string:
		return typed == "true"
	case bool:
		return typed
	default:
		return defaultValue
	}
}

func normalizeProductDocument(raw bson.M) (models.Product, error) {
	if cat, ok := raw["category"].(string); ok {
		raw["category"] = []string{cat}
	}

	raw["isCampaign"] = parseOptionalBool(raw["isCampaign"], false)
	raw["isActive"] = parseOptionalBool(raw["isActive"], true)
	raw["isDeleted"] = parseOptionalBool(raw["isDeleted"], false)

	data, err := bson.Marshal(raw)
	if err != nil {
		return models.Product{}, err
	}

	var p models.Product
	if err := bson.Unmarshal(data, &p); err != nil {
		return models.Product{}, err
	}

	return p, nil
}

func decodeProducts(ctx context.Context, cursor *mongo.Cursor) ([]models.Product, error) {
	products := make([]models.Product, 0)

	for cursor.Next(ctx) {
		var raw bson.M
		if err := cursor.Decode(&raw); err != nil {
			return nil, err
		}

		product, err := normalizeProductDocument(raw)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
