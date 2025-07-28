package repositories

import (
	"context"
	"ponziworld/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type investmentRepository struct {
	collection *mongo.Collection
}

// NewInvestmentRepository creates a new investment repository
func NewInvestmentRepository(database *mongo.Database) InvestmentRepository {
	return &investmentRepository{
		collection: database.Collection("investments"),
	}
}

func (r *investmentRepository) Create(ctx context.Context, investment *models.Investment) error {
	_, err := r.collection.InsertOne(ctx, investment)
	return err
}

func (r *investmentRepository) FindBySourceBankID(ctx context.Context, sourceBankID primitive.ObjectID) ([]models.Investment, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"sourceBankId": sourceBankID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var investments []models.Investment
	if err = cursor.All(ctx, &investments); err != nil {
		return nil, err
	}
	return investments, nil
}

func (r *investmentRepository) FindBySourceIdAndTargetId(
	ctx context.Context,
	sourceBankID,
	targetAssetId primitive.ObjectID,
) (*models.Investment, error) {
	var investment models.Investment
	err := r.collection.FindOne(ctx, bson.M{"sourceBankId": sourceBankID, "targetAssetId": targetAssetId}).Decode(&investment)
	if err != nil {
		return nil, err
	}
	return &investment, nil
}

func (r *investmentRepository) CalculateActualCapital(ctx context.Context, bankID primitive.ObjectID) (int64, error) {
	investments, err := r.FindBySourceBankID(ctx, bankID)
	if err != nil {
		return 0, err
	}

	var actualCapital int64 = 0
	for _, investment := range investments {
		actualCapital += investment.Amount
	}
	return actualCapital, nil
}

func (r *investmentRepository) UpdateAmount(
	ctx context.Context,
	sourceBankID,
	targetAssetId primitive.ObjectID,
	newAmount int64,
) error {
	filter := bson.M{"sourceBankId": sourceBankID, "targetAssetId": targetAssetId}
	update := bson.M{"$set": bson.M{"amount": newAmount}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *investmentRepository) DeleteBySourceIdAndTargetId(
	ctx context.Context,
	sourceBankID,
	targetAssetId primitive.ObjectID,
) error {
	filter := bson.M{"sourceBankId": sourceBankID, "targetAssetId": targetAssetId}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
