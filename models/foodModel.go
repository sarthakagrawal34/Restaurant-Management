package models

import (
	"restaurant-management/database"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Food struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	FoodId    string             `json:"food_id" bson:"food_id" validate:"required"`
	MenuId    *string            `json:"menu_id" bson:"menu_id" validate:"required"`
	Name      *string            `json:"name" bson:"name" validate:"required,min=2,max=100"`
	Price     *float64           `json:"price" bson:"price" validate:"required"`
	FoodImage *string            `json:"food_image" bson:"food_image" validate:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

var FoodCollection *mongo.Collection = database.OpenCollection(database.MongoClient, "food")
