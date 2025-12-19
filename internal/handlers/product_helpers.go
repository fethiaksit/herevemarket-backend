package handlers

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"backend/internal/models"
)

func normalizeProductDocument(raw bson.M) (models.Product, error) {
	if cat, ok := raw["category"].(string); ok {
		raw["category"] = []string{cat}
	}

	if val, ok := raw["isCampaign"]; ok {
		switch typed := val.(type) {
		case string:
			raw["isCampaign"] = typed == "true"
		case bool:
			// already bool, keep as is
		default:
			raw["isCampaign"] = false
		}
	} else {
		raw["isCampaign"] = false
	}

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
