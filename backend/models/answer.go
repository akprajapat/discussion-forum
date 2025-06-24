package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Answer struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Body      string               `bson:"body" json:"body"`
	Author    primitive.ObjectID   `bson:"author" json:"author"`
	Question  primitive.ObjectID   `bson:"question" json:"question"`
	Votes     int                  `bson:"votes" json:"votes"`
	Comments  []primitive.ObjectID `bson:"comments" json:"comments"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
}
