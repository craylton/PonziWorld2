package services

import (
	"context"
	"ponziworld/backend/models"
	"ponziworld/backend/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssetService struct {
	assetRepo                 repositories.AssetRepository
	assetTypeRepo             repositories.AssetTypeRepository
	bankRepo                  repositories.BankRepository
	bankService               *BankService
	pendingTransactionRepo    repositories.PendingTransactionRepository
	historicalPerformanceService *HistoricalPerformanceService
}

func NewAssetService(
	assetRepo repositories.AssetRepository,
	assetTypeRepo repositories.AssetTypeRepository,
	bankRepo repositories.BankRepository,
	bankService *BankService,
	pendingTransactionRepo repositories.PendingTransactionRepository,
	historicalPerformanceService *HistoricalPerformanceService,
) *AssetService {
	return &AssetService{
		assetRepo:                 assetRepo,
		assetTypeRepo:             assetTypeRepo,
		bankRepo:                  bankRepo,
		bankService:               bankService,
		pendingTransactionRepo:    pendingTransactionRepo,
		historicalPerformanceService: historicalPerformanceService,
	}
}

func (s *AssetService) CreateInitialAsset(
	ctx context.Context,
	bankID primitive.ObjectID,
	amount int64,
) (*models.Asset, error) {
	// Get the Cash asset type
	cashAssetType, err := s.assetTypeRepo.FindByName(ctx, "Cash")
	if err != nil {
		return nil, err
	}

	asset := &models.Asset{
		Id:          primitive.NewObjectID(),
		BankId:      bankID,
		Amount:      amount,
		AssetTypeId: cashAssetType.Id,
	}

	err = s.assetRepo.Create(ctx, asset)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func (s *AssetService) GetAssetDetails(ctx context.Context, username string, assetID primitive.ObjectID, bankID primitive.ObjectID) (*models.AssetDetailsResponse, error) {
	// 1. Look up the asset by asset ID - first try asset types, then banks
	var isBank bool
	_, err := s.assetTypeRepo.FindByID(ctx, assetID)
	if err != nil {
		// Asset type not found, try to find it as a bank
		_, err = s.bankRepo.FindByID(ctx, assetID)
		if err != nil {
			return nil, ErrAssetNotFound
		}
		isBank = true
	}

	// 2. Look up the bank by bank ID and validate ownership
	err = s.bankService.ValidateBankOwnership(ctx, username, bankID)
	if err != nil {
		return nil, err
	}

	// 3. Find out whether this bank has invested in this asset
	var investedAmount int64
	asset, err := s.assetRepo.FindByBankIDAndAssetTypeID(ctx, bankID, assetID)
	if err != nil {
		// If no asset found, invested amount is 0
		investedAmount = 0
	} else {
		investedAmount = asset.Amount
	}

	// 4. Find out whether this bank has any pending transactions for this asset
	pendingAmount, err := s.pendingTransactionRepo.SumPendingAmountByBankIDAndAssetID(ctx, bankID, assetID)
	if err != nil {
		return nil, err
	}

	// 5. Get the past 8 days of historical performance for this asset
	// For banks, we use the bank ID as the asset ID for historical performance
	historicalAssetID := assetID
	if isBank {
		// For bank assets, use the bank ID directly for historical performance
		historicalAssetID = assetID
	}
	
	historicalData, err := s.historicalPerformanceService.GetAssetHistoricalPerformance(ctx, historicalAssetID, 8)
	if err != nil {
		return nil, err
	}

	// 6. Return the response
	return &models.AssetDetailsResponse{
		InvestedAmount: investedAmount,
		PendingAmount:  pendingAmount,
		HistoricalData: historicalData,
	}, nil
}
