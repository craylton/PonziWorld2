package services

import (
	"context"
	"errors"
	"ponziworld/backend/models"
	"ponziworld/backend/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidAssetID     = errors.New("invalid asset ID")
	ErrInvalidAmount      = errors.New("amount must not be zero")
	ErrAssetNotFound      = errors.New("asset not found")
	ErrInvalidBankID      = errors.New("invalid bank ID")
	ErrTransactionExists  = errors.New("pending transaction already exists")
	ErrSelfInvestment     = errors.New("bank cannot invest in itself")
	ErrUnauthorizedBank   = errors.New("bank is not owned by the current player")
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

func (s *PendingTransactionService) CreateTransaction(ctx context.Context, buyerBankId, assetId primitive.ObjectID, amount int64, username string) error {
	// Validate amount is not zero
	if amount == 0 {
		return ErrInvalidAmount
	}

	// Validate buyer bank exists and is owned by the current player
	buyerBank, err := s.bankRepo.FindByID(ctx, buyerBankId)
	if err != nil {
		return ErrInvalidBankID
	}

	// Validate bank ownership
	player, err := s.playerRepo.FindByUsername(ctx, username)
	if err != nil {
		return ErrInvalidBankID // Player not found
	}

	if buyerBank.PlayerId != player.Id {
		return ErrUnauthorizedBank
	}

	// Validate asset exists (check both asset types and banks since banks are also assets)
	assetExists, err := s.validateAssetExists(ctx, assetId)
	if err != nil {
		return err
	}
	if !assetExists {
		return ErrAssetNotFound
	}

	// Validate that bank is not investing in itself
	if buyerBankId == assetId {
		return ErrSelfInvestment
	}

	// Check if there's an existing pending transaction for this bank-asset combination
	existingTransactions, err := s.pendingTransactionRepo.FindByBuyerBankIDAndAssetID(ctx, buyerBankId, assetId)
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
	transaction := &models.PendingTransaction{
		BuyerBankId: buyerBankId,
		AssetId:     assetId,
		Amount:      amount,
	}

	return s.pendingTransactionRepo.Create(ctx, transaction)
}

func (s *PendingTransactionService) validateAssetExists(ctx context.Context, assetId primitive.ObjectID) (bool, error) {
	// First check if it's an asset type
	_, err := s.assetTypeRepo.FindByID(ctx, assetId)
	if err == nil {
		return true, nil
	}

	// Then check if it's a bank (since banks are also assets)
	_, err = s.bankRepo.FindByID(ctx, assetId)
	if err == nil {
		return true, nil
	}

	return false, nil
}

func (s *PendingTransactionService) GetTransactionsByBuyerBankID(ctx context.Context, buyerBankID primitive.ObjectID) ([]models.PendingTransaction, error) {
	return s.pendingTransactionRepo.FindByBuyerBankID(ctx, buyerBankID)
}

func (s *PendingTransactionService) GetTransactionsByAssetID(ctx context.Context, assetID primitive.ObjectID) ([]models.PendingTransaction, error) {
	return s.pendingTransactionRepo.FindByAssetID(ctx, assetID)
}
