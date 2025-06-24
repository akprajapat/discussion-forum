package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Question struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Title     string               `bson:"title" json:"title"`
	Body      string               `bson:"body" json:"body"`
	Author    primitive.ObjectID   `bson:"author" json:"author"`
	Votes     int                  `bson:"votes" json:"votes"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	Tags      []string             `bson:"tags" json:"tags"`
	Answers   []primitive.ObjectID `bson:"answers" json:"answers"`
}
