package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Asset struct {
	Id          primitive.ObjectID `bson:"_id" json:"-"`
	BankId      primitive.ObjectID `bson:"bankId" json:"-"`
	Amount      int64              `bson:"amount" json:"amount"`
	AssetTypeId primitive.ObjectID `bson:"assetTypeId" json:"assetTypeId"`
}

type AvailableAssetResponse struct {
	AssetTypeId         string `json:"assetTypeId"`
	AssetType           string `json:"assetType"`
	IsInvestedOrPending bool   `json:"isInvestedOrPending"`
}

type AssetDetailsResponse struct {
	AssetId        string                          `json:"assetId"`
	Name           string                          `json:"name"`
	InvestedAmount int64                           `json:"investedAmount"`
	PendingAmount  int64                           `json:"pendingAmount"`
	HistoricalData []HistoricalPerformanceResponse `json:"historicalData"`
}
