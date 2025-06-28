package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"-"` // Omit ID from JSON response
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"-"` // Hashed, but even so, don't include password in JSON response
}

type Bank struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	UserID         primitive.ObjectID `bson:"userId" json:"-"`
	BankName       string             `bson:"bankName" json:"bankName"`
	ClaimedCapital int64              `bson:"claimedCapital" json:"claimedCapital"`
}

type Asset struct {
	ID       primitive.ObjectID `bson:"_id" json:"-"`
	BankID   primitive.ObjectID `bson:"bankId" json:"-"`
	Amount   int64              `bson:"amount" json:"amount"`
	AssetType string            `bson:"assetType" json:"assetType"` // Using string for asset type (Cash, Stocks, Bonds, Crypto, etc.)
}

type HistoricalPerformance struct {
	ID        primitive.ObjectID `bson:"_id" json:"-"`
	Day       int               `bson:"day" json:"day"`
	BankID    primitive.ObjectID `bson:"bankId" json:"-"`
	Value     int64             `bson:"value" json:"value"`
	IsClaimed bool              `bson:"isClaimed" json:"isClaimed"`
}

type BankResponse struct {
	ID             string  `json:"id"`
	BankName       string  `json:"bankName"`
	ClaimedCapital int64   `json:"claimedCapital"`
	ActualCapital  int64   `json:"actualCapital"`
	Assets         []Asset `json:"assets"`
}

// PerformanceHistoryResponse represents the response structure for performance history
type PerformanceHistoryResponse struct {
	ClaimedHistory []HistoricalPerformanceResponse `json:"claimedHistory"`
	ActualHistory  []HistoricalPerformanceResponse `json:"actualHistory"`
}

// HistoricalPerformanceResponse represents a single day's performance value
type HistoricalPerformanceResponse struct {
	Day   int   `json:"day"`
	Value int64 `json:"value"`
}
