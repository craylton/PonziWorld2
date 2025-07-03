package services

import (
	"context"
	"ponziworld/backend/repositories"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GameService struct {
	gameRepo repositories.GameRepository
}

func NewGameService(gameRepo repositories.GameRepository) *GameService {
	return &GameService{
		gameRepo: gameRepo,
	}
}

func (s *GameService) GetCurrentDay(ctx context.Context) (int, error) {
	currentDay, err := s.gameRepo.GetCurrentDay(ctx)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If no game state exists, create it starting at day 0
			err = s.gameRepo.CreateInitialGame(ctx, 0)
			if err != nil {
				return 0, err
			}
			return 0, nil
		}
		return 0, err
	}
	return currentDay, nil
}

func (s *GameService) NextDay(ctx context.Context) (int, error) {
	// Try to increment the day
	newDay, err := s.gameRepo.IncrementDay(ctx)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If no game state exists, create it starting at day 1 (since we're incrementing from 0)
			err = s.gameRepo.CreateInitialGame(ctx, 1)
			if err != nil {
				return 0, err
			}
			return 1, nil
		}
		return 0, err
	}
	return newDay, nil
}
