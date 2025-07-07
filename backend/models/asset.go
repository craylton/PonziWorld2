package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Asset struct {
	Id          primitive.ObjectID `bson:"_id" json:"-"`
	BankId      primitive.ObjectID `bson:"bankId" json:"-"`
	Amount      int64              `bson:"amount" json:"amount"`
	AssetTypeId primitive.ObjectID `bson:"assetTypeId" json:"assetTypeId"`
}

type AssetResponse struct {
	Amount      int64  `json:"amount"`
	AssetTypeId string `json:"assetTypeId"`
	AssetType   string `json:"assetType"`
}
