package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Asset struct {
	Id        primitive.ObjectID `bson:"_id" json:"-"`
	BankId    primitive.ObjectID `bson:"bankId" json:"-"`
	Amount    int64              `bson:"amount" json:"amount"`
	AssetType string             `bson:"assetType" json:"assetType"`
}
