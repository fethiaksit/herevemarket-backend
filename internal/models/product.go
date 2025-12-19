package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	Price      float64            `bson:"price" json:"price"`
	Category   StringList         `bson:"category" json:"category"`
	ImageURL   string             `bson:"imageUrl" json:"imageUrl"`
	IsActive   bool               `bson:"isActive" json:"isActive"`
	IsCampaign bool               `bson:"isCampaign" json:"isCampaign"`
	IsDeleted  bool               `bson:"isDeleted" json:"isDeleted"`
	DeletedAt  *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}
