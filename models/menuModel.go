package models

import (
	"restaurant-management/database"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Menu struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	MenuId    string             `json:"menu_id" bson:"menu_id"`
	Name      *string            `json:"name" bson:"name" validate:"required,min=2,max=100"`
	Category  *string            `json:"category" bson:"category" validate:"required,min=2,max=100"`
	StartDate *time.Time         `json:"start_date" bson:"start_date" validate:"required"`
	EndDate   *time.Time         `json:"end_date" bson:"end_date" validate:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

var MenuCollection *mongo.Collection = database.OpenCollection(database.MongoClient, "menu")
