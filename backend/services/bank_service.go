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
	ErrBankNotFound = errors.New("bank not found")
	ErrUnauthorized = errors.New("unauthorized access")
)

type BankService struct {
	playerRepo             repositories.PlayerRepository
	bankRepo               repositories.BankRepository
	assetRepo              repositories.InvestmentRepository
	assetTypeRepo          repositories.AssetTypeRepository
	pendingTransactionRepo repositories.PendingTransactionRepository
}

func NewBankService(
	playerRepo repositories.PlayerRepository,
	bankRepo repositories.BankRepository,
	assetRepo repositories.InvestmentRepository,
	assetTypeRepo repositories.AssetTypeRepository,
	pendingTransactionRepo repositories.PendingTransactionRepository,
) *BankService {
	return &BankService{
		playerRepo:             playerRepo,
		bankRepo:               bankRepo,
		assetRepo:              assetRepo,
		assetTypeRepo:          assetTypeRepo,
		pendingTransactionRepo: pendingTransactionRepo,
	}
}

func (s *BankService) GetAllBanksByUsername(ctx context.Context, username string) ([]models.BankResponse, error) {
	// Find the player
	player, err := s.playerRepo.FindByUsername(ctx, username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrPlayerNotFound
		}
		return nil, err
	}

	// Find all banks for this player
	banks, err := s.bankRepo.FindAllByPlayerID(ctx, player.Id)
	if err != nil {
		return nil, err
	}

	// Convert each bank to BankResponse
	bankResponses := make([]models.BankResponse, len(banks))
	for i, bank := range banks {
		// Get all assets for this bank
		assets, err := s.assetRepo.FindBySourceBankID(ctx, bank.Id)
		if err != nil {
			return nil, err
		}

		// Calculate actual capital
		actualCapital, err := s.assetRepo.CalculateActualCapital(ctx, bank.Id)
		if err != nil {
			return nil, err
		}

		// Get all asset types
		allAssetTypes, err := s.assetTypeRepo.FindAll(ctx)
		if err != nil {
			return nil, err
		}

		// Get all pending transactions for this bank
		pendingTransactions, err := s.pendingTransactionRepo.FindBySourceBankID(ctx, bank.Id)
		if err != nil {
			return nil, err
		}

		// Create maps for quick lookup
		investedAssetTypes := make(map[string]bool)
		for _, asset := range assets {
			investedAssetTypes[asset.TargetAssetId.Hex()] = true
		}

		pendingAssetTypes := make(map[string]bool)
		for _, transaction := range pendingTransactions {
			pendingAssetTypes[transaction.TargetAssetId.Hex()] = true
		}

		// Create available assets response
		availableAssets := make([]models.AvailableAssetResponse, len(allAssetTypes))
		for j, assetType := range allAssetTypes {
			assetTypeIdStr := assetType.Id.Hex()
			isInvestedOrPending := investedAssetTypes[assetTypeIdStr] || pendingAssetTypes[assetTypeIdStr]

			availableAssets[j] = models.AvailableAssetResponse{
				AssetTypeId:         assetTypeIdStr,
				AssetName:           assetType.Name,
				IsInvestedOrPending: isInvestedOrPending,
			}
		}

		// Create response
		bankResponses[i] = models.BankResponse{
			Id:              bank.Id.Hex(),
			BankName:        bank.BankName,
			ClaimedCapital:  bank.ClaimedCapital,
			ActualCapital:   actualCapital,
			AvailableAssets: availableAssets,
		}
	}

	return bankResponses, nil
}

func (s *BankService) CreateBankForUsername(
	ctx context.Context,
	username string,
	bankName string,
	claimedCapital int64,
) (*models.Bank, error) {
	// Find the player
	player, err := s.playerRepo.FindByUsername(ctx, username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrPlayerNotFound
		}
		return nil, err
	}

	// Create the bank
	bank := &models.Bank{
		Id:             primitive.NewObjectID(),
		PlayerId:       player.Id,
		BankName:       bankName,
		ClaimedCapital: claimedCapital,
	}

	err = s.bankRepo.Create(ctx, bank)
	if err != nil {
		return nil, err
	}

	// Create initial cash asset
	err = s.createInitialCashAsset(ctx, bank.Id, claimedCapital)
	if err != nil {
		return nil, err
	}

	return bank, nil
}

func (s *BankService) createInitialCashAsset(ctx context.Context, bankID primitive.ObjectID, amount int64) error {
	// Get the Cash asset type
	cashAssetType, err := s.assetTypeRepo.FindByName(ctx, "Cash")
	if err != nil {
		return err
	}

	asset := &models.Investment{
		Id:            primitive.NewObjectID(),
		SourceBankId:  bankID,
		Amount:        amount,
		TargetAssetId: cashAssetType.Id,
	}

	return s.assetRepo.Create(ctx, asset)
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
