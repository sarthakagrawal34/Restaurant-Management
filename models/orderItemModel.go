package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	OrderItemId string             `json:"order_item_id" bson:"order_item_id" validate:"required"`
	FoodId      *string            `json:"food_id" bson:"food_id" validate:"required"`
	OrderId     *string            `json:"order_id" bson:"order_id" validate:"required"`
	Quantity    *string            `json:"quantity" bson:"quantity" validate:"required,eq=S|eq=M|eq=L"`
	UnitPrice   *float64           `json:"unit_price" bson:"unit_price" validate:"required"`
	OrderDate   time.Time          `json:"order_date" bson:"order_date" validate:"required"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
