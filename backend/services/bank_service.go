package services

import (
	"context"
	"errors"
	"ponziworld/backend/models"
	"ponziworld/backend/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrBankNotFound    = errors.New("bank not found")
	ErrUnauthorized    = errors.New("unauthorized access")
)

type BankService struct {
	playerRepo repositories.PlayerRepository
	bankRepo   repositories.BankRepository
	assetRepo  repositories.AssetRepository
}

func NewBankService(
	playerRepo repositories.PlayerRepository,
	bankRepo repositories.BankRepository,
	assetRepo repositories.AssetRepository,
) *BankService {
	return &BankService{
		playerRepo: playerRepo,
		bankRepo:   bankRepo,
		assetRepo:  assetRepo,
	}
}

func (s *BankService) GetBankByUsername(ctx context.Context, username string) (*models.BankResponse, error) {
	// Find the player
	player, err := s.playerRepo.FindByUsername(ctx, username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrPlayerNotFound
		}
		return nil, err
	}

	// Find the bank for this player
	bank, err := s.bankRepo.FindByPlayerID(ctx, player.Id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrBankNotFound
		}
		return nil, err
	}

	// Get all assets for this bank
	assets, err := s.assetRepo.FindByBankID(ctx, bank.Id)
	if err != nil {
		return nil, err
	}

	// Calculate actual capital
	actualCapital, err := s.assetRepo.CalculateActualCapital(ctx, bank.Id)
	if err != nil {
		return nil, err
	}

	// Create response
	response := &models.BankResponse{
		Id:             bank.Id.Hex(),
		BankName:       bank.BankName,
		ClaimedCapital: bank.ClaimedCapital,
		ActualCapital:  actualCapital,
		Assets:         assets,
	}

	return response, nil
}

func (s *BankService) CreateBank(ctx context.Context, playerID primitive.ObjectID, bankName string, initialCapital int64) (*models.Bank, error) {
	bank := &models.Bank{
		Id:             primitive.NewObjectID(),
		PlayerId:       playerID,
		BankName:       bankName,
		ClaimedCapital: initialCapital,
	}

	err := s.bankRepo.Create(ctx, bank)
	if err != nil {
		return nil, err
	}

	return bank, nil
}

func (s *BankService) ValidateBankOwnership(ctx context.Context, username string, bankID primitive.ObjectID) error {
	player, err := s.playerRepo.FindByUsername(ctx, username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrPlayerNotFound
		}
		return err
	}

	bank, err := s.bankRepo.FindByID(ctx, bankID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrBankNotFound
		}
		return err
	}

	if bank.PlayerId != player.Id {
		return ErrUnauthorized
	}

	return nil
}
