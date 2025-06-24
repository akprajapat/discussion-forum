package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Body      string             `bson:"body" json:"body"`
	Author    primitive.ObjectID `bson:"author" json:"author"`
	Answer    primitive.ObjectID `bson:"answer" json:"answer"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
