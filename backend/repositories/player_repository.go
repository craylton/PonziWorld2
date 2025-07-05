package repositories

import (
	"context"
	"ponziworld/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type playerRepository struct {
	collection *mongo.Collection
}

// NewPlayerRepository creates a new player repository
func NewPlayerRepository(database *mongo.Database) PlayerRepository {
	return &playerRepository{
		collection: database.Collection("players"),
	}
}

func (r *playerRepository) Create(ctx context.Context, player *models.Player) error {
	_, err := r.collection.InsertOne(ctx, player)
	return err
}

func (r *playerRepository) FindByUsername(ctx context.Context, username string) (*models.Player, error) {
	var player models.Player
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&player)
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *playerRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Player, error) {
	var player models.Player
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&player)
	if err != nil {
		return nil, err
	}
	return &player, nil
}
