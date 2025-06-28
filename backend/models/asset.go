package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Asset struct {
	ID        primitive.ObjectID `bson:"_id" json:"-"`
	BankID    primitive.ObjectID `bson:"bankId" json:"-"`
	Amount    int64              `bson:"amount" json:"amount"`
	AssetType string             `bson:"assetType" json:"assetType"`
}
