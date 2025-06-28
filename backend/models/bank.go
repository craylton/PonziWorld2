package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Bank struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	UserID         primitive.ObjectID `bson:"userId" json:"-"`
	BankName       string             `bson:"bankName" json:"bankName"`
	ClaimedCapital int64              `bson:"claimedCapital" json:"claimedCapital"`
}

type BankResponse struct {
	ID             string  `json:"id"`
	BankName       string  `json:"bankName"`
	ClaimedCapital int64   `json:"claimedCapital"`
	ActualCapital  int64   `json:"actualCapital"`
	Assets         []Asset `json:"assets"`
}
