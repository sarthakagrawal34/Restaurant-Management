package models

import (
	"restaurant-management/database"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Order struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	OrderDate time.Time          `json:"order_date" bson:"order_date" validate:"required"`
	OrderId   string             `json:"order_id" bson:"order_id"`
	TableId   *string            `json:"table_id" bson:"table_id" validate:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

var OrderCollection *mongo.Collection = database.OpenCollection(database.MongoClient, "order")
