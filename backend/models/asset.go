package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Investment struct {
	Id            primitive.ObjectID `bson:"_id" json:"-"`
	SourceBankId  primitive.ObjectID `bson:"sourceBankId" json:"-"`
	Amount        int64              `bson:"amount" json:"amount"`
	TargetAssetId primitive.ObjectID `bson:"targetAssetId" json:"targetAssetId"`
}

type InvestmentDetailsResponse struct {
	TargetAssetId  string                          `json:"targetAssetId"`
	Name           string                          `json:"name"`
	InvestedAmount int64                           `json:"investedAmount"`
	PendingAmount  int64                           `json:"pendingAmount"`
	HistoricalData []HistoricalPerformanceResponse `json:"historicalData"`
}
