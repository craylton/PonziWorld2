package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"-"` // Omit ID from JSON response
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"-"` // Don't include password in JSON response
}

type Bank struct {
	ID             primitive.ObjectID `bson:"_id" json:"-"` // Omit ID from JSON response
	UserID         primitive.ObjectID `bson:"userId" json:"-"` // Link to User
	BankName       string             `bson:"bankName" json:"bankName"`
	ClaimedCapital int64              `bson:"claimedCapital" json:"claimedCapital"`
}

type Asset struct {
	ID       primitive.ObjectID `bson:"_id" json:"-"` // Omit ID from JSON response
	BankID   primitive.ObjectID `bson:"bankId" json:"-"` // Link to Bank
	Amount   int64              `bson:"amount" json:"amount"`
	AssetType string            `bson:"assetType" json:"assetType"` // Using string for asset type (Cash, Stocks, Bonds, Crypto, etc.)
}

// BankResponse represents the response structure for bank data
type BankResponse struct {
	BankName       string  `json:"bankName"`
	ClaimedCapital int64   `json:"claimedCapital"`
	ActualCapital  int64   `json:"actualCapital"` // Calculated from assets
	Assets         []Asset `json:"assets"`
}
