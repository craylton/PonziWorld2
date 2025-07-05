package repositories

import (
	"context"
	"ponziworld/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type assetRepository struct {
	collection *mongo.Collection
}

// NewAssetRepository creates a new asset repository
func NewAssetRepository(database *mongo.Database) AssetRepository {
	return &assetRepository{
		collection: database.Collection("assets"),
	}
}

func (r *assetRepository) Create(ctx context.Context, asset *models.Asset) error {
	_, err := r.collection.InsertOne(ctx, asset)
	return err
}

func (r *assetRepository) FindByBankID(ctx context.Context, bankID primitive.ObjectID) ([]models.Asset, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"bankId": bankID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var assets []models.Asset
	if err = cursor.All(ctx, &assets); err != nil {
		return nil, err
	}
	return assets, nil
}

func (r *assetRepository) CalculateActualCapital(ctx context.Context, bankID primitive.ObjectID) (int64, error) {
	assets, err := r.FindByBankID(ctx, bankID)
	if err != nil {
		return 0, err
	}

	var actualCapital int64 = 0
	for _, asset := range assets {
		actualCapital += asset.Amount
	}
	return actualCapital, nil
}
