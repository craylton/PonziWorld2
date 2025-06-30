package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Game struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CurrentDay int                `bson:"currentDay" json:"currentDay"`
}
