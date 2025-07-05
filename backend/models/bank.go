package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Bank struct {
	Id             primitive.ObjectID `bson:"_id" json:"id"`
	PlayerId       primitive.ObjectID `bson:"playerId" json:"-"`
	BankName       string             `bson:"bankName" json:"bankName"`
	ClaimedCapital int64              `bson:"claimedCapital" json:"claimedCapital"`
}

type AssetResponse struct {
	Amount      int64  `json:"amount"`
	AssetTypeId string `json:"assetTypeId"`
	AssetType   string `json:"assetType"` // Include the asset type name for convenience
}

type BankResponse struct {
	Id             string          `json:"id"`
	BankName       string          `json:"bankName"`
	ClaimedCapital int64           `json:"claimedCapital"`
	ActualCapital  int64           `json:"actualCapital"`
	Assets         []AssetResponse `json:"assets"`
}
