package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Player struct {
	Id       primitive.ObjectID `bson:"_id" json:"-"` // Omit ID from JSON response
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"-"` // Hashed, but even so, don't include password in JSON response
}
