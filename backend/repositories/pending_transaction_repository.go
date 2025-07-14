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

func (r *PendingTransactionRepositoryImpl) Create(ctx context.Context, transaction *models.PendingTransaction) error {
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

func (r *PendingTransactionRepositoryImpl) FindByBuyerBankID(
	ctx context.Context,
	buyerBankID primitive.ObjectID,
) ([]models.PendingTransaction, error) {
	filter := bson.M{"buyerBankId": buyerBankID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []models.PendingTransaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *PendingTransactionRepositoryImpl) FindByBuyerBankIDAndAssetID(
	ctx context.Context,
	buyerBankID,
	assetID primitive.ObjectID,
) ([]models.PendingTransaction, error) {
	filter := bson.M{
		"buyerBankId": buyerBankID,
		"assetId":     assetID,
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []models.PendingTransaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
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
