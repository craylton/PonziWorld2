package main

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents a user in the database
type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username       string             `bson:"username" json:"username"`
	BankName       string             `bson:"bankName" json:"bankName"`
	ClaimedCapital int64              `bson:"claimedCapital" json:"claimedCapital"`
	ActualCapital  int64              `bson:"actualCapital" json:"actualCapital"`
}
