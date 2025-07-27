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

	// Get cash asset type ID
	cashAssetType, err := s.assetTypeRepo.FindByName(ctx, "Cash")
	if err != nil {
		return err
	}

	// Check cash balance considering existing pending transactions
	err = s.validateSufficientCashForTransaction(ctx, sourceBankId, amount)
	if err != nil {
		return err
	}

	// Create the asset purchase transaction (positive amount)
	err = s.createTransaction(ctx, sourceBankId, targetAssetId, amount, username)
	if err != nil {
		return err
	}

	// Create the cash transaction (negative amount to represent cash spent)
	return s.createTransaction(ctx, sourceBankId, cashAssetType.Id, -amount, username)
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

	// First validate transaction requirements (this will prevent selling cash directly)
	err := s.validateTransactionRequirements(ctx, sourceBankId, targetAssetId, username)
	if err != nil {
		return err
	}

	// Get cash asset type ID
	cashAssetType, err := s.assetTypeRepo.FindByName(ctx, "Cash")
	if err != nil {
		return err
	}

	// Create the asset sale transaction (negative amount)
	err = s.createTransaction(ctx, sourceBankId, targetAssetId, -amount, username)
	if err != nil {
		return err
	}

	// Create the cash transaction (positive amount to represent cash received)
	return s.createTransaction(ctx, sourceBankId, cashAssetType.Id, amount, username)
}

func (s *PendingTransactionService) validateSufficientCashForTransaction(
	ctx context.Context,
	sourceBankId primitive.ObjectID,
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

	// Get all pending cash transactions for this bank
	pendingCashTransactions, err := s.pendingTransactionRepo.FindBySourceBankIDAndTargetAssetID(ctx, sourceBankId, cashAssetType.Id)
	if err != nil {
		return err
	}

	// Calculate net pending cash changes
	var pendingCashChange int64
	for _, transaction := range pendingCashTransactions {
		pendingCashChange += transaction.Amount
	}

	// Calculate what the cash balance would be after all pending transactions and this new transaction
	projectedCashBalance := currentCashBalance + pendingCashChange - amount

	// Check if there's enough cash
	if projectedCashBalance < 0 {
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
	// Validate basic requirements - for internal use, we skip cash restriction
	// since cash transactions are created automatically by buy/sell operations
	
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

// validateTransactionRequirements performs validation for transactions (prevents direct cash trading)
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

	// Check if target asset is cash - cash cannot be bought or sold directly by users
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