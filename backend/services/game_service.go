package services

import (
	"context"
	"errors"
	"ponziworld/backend/models"
	"ponziworld/backend/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GameService struct {
	gameRepo               repositories.GameRepository
	pendingTransactionRepo repositories.PendingTransactionRepository
	investmentRepo         repositories.InvestmentRepository
	bankRepo               repositories.BankRepository
	assetTypeRepo          repositories.AssetTypeRepository
	playerRepo             repositories.PlayerRepository
}

func NewGameService(
	gameRepo repositories.GameRepository,
	pendingTransactionRepo repositories.PendingTransactionRepository,
	investmentRepo repositories.InvestmentRepository,
	bankRepo repositories.BankRepository,
	assetTypeRepo repositories.AssetTypeRepository,
	playerRepo repositories.PlayerRepository,
) *GameService {
	return &GameService{
		gameRepo:               gameRepo,
		pendingTransactionRepo: pendingTransactionRepo,
		investmentRepo:         investmentRepo,
		bankRepo:               bankRepo,
		assetTypeRepo:          assetTypeRepo,
		playerRepo:             playerRepo,
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
	// Process all pending transactions before advancing to the next day
	err := s.ProcessPendingTransactions(ctx)
	if err != nil {
		return 0, err
	}

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

// ProcessPendingTransactions converts all pending transactions into actual investments
func (s *GameService) ProcessPendingTransactions(ctx context.Context) error {
	// Get all pending transactions
	pendingTransactions, err := s.pendingTransactionRepo.FindAll(ctx)
	if err != nil {
		return err
	}

	if len(pendingTransactions) == 0 {
		return nil // Nothing to process
	}

	// Get cash asset type once for all transactions
	cashAssetType, err := s.assetTypeRepo.FindByName(ctx, "Cash")
	if err != nil {
		return err
	}

	// Pre-fetch all banks and asset types that appear in transactions to minimize DB calls
	bankIDs := make(map[primitive.ObjectID]bool)
	assetIDs := make(map[primitive.ObjectID]bool)

	for _, tx := range pendingTransactions {
		bankIDs[tx.SourceBankId] = true
		assetIDs[tx.TargetAssetId] = true
	}

	// Process each pending transaction
	for _, pendingTx := range pendingTransactions {
		err := s.processSinglePendingTransaction(ctx, &pendingTx, cashAssetType.Id)
		if err != nil {
			// Log error but continue processing other transactions
			// In a production system, you might want to handle this differently
			continue
		}

		// Remove the processed pending transaction
		err = s.pendingTransactionRepo.Delete(ctx, pendingTx.Id)
		if err != nil {
			// Log error but continue
			continue
		}
	}

	return nil
}

// processSinglePendingTransaction processes a single pending transaction
func (s *GameService) processSinglePendingTransaction(
	ctx context.Context,
	pendingTx *models.PendingTransactionResponse,
	cashAssetTypeId primitive.ObjectID,
) error {
	// For cash transactions, skip most validation - allow negative cash
	isCashTransaction := pendingTx.TargetAssetId == cashAssetTypeId

	if !isCashTransaction {
		// Only validate non-cash transactions
		err := s.validateNonCashPendingTransaction(ctx, pendingTx)
		if err != nil {
			return err
		}
	}

	// Find existing investment
	existingInvestment, err := s.investmentRepo.FindBySourceIdAndTargetId(
		ctx,
		pendingTx.SourceBankId,
		pendingTx.TargetAssetId,
	)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	if existingInvestment != nil {
		// Update existing investment
		newAmount := existingInvestment.Amount + pendingTx.Amount

		if newAmount == 0 {
			// Remove investment entirely if amount becomes zero
			return s.investmentRepo.DeleteBySourceIdAndTargetId(ctx, pendingTx.SourceBankId, pendingTx.TargetAssetId)
		} else {
			// Update with new amount (allow negative for cash)
			return s.investmentRepo.UpdateAmount(ctx, pendingTx.SourceBankId, pendingTx.TargetAssetId, newAmount)
		}
	} else {
		// Create new investment (allow negative amounts for cash)
		if pendingTx.Amount == 0 {
			// Skip zero-amount transactions
			return nil
		}

		investment := &models.Investment{
			Id:            primitive.NewObjectID(),
			SourceBankId:  pendingTx.SourceBankId,
			TargetAssetId: pendingTx.TargetAssetId,
			Amount:        pendingTx.Amount,
		}

		return s.investmentRepo.Create(ctx, investment)
	}
}

// validateNonCashPendingTransaction validates non-cash pending transactions
func (s *GameService) validateNonCashPendingTransaction(
	ctx context.Context,
	pendingTx *models.PendingTransactionResponse,
) error {
	// Validate that the source bank still exists
	_, err := s.bankRepo.FindByID(ctx, pendingTx.SourceBankId)
	if err != nil {
		return errors.New("source bank no longer exists")
	}

	// Validate that the target asset still exists
	assetExists, err := s.validateTargetAssetExists(ctx, pendingTx.TargetAssetId)
	if err != nil {
		return err
	}
	if !assetExists {
		return errors.New("target asset no longer exists")
	}

	// For sell transactions (negative amounts), validate that we have enough to sell
	if pendingTx.Amount < 0 {
		err := s.validateSufficientAssetForSale(ctx, pendingTx.SourceBankId, pendingTx.TargetAssetId, -pendingTx.Amount)
		if err != nil {
			return err
		}
	}

	return nil
}

// validateTargetAssetExists checks if the target asset (asset type or bank) exists
func (s *GameService) validateTargetAssetExists(ctx context.Context, targetAssetId primitive.ObjectID) (bool, error) {
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

// validateSufficientAssetForSale checks if there's enough of an asset to sell
func (s *GameService) validateSufficientAssetForSale(
	ctx context.Context,
	sourceBankId,
	targetAssetId primitive.ObjectID,
	amountToSell int64,
) error {
	// Get current investment amount
	investment, err := s.investmentRepo.FindBySourceIdAndTargetId(ctx, sourceBankId, targetAssetId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("insufficient asset balance for sale - no investment found")
		}
		return err
	}

	if investment.Amount < amountToSell {
		return errors.New("insufficient asset balance for sale")
	}

	return nil
}
