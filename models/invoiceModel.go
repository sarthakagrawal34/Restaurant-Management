package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	InvoiceId      string             `json:"invoice_id" bson:"invoice_id" validate:"required"`
	OrderId        *string            `json:"order_id" bson:"order_id" validate:"required"`
	PaymentMethod  *string            `json:"payment_method" bson:"payment_method" validate:"required,eq=CARD|eq=CASH|eq="`
	PaymentStatus  *float64           `json:"payment_status" bson:"payment_status" validate:"required,eq=PENDING|eq=PAID|eq=CANCELLED"`
	PaymentDueDate time.Time          `json:"payment_due_date" bson:"payment_due_date"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}
