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
)

type PendingTransactionService struct {
	pendingTransactionRepo repositories.PendingTransactionRepository
	bankRepo               repositories.BankRepository
	assetTypeRepo          repositories.AssetTypeRepository
}

func NewPendingTransactionService(
	pendingTransactionRepo repositories.PendingTransactionRepository,
	bankRepo repositories.BankRepository,
	assetTypeRepo repositories.AssetTypeRepository,
) *PendingTransactionService {
	return &PendingTransactionService{
		pendingTransactionRepo: pendingTransactionRepo,
		bankRepo:               bankRepo,
		assetTypeRepo:          assetTypeRepo,
	}
}

func (s *PendingTransactionService) CreateTransaction(ctx context.Context, buyerBankId, assetId primitive.ObjectID, amount int64) error {
	// Validate amount is not zero
	if amount == 0 {
		return ErrInvalidAmount
	}

	// Validate buyer bank exists
	_, err := s.bankRepo.FindByID(ctx, buyerBankId)
	if err != nil {
		return ErrInvalidBankID
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

	// Create the pending transaction
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
