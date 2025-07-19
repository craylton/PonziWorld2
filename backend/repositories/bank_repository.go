package repositories

import (
	"context"
	"ponziworld/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type bankRepository struct {
	collection *mongo.Collection
}

// NewBankRepository creates a new bank repository
func NewBankRepository(database *mongo.Database) BankRepository {
	return &bankRepository{
		collection: database.Collection("banks"),
	}
}

func (r *bankRepository) Create(ctx context.Context, bank *models.Bank) error {
	_, err := r.collection.InsertOne(ctx, bank)
	return err
}

func (r *bankRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Bank, error) {
	var bank models.Bank
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&bank)
	if err != nil {
		return nil, err
	}
	return &bank, nil
}

func (r *bankRepository) FindAllByPlayerID(ctx context.Context, playerID primitive.ObjectID) ([]models.Bank, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"playerId": playerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var banks []models.Bank
	if err = cursor.All(ctx, &banks); err != nil {
		return nil, err
	}
	return banks, nil
}
