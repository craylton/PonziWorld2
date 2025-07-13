package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PendingTransaction struct {
	Id           primitive.ObjectID `bson:"_id" json:"id"`
	BuyerBankId  primitive.ObjectID `bson:"buyerBankId" json:"buyerBankId"`
	AssetId      primitive.ObjectID `bson:"assetId" json:"assetId"`
	Amount       int64              `bson:"amount" json:"amount"` // Internal: Positive = buy, negative = sell
	CreatedAt    primitive.DateTime `bson:"createdAt" json:"createdAt"`
}

type PendingTransactionRequest struct {
	BuyerBankId string `json:"buyerBankId"`
	AssetId     string `json:"assetId"`
	Amount      int64  `json:"amount"` // Always positive value from client
}
