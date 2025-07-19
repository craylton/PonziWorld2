package repositories

import (
	"context"
	"errors"
	"ponziworld/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PendingTransactionRepositoryImpl struct {
	collection *mongo.Collection
}

func NewPendingTransactionRepository(database *mongo.Database) *PendingTransactionRepositoryImpl {
	return &PendingTransactionRepositoryImpl{
		collection: database.Collection("pendingTransactions"),
	}
}

func (r *PendingTransactionRepositoryImpl) Create(ctx context.Context, transaction *models.PendingTransactionResponse) error {
	if transaction.Id.IsZero() {
		transaction.Id = primitive.NewObjectID()
	}

	_, err := r.collection.InsertOne(ctx, transaction)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("pending transaction already exists")
		}
		return err
	}
	return nil
}

func (r *PendingTransactionRepositoryImpl) FindBySourceBankID(
	ctx context.Context,
	sourceBankID primitive.ObjectID,
) ([]models.PendingTransactionResponse, error) {
	filter := bson.M{"sourceBankId": sourceBankID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []models.PendingTransactionResponse
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *PendingTransactionRepositoryImpl) FindBySourceBankIDAndTargetAssetID(
	ctx context.Context,
	sourceBankID,
	targetAssetID primitive.ObjectID,
) ([]models.PendingTransactionResponse, error) {
	filter := bson.M{
		"sourceBankId": sourceBankID,
		"targetAssetId":     targetAssetID,
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []models.PendingTransactionResponse
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *PendingTransactionRepositoryImpl) SumPendingAmountBySourceBankIdAndTargetAssetId(
	ctx context.Context,
	sourceBankID,
	targetAssetID primitive.ObjectID,
) (int64, error) {
	transactions, err := r.FindBySourceBankIDAndTargetAssetID(ctx, sourceBankID, targetAssetID)
	if err != nil {
		return 0, err
	}

	var totalAmount int64
	for _, transaction := range transactions {
		totalAmount += transaction.Amount
	}

	return totalAmount, nil
}

func (r *PendingTransactionRepositoryImpl) UpdateAmount(
	ctx context.Context,
	id primitive.ObjectID,
	newAmount int64,
) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"amount": newAmount}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("pending transaction not found")
	}

	return nil
}

func (r *PendingTransactionRepositoryImpl) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("pending transaction not found")
	}
	return nil
}
