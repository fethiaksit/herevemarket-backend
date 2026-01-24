package handlers

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"backend/internal/models"
)

func normalizeProductDocument(raw bson.M) (models.Product, error) {
	if cat, ok := raw["category"].(string); ok {
		raw["category"] = []string{cat}
	}

	if rawCategoryIDs, ok := raw["categoryIds"]; ok {
		converted, err := coerceCategoryIDs(rawCategoryIDs)
		if err != nil {
			return models.Product{}, err
		}
		if converted != nil {
			raw["categoryIds"] = converted
		}
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

	if val, ok := raw["stock"]; ok {
		switch typed := val.(type) {
		case int32:
			raw["stock"] = int(typed)
		case int64:
			raw["stock"] = int(typed)
		case float64:
			raw["stock"] = int(typed)
		case int:
			raw["stock"] = typed
		default:
			raw["stock"] = 0
		}
	} else {
		raw["stock"] = 0
	}

	data, err := bson.Marshal(raw)
	if err != nil {
		return models.Product{}, err
	}

	var p models.Product
	if err := bson.Unmarshal(data, &p); err != nil {
		return models.Product{}, err
	}

	p.InStock = p.Stock > 0

	return p, nil
}

func coerceCategoryIDs(value interface{}) ([]primitive.ObjectID, error) {
	switch typed := value.(type) {
	case []primitive.ObjectID:
		return typed, nil
	case []string:
		if len(typed) == 0 {
			return []primitive.ObjectID{}, nil
		}
		ids, err := parseCategoryObjectIDs(typed)
		if err != nil {
			return nil, err
		}
		return ids, nil
	case []interface{}:
		values := make([]string, 0, len(typed))
		for _, item := range typed {
			str, ok := item.(string)
			if !ok {
				return nil, nil
			}
			values = append(values, str)
		}
		if len(values) == 0 {
			return []primitive.ObjectID{}, nil
		}
		ids, err := parseCategoryObjectIDs(values)
		if err != nil {
			return nil, err
		}
		return ids, nil
	default:
		return nil, nil
	}
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
