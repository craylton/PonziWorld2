package repositories

import (
	"context"
	"ponziworld/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type assetTypeRepository struct {
	collection *mongo.Collection
}

// NewAssetTypeRepository creates a new asset type repository
func NewAssetTypeRepository(database *mongo.Database) AssetTypeRepository {
	return &assetTypeRepository{
		collection: database.Collection("assetTypes"),
	}
}

func (r *assetTypeRepository) Create(ctx context.Context, assetType *models.AssetType) error {
	_, err := r.collection.InsertOne(ctx, assetType)
	return err
}

func (r *assetTypeRepository) FindAll(ctx context.Context) ([]models.AssetType, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var assetTypes []models.AssetType
	if err = cursor.All(ctx, &assetTypes); err != nil {
		return nil, err
	}
	return assetTypes, nil
}

func (r *assetTypeRepository) FindByName(ctx context.Context, name string) (*models.AssetType, error) {
	var assetType models.AssetType
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&assetType)
	if err != nil {
		return nil, err
	}
	return &assetType, nil
}

func (r *assetTypeRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.AssetType, error) {
	var assetType models.AssetType
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&assetType)
	if err != nil {
		return nil, err
	}
	return &assetType, nil
}

func (r *assetTypeRepository) UpsertByName(ctx context.Context, assetType *models.AssetType) error {
	// Try to find existing asset type
	existing, err := r.FindByName(ctx, assetType.Name)
	if err == nil {
		// Update existing
		assetType.Id = existing.Id
		filter := bson.M{"_id": existing.Id}
		update := bson.M{"$set": assetType}
		_, err := r.collection.UpdateOne(ctx, filter, update)
		return err
	}
	
	// Create new if not found
	if err == mongo.ErrNoDocuments {
		return r.Create(ctx, assetType)
	}
	
	return err
}
