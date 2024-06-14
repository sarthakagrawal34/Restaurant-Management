package models

import (
	"restaurant-management/database"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderItem struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	OrderItemId string             `json:"order_item_id" bson:"order_item_id"`
	FoodId      *string            `json:"food_id" bson:"food_id" validate:"required"`
	OrderId     *string            `json:"order_id" bson:"order_id" validate:"required"`
	Quantity    *string            `json:"quantity" bson:"quantity" validate:"required,eq=S|eq=M|eq=L"`
	UnitPrice   *float64           `json:"unit_price" bson:"unit_price" validate:"required"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type OrderItemPack struct {
	TableId    *string     `json:"table_id"`
	OrderItems []OrderItem `json:"order_items"`
}

var OrderItemCollection *mongo.Collection = database.OpenCollection(database.MongoClient, "order_item")
