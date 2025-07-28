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

func (s *GameService) AdvanceToNextDay(ctx context.Context) (int, error) {
	// Process all pending transactions before advancing to the next day
	err := s.ProcessPendingTransactions(ctx)
	if err != nil {
		return 0, err
	}

	// Try to increment the day
	newDay, err := s.incrementDay(ctx)
	if err != nil {
		return 0, err
	}
	return newDay, nil
}

func (s *GameService) incrementDay(ctx context.Context) (int, error) {
	// Increment the day in the game state
	newDay, err := s.gameRepo.IncrementDay(ctx)
	if err == nil {
		return newDay, nil
	}
	// Some unexpected error occurred
	if err != mongo.ErrNoDocuments {
		return 0, err
	}

	// No game state exists, create it starting at day 1
	err = s.gameRepo.CreateInitialGame(ctx, 1)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

// ProcessPendingTransactions converts all pending transactions into actual investments
func (s *GameService) ProcessPendingTransactions(ctx context.Context) error {
	// Get all pending transactions
	pendingTransactions, err := s.pendingTransactionRepo.FindAll(ctx)
	if err != nil {
		return err
	}

	if len(pendingTransactions) == 0 {
		return nil
	}

	// Get cash asset type once for all transactions
	cashAssetType, err := s.assetTypeRepo.FindByName(ctx, "Cash")
	if err != nil {
		return err
	}

	// Process each pending transaction
	for _, pendingTransaction := range pendingTransactions {
		// Validate the pending transaction
		err := s.validatePendingTransaction(ctx, &pendingTransaction, cashAssetType.Id)
		if err != nil {
			// Todo: log error
			continue
		}

		// Process the pending transaction
		err = s.processSinglePendingTransaction(ctx, &pendingTransaction)
		if err != nil {
			// Todo: log error
			continue
		}

		// Remove the processed pending transaction
		err = s.pendingTransactionRepo.Delete(ctx, pendingTransaction.Id)
		if err != nil {
			// Todo: log error
			continue
		}
	}

	return nil
}

func (s *GameService) validatePendingTransaction(
	ctx context.Context,
	pendingTransaction *models.PendingTransactionResponse,
	cashAssetTypeId primitive.ObjectID,
) error {
	// Cash is handled slightly differently as we allow negative cash
	if pendingTransaction.TargetAssetId == cashAssetTypeId {
		err := s.validateCashPendingTransaction(ctx, pendingTransaction)
		if err != nil {
			return err
		}
	} else {
		err := s.validateNonCashPendingTransaction(ctx, pendingTransaction)
		if err != nil {
			return err
		}
	}

	return nil
}

// processSinglePendingTransaction processes a single pending transaction
func (s *GameService) processSinglePendingTransaction(
	ctx context.Context,
	pendingTransaction *models.PendingTransactionResponse,
) error {
	// Find existing investment, if any
	existingInvestment, err := s.investmentRepo.FindBySourceIdAndTargetId(
		ctx,
		pendingTransaction.SourceBankId,
		pendingTransaction.TargetAssetId,
	)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	if existingInvestment != nil {
		// Update existing investment
		newAmount := existingInvestment.Amount + pendingTransaction.Amount

		if newAmount == 0 {
			// Remove investment entirely if amount becomes zero
			return s.investmentRepo.DeleteBySourceIdAndTargetId(
				ctx,
				pendingTransaction.SourceBankId,
				pendingTransaction.TargetAssetId,
			)
		} else {
			// Update with new amount (allow negative for cash)
			return s.investmentRepo.UpdateAmount(
				ctx,
				pendingTransaction.SourceBankId,
				pendingTransaction.TargetAssetId,
				newAmount,
			)
		}
	} else {
		// Create new investment (allow negative amounts for cash)
		if pendingTransaction.Amount == 0 {
			// Skip zero-amount transactions
			return nil
		}

		investment := &models.Investment{
			Id:            primitive.NewObjectID(),
			SourceBankId:  pendingTransaction.SourceBankId,
			TargetAssetId: pendingTransaction.TargetAssetId,
			Amount:        pendingTransaction.Amount,
		}

		return s.investmentRepo.Create(ctx, investment)
	}
}

func (s *GameService) validateCashPendingTransaction(
	ctx context.Context,
	pendingTransaction *models.PendingTransactionResponse,
) error {
	// Validate that the source bank still exists
	_, err := s.bankRepo.FindByID(ctx, pendingTransaction.SourceBankId)
	if err != nil {
		return errors.New("source bank no longer exists")
	}
	return nil
}

// validateNonCashPendingTransaction validates non-cash pending transactions
func (s *GameService) validateNonCashPendingTransaction(
	ctx context.Context,
	pendingTransaction *models.PendingTransactionResponse,
) error {
	// Validate that the source bank still exists
	_, err := s.bankRepo.FindByID(ctx, pendingTransaction.SourceBankId)
	if err != nil {
		return errors.New("source bank no longer exists")
	}

	// Validate that the target asset still exists
	assetExists, err := s.validateTargetAssetExists(ctx, pendingTransaction.TargetAssetId)
	if err != nil {
		return err
	}
	if !assetExists {
		return errors.New("target asset no longer exists")
	}

	// For sell transactions (negative amounts), validate that we have enough to sell
	if pendingTransaction.Amount < 0 {
		err := s.validateSufficientAssetForSale(
			ctx,
			pendingTransaction.SourceBankId,
			pendingTransaction.TargetAssetId,
			-pendingTransaction.Amount,
		)
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
