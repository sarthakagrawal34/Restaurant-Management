package models

import (
	"restaurant-management/database"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RestaurantTable struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	TableId        string             `json:"table_id" bson:"table_id" validate:"required"`
	NumberOfGuests *int               `json:"number_of_guests" bson:"number_of_guests" validate:"required"`
	TableNumber    *int               `json:"table_number" bson:"table_number" validate:"required"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}

var RestaurantTableCollection *mongo.Collection = database.OpenCollection(database.MongoClient, "restaurant_table")
