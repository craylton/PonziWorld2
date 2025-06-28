package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type HistoricalPerformance struct {
	ID        primitive.ObjectID `bson:"_id" json:"-"`
	Day       int               `bson:"day" json:"day"`
	BankID    primitive.ObjectID `bson:"bankId" json:"-"`
	Value     int64             `bson:"value" json:"value"`
	IsClaimed bool              `bson:"isClaimed" json:"isClaimed"`
}

type PerformanceHistoryResponse struct {
	ClaimedHistory []HistoricalPerformanceResponse `json:"claimedHistory"`
	ActualHistory  []HistoricalPerformanceResponse `json:"actualHistory"`
}

type HistoricalPerformanceResponse struct {
	Day   int   `json:"day"`
	Value int64 `json:"value"`
}
