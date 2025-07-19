package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PendingTransactionResponse struct {
	Id            primitive.ObjectID `bson:"_id" json:"id"`
	SourceBankId  primitive.ObjectID `bson:"sourceBankId" json:"sourceBankId"`
	TargetAssetId primitive.ObjectID `bson:"targetAssetId" json:"targetAssetId"`
	Amount        int64              `bson:"amount" json:"amount"` // Internal: Positive = buy, negative = sell
}

type PendingTransactionRequest struct {
	SourceBankId  string `json:"sourceBankId"`
	TargetAssetId string `json:"targetAssetId"`
	Amount        int64  `json:"amount"` // Always positive value from client
}
