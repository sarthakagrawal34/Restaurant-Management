package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	OrderDate      time.Time          `json:"order_date" bson:"order_date" validate:"required"`
	OrderId        string             `json:"order_id" bson:"order_id" validate:"required"`
	TableId        *string            `json:"table_id" bson:"table_id" validate:"required"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}
