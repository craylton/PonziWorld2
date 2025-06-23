package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username       string             `bson:"username" json:"username"`
	Password       string             `bson:"password" json:"-"` // Don't include password in JSON response
	BankName       string             `bson:"bankName" json:"bankName"`
	ClaimedCapital int64              `bson:"claimedCapital" json:"claimedCapital"`
	ActualCapital  int64              `bson:"actualCapital" json:"actualCapital"`
}
