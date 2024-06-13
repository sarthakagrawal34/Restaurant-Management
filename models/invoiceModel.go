package models

import (
	"restaurant-management/database"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Invoice struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	InvoiceId      string             `json:"invoice_id" bson:"invoice_id" validate:"required"`
	OrderId        *string            `json:"order_id" bson:"order_id" validate:"required"`
	PaymentMethod  *string            `json:"payment_method" bson:"payment_method" validate:"required,eq=CARD|eq=CASH|eq="`
	PaymentStatus  *string            `json:"payment_status" bson:"payment_status" validate:"required,eq=PENDING|eq=PAID|eq=CANCELLED"`
	PaymentDueDate time.Time          `json:"payment_due_date" bson:"payment_due_date"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}

type InvoiceViewFormat struct {
	InvoiceId      string
	PaymentMethod  string
	OrderId        string
	PaymentStatus  *string
	PaymentDue     interface{}
	TableNumber    interface{}
	PaymentDueDate time.Time
	OrderDetails   interface{}
}

var InvoiceCollection *mongo.Collection = database.OpenCollection(database.MongoClient, "invoice")
