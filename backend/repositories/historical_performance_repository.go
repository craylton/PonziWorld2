package repositories

import (
	"context"
	"ponziworld/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type historicalPerformanceRepository struct {
	collection *mongo.Collection
}

// NewHistoricalPerformanceRepository creates a new historical performance repository
func NewHistoricalPerformanceRepository(database *mongo.Database) HistoricalPerformanceRepository {
	return &historicalPerformanceRepository{
		collection: database.Collection("historicalPerformance"),
	}
}

func (r *historicalPerformanceRepository) Create(ctx context.Context, performance *models.HistoricalPerformance) error {
	_, err := r.collection.InsertOne(ctx, performance)
	return err
}

func (r *historicalPerformanceRepository) CreateMany(ctx context.Context, performances []models.HistoricalPerformance) error {
	if len(performances) == 0 {
		return nil
	}

	documents := make([]interface{}, len(performances))
	for i, performance := range performances {
		documents[i] = performance
	}

	_, err := r.collection.InsertMany(ctx, documents)
	return err
}

func (r *historicalPerformanceRepository) FindByBankIDAndDateRange(ctx context.Context, bankID primitive.ObjectID, startDay, endDay int) ([]models.HistoricalPerformance, error) {
	filter := bson.M{
		"bankId": bankID,
		"day":    bson.M{"$gt": startDay, "$lte": endDay},
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"day": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var history []models.HistoricalPerformance
	if err = cursor.All(ctx, &history); err != nil {
		return nil, err
	}
	return history, nil
}

func (r *historicalPerformanceRepository) FindClaimedByBankIDAndDateRange(ctx context.Context, bankID primitive.ObjectID, startDay, endDay int) ([]models.HistoricalPerformance, error) {
	filter := bson.M{
		"bankId":    bankID,
		"day":       bson.M{"$gt": startDay, "$lte": endDay},
		"isClaimed": true,
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"day": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var history []models.HistoricalPerformance
	if err = cursor.All(ctx, &history); err != nil {
		return nil, err
	}
	return history, nil
}

func (r *historicalPerformanceRepository) FindActualByBankIDAndDateRange(ctx context.Context, bankID primitive.ObjectID, startDay, endDay int) ([]models.HistoricalPerformance, error) {
	filter := bson.M{
		"bankId":    bankID,
		"day":       bson.M{"$gt": startDay, "$lte": endDay},
		"isClaimed": false,
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"day": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var history []models.HistoricalPerformance
	if err = cursor.All(ctx, &history); err != nil {
		return nil, err
	}
	return history, nil
}
