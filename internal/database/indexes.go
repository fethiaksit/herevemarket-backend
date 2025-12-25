package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureProductIndexes(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	indexes := db.Collection("products").Indexes()

	barcodeIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "barcode", Value: 1}},
		Options: options.Index().
			SetName("barcode_unique").
			SetUnique(true).
			SetPartialFilterExpression(bson.M{
				"barcode": bson.M{
					"$exists": true,
					"$ne":     "",
				},
			}),
	}

	log.Println("EnsureProductIndexes: creating barcode_unique index")
	_, err := indexes.CreateOne(ctx, barcodeIndex)
	if err != nil {
		log.Println("EnsureProductIndexes: barcode index error:", err)
		return err
	}
	log.Println("EnsureProductIndexes: barcode_unique index created")
	return nil
}
