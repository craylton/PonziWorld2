package services

import (
	"context"
	"errors"
	"ponziworld/backend/models"
	"ponziworld/backend/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPlayerNotFound     = errors.New("player not found")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUsernameExists     = errors.New("username already exists")
)

type AuthService struct {
	playerRepo repositories.PlayerRepository
}

func NewAuthService(playerRepo repositories.PlayerRepository) *AuthService {
	return &AuthService{
		playerRepo: playerRepo,
	}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*models.Player, error) {
	player, err := s.playerRepo.FindByUsername(ctx, username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return player, nil
}

func (s *AuthService) CreatePlayer(ctx context.Context, username, password string) (*models.Player, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	player := &models.Player{
		Id:       primitive.NewObjectID(),
		Username: username,
		Password: string(hashedPassword),
		IsAdmin:  false,
	}

	err = s.playerRepo.Create(ctx, player)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrUsernameExists
		}
		return nil, err
	}

	return player, nil
}

func (s *AuthService) GetPlayerByUsername(ctx context.Context, username string) (*models.Player, error) {
	player, err := s.playerRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return player, nil
}
