package services

import (
	"context"
	"errors"
	"ponziworld/backend/models"
	"ponziworld/backend/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidAssetID      = errors.New("invalid asset ID")
	ErrInvalidAmount       = errors.New("amount must be positive")
	ErrTargetAssetNotFound = errors.New("target asset not found")
	ErrInvalidBankID       = errors.New("invalid bank ID")
	ErrSelfInvestment      = errors.New("bank cannot invest in itself")
	ErrUnauthorizedBank    = errors.New("bank is not owned by the current player")
)

type PendingTransactionService struct {
	pendingTransactionRepo repositories.PendingTransactionRepository
	bankRepo               repositories.BankRepository
	assetTypeRepo          repositories.AssetTypeRepository
	playerRepo             repositories.PlayerRepository
}

func NewPendingTransactionService(
	pendingTransactionRepo repositories.PendingTransactionRepository,
	bankRepo repositories.BankRepository,
	assetTypeRepo repositories.AssetTypeRepository,
	playerRepo repositories.PlayerRepository,
) *PendingTransactionService {
	return &PendingTransactionService{
		pendingTransactionRepo: pendingTransactionRepo,
		bankRepo:               bankRepo,
		assetTypeRepo:          assetTypeRepo,
		playerRepo:             playerRepo,
	}
}

func (s *PendingTransactionService) CreateBuyTransaction(
	ctx context.Context,
	sourceBankId,
	targetAssetId primitive.ObjectID,
	amount int64,
	username string,
) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	return s.createTransaction(ctx, sourceBankId, targetAssetId, amount, username)
}

func (s *PendingTransactionService) CreateSellTransaction(
	ctx context.Context,
	sourceBankId,
	targetAssetId primitive.ObjectID,
	amount int64,
	username string,
) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	// Convert to negative amount for internal use
	return s.createTransaction(ctx, sourceBankId, targetAssetId, -amount, username)
}

func (s *PendingTransactionService) createTransaction(
	ctx context.Context,
	sourceBankId,
	targetAssetId primitive.ObjectID,
	amount int64,
	username string,
) error {
	// Validate buyer bank exists and is owned by the current player
	sourceBank, err := s.bankRepo.FindByID(ctx, sourceBankId)
	if err != nil {
		return ErrInvalidBankID
	}

	// Validate bank ownership
	player, err := s.playerRepo.FindByUsername(ctx, username)
	if err != nil {
		return ErrPlayerNotFound
	}

	if sourceBank.PlayerId != player.Id {
		return ErrUnauthorizedBank
	}

	// Validate asset exists
	assetExists, err := s.validateTargetAssetExists(ctx, targetAssetId)
	if err != nil {
		return err
	}
	if !assetExists {
		return ErrTargetAssetNotFound
	}

	// Validate that bank is not investing in itself
	if sourceBankId == targetAssetId {
		return ErrSelfInvestment
	}

	// Check if there's an existing pending transaction for this bank-asset combination
	existingTransactions, err := s.pendingTransactionRepo.FindBySourceBankIDAndTargetAssetID(
		ctx,
		sourceBankId,
		targetAssetId,
	)
	if err != nil {
		return err
	}

	// If there's an existing transaction, combine them
	if len(existingTransactions) > 0 {
		existingTransaction := existingTransactions[0]
		newAmount := existingTransaction.Amount + amount

		// If the new amount is zero, delete the transaction
		if newAmount == 0 {
			return s.pendingTransactionRepo.Delete(ctx, existingTransaction.Id)
		}

		// Otherwise, update the existing transaction
		return s.pendingTransactionRepo.UpdateAmount(ctx, existingTransaction.Id, newAmount)
	}

	// Create new pending transaction if none exists
	transaction := &models.PendingTransactionResponse{
		SourceBankId:  sourceBankId,
		TargetAssetId: targetAssetId,
		Amount:        amount,
	}

	return s.pendingTransactionRepo.Create(ctx, transaction)
}

func (s *PendingTransactionService) validateTargetAssetExists(
	ctx context.Context,
	targetAssetId primitive.ObjectID,
) (bool, error) {
	// First check if it's an asset type
	_, err := s.assetTypeRepo.FindByID(ctx, targetAssetId)
	if err == nil {
		return true, nil
	}

	// Then check if it's a bank (since banks are also assets)
	_, err = s.bankRepo.FindByID(ctx, targetAssetId)
	if err == nil {
		return true, nil
	}

	return false, nil
}

func (s *PendingTransactionService) GetTransactionsByBuyerBankID(
	ctx context.Context,
	bankID primitive.ObjectID,
	username string,
) ([]models.PendingTransactionResponse, error) {
	// Validate bank exists and is owned by the current player
	bank, err := s.bankRepo.FindByID(ctx, bankID)
	if err != nil {
		return nil, ErrInvalidBankID
	}

	// Validate bank ownership
	player, err := s.playerRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, ErrPlayerNotFound
	}

	if bank.PlayerId != player.Id {
		return nil, ErrUnauthorizedBank
	}

	return s.pendingTransactionRepo.FindBySourceBankID(ctx, bankID)
}
