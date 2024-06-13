package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	NoteId    string             `json:"note_id" bson:"note_id" validate:"required"`
	Text      *string            `json:"text" bson:"text" validate:"required"`
	Title     *string            `json:"title" bson:"title" validate:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
