package repositories

import (
	"context"
	"ponziworld/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type gameRepository struct {
	collection *mongo.Collection
}

func NewGameRepository(database *mongo.Database) GameRepository {
	return &gameRepository{
		collection: database.Collection("game"),
	}
}

func (r *gameRepository) GetCurrentDay(ctx context.Context) (int, error) {
	var game models.Game
	err := r.collection.FindOne(ctx, bson.M{}).Decode(&game)
	if err != nil {
		return 0, err // Return the error (including ErrNoDocuments) to the service
	}
	return game.CurrentDay, nil
}

func (r *gameRepository) IncrementDay(ctx context.Context) (int, error) {
	filter := bson.M{}
	update := bson.M{"$inc": bson.M{"currentDay": 1}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var game models.Game
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&game)
	if err != nil {
		return 0, err
	}
	return game.CurrentDay, nil
}

func (r *gameRepository) CreateInitialGame(ctx context.Context, initialDay int) error {
	game := models.Game{
		ID:         primitive.NewObjectID(),
		CurrentDay: initialDay,
	}
	_, err := r.collection.InsertOne(ctx, game)
	return err
}
