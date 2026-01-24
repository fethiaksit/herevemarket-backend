package database

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"backend/internal/models"
)

func RunCategoryIDMigration(ctx context.Context, db *mongo.Database) error {
	log.Println("RunCategoryIDMigration: loading categories")
	categoryCursor, err := db.Collection("categories").Find(ctx, bson.M{})
	if err != nil {
		return err
	}

	var categories []models.Category
	if err := categoryCursor.All(ctx, &categories); err != nil {
		return err
	}

	nameToID := make(map[string]primitive.ObjectID, len(categories))
	for _, category := range categories {
		nameToID[category.Name] = category.ID
	}

	filter := bson.M{
		"category": bson.M{"$exists": true},
		"$or": []bson.M{
			{"categoryIds": bson.M{"$exists": false}},
			{"categoryIds": bson.M{"$size": 0}},
		},
	}

	log.Println("RunCategoryIDMigration: scanning products")
	productCursor, err := db.Collection("products").Find(ctx, filter)
	if err != nil {
		return err
	}
	defer productCursor.Close(ctx)

	updatedCount := 0
	skippedCount := 0
	missingNames := map[string]struct{}{}

	for productCursor.Next(ctx) {
		var raw bson.M
		if err := productCursor.Decode(&raw); err != nil {
			return err
		}

		productID, ok := raw["_id"].(primitive.ObjectID)
		if !ok {
			skippedCount++
			continue
		}

		legacyNames := extractLegacyCategoryNames(raw["category"])
		normalizedNames := normalizeCategoryNames(legacyNames)
		if len(normalizedNames) == 0 {
			skippedCount++
			continue
		}

		categoryIDs := make([]primitive.ObjectID, 0, len(normalizedNames))
		missing := false
		for _, name := range normalizedNames {
			objectID, ok := nameToID[name]
			if !ok {
				missingNames[name] = struct{}{}
				missing = true
				break
			}
			categoryIDs = append(categoryIDs, objectID)
		}

		if missing {
			skippedCount++
			continue
		}

		update := bson.M{
			"$set": bson.M{
				"categoryIds":   categoryIDs,
				"categoryNames": normalizedNames,
			},
		}

		if _, err := db.Collection("products").UpdateOne(ctx, bson.M{"_id": productID}, update); err != nil {
			return err
		}

		updatedCount++
		if updatedCount%100 == 0 {
			log.Printf("RunCategoryIDMigration: updated %d products", updatedCount)
		}
	}

	if err := productCursor.Err(); err != nil {
		return err
	}

	report := formatMissingCategoryNames(missingNames)
	if report != "" {
		log.Printf("RunCategoryIDMigration: missing category names: %s", report)
	}

	log.Printf("RunCategoryIDMigration: done updated=%d skipped=%d", updatedCount, skippedCount)
	return nil
}

func extractLegacyCategoryNames(value interface{}) []string {
	switch typed := value.(type) {
	case string:
		return []string{typed}
	case []string:
		return typed
	case []interface{}:
		values := make([]string, 0, len(typed))
		for _, item := range typed {
			name, ok := item.(string)
			if !ok {
				continue
			}
			values = append(values, name)
		}
		return values
	default:
		return nil
	}
}

func normalizeCategoryNames(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))

	for _, value := range values {
		name := strings.TrimSpace(value)
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

func formatMissingCategoryNames(missing map[string]struct{}) string {
	if len(missing) == 0 {
		return ""
	}

	names := make([]string, 0, len(missing))
	for name := range missing {
		names = append(names, name)
	}
	sort.Strings(names)
	return fmt.Sprintf("%s (count=%d) at %s", strings.Join(names, ", "), len(names), time.Now().Format(time.RFC3339))
}
