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
	ErrCashNotTradable     = errors.New("cash cannot be bought or sold")
	ErrInsufficientFunds   = errors.New("insufficient cash balance for this purchase")
)

type PendingTransactionService struct {
	pendingTransactionRepo repositories.PendingTransactionRepository
	bankRepo               repositories.BankRepository
	assetTypeRepo          repositories.AssetTypeRepository
	playerRepo             repositories.PlayerRepository
	investmentRepo         repositories.InvestmentRepository
}

func NewPendingTransactionService(
	pendingTransactionRepo repositories.PendingTransactionRepository,
	bankRepo repositories.BankRepository,
	assetTypeRepo repositories.AssetTypeRepository,
	playerRepo repositories.PlayerRepository,
	investmentRepo repositories.InvestmentRepository,
) *PendingTransactionService {
	return &PendingTransactionService{
		pendingTransactionRepo: pendingTransactionRepo,
		bankRepo:               bankRepo,
		assetTypeRepo:          assetTypeRepo,
		playerRepo:             playerRepo,
		investmentRepo:         investmentRepo,
	}
}

func (s *PendingTransactionService) CreateBuyTransaction(
	ctx context.Context,
	sourceBankId,
	targetAssetId primitive.ObjectID,
	amount int64,
	username string,
) error {
	// First validate basic transaction requirements (bank ownership, asset existence, etc.)
	err := s.validateTransactionRequirements(ctx, sourceBankId, targetAssetId, username)
	if err != nil {
		return err
	}

	// Check cash balance considering existing pending transactions for this specific asset
	err = s.validateSufficientCashForTransaction(ctx, sourceBankId, targetAssetId, amount)
	if err != nil {
		return err
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

func (s *PendingTransactionService) validateSufficientCashForTransaction(
	ctx context.Context,
	sourceBankId,
	targetAssetId primitive.ObjectID,
	amount int64,
) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	// Get the Cash asset type
	cashAssetType, err := s.assetTypeRepo.FindByName(ctx, "Cash")
	if err != nil {
		return err
	}

	// Get current cash investment (this represents the cash balance)
	var currentCashBalance int64
	investment, err := s.investmentRepo.FindBySourceIdAndTargetId(ctx, sourceBankId, cashAssetType.Id)
	if err != nil {
		// No cash investment found, balance is 0
		currentCashBalance = 0
	} else {
		currentCashBalance = investment.Amount
	}

	// Get all pending transactions for this bank to calculate cash usage
	pendingTransactions, err := s.pendingTransactionRepo.FindBySourceBankID(ctx, sourceBankId)
	if err != nil {
		return err
	}

	// Calculate pending cash usage (all pending buy transactions consume cash)
	var pendingCashUsage int64
	var existingPendingForThisAsset int64
	for _, transaction := range pendingTransactions {
		if transaction.Amount > 0 { // Positive amount = buy transaction = cash usage
			if transaction.TargetAssetId == targetAssetId {
				// This is existing pending for the same asset we're buying
				existingPendingForThisAsset = transaction.Amount
			} else {
				// This is cash usage for other assets
				pendingCashUsage += transaction.Amount
			}
		}
	}

	// Calculate what the new pending amount would be for this asset
	newPendingForThisAsset := existingPendingForThisAsset + amount

	// Calculate total cash that would be used
	totalCashUsage := pendingCashUsage + newPendingForThisAsset

	// Check if there's enough cash
	if currentCashBalance < totalCashUsage {
		return ErrInsufficientFunds
	}

	return nil
}

func (s *PendingTransactionService) createTransaction(
	ctx context.Context,
	sourceBankId,
	targetAssetId primitive.ObjectID,
	amount int64,
	username string,
) error {
	// Validate basic transaction requirements
	err := s.validateTransactionRequirements(ctx, sourceBankId, targetAssetId, username)
	if err != nil {
		return err
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

// validateTransactionRequirements performs basic validation for all transactions
func (s *PendingTransactionService) validateTransactionRequirements(
	ctx context.Context,
	sourceBankId,
	targetAssetId primitive.ObjectID,
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

	// Check if target asset is cash - cash cannot be bought or sold
	cashAssetType, err := s.assetTypeRepo.FindByName(ctx, "Cash")
	if err == nil && cashAssetType.Id == targetAssetId {
		return ErrCashNotTradable
	}

	// Validate that bank is not investing in itself
	if sourceBankId == targetAssetId {
		return ErrSelfInvestment
	}

	return nil
}
