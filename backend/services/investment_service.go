package services

import (
	"context"
	"ponziworld/backend/models"
	"ponziworld/backend/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvestmentService struct {
	investmentRepo               repositories.InvestmentRepository
	assetTypeRepo                repositories.AssetTypeRepository
	bankRepo                     repositories.BankRepository
	bankService                  *BankService
	pendingTransactionRepo       repositories.PendingTransactionRepository
	historicalPerformanceService *HistoricalPerformanceService
}

func NewInvestmentService(
	investmentRepo repositories.InvestmentRepository,
	assetTypeRepo repositories.AssetTypeRepository,
	bankRepo repositories.BankRepository,
	bankService *BankService,
	pendingTransactionRepo repositories.PendingTransactionRepository,
	historicalPerformanceService *HistoricalPerformanceService,
) *InvestmentService {
	return &InvestmentService{
		investmentRepo:               investmentRepo,
		assetTypeRepo:                assetTypeRepo,
		bankRepo:                     bankRepo,
		bankService:                  bankService,
		pendingTransactionRepo:       pendingTransactionRepo,
		historicalPerformanceService: historicalPerformanceService,
	}
}

func (s *InvestmentService) CreateInitialInvestment(
	ctx context.Context,
	bankID primitive.ObjectID,
	amount int64,
) (*models.Investment, error) {
	// Get the Cash asset type
	cashAssetType, err := s.assetTypeRepo.FindByName(ctx, "Cash")
	if err != nil {
		return nil, err
	}

	investment := &models.Investment{
		Id:            primitive.NewObjectID(),
		SourceBankId:  bankID,
		Amount:        amount,
		TargetAssetId: cashAssetType.Id,
	}

	err = s.investmentRepo.Create(ctx, investment)
	if err != nil {
		return nil, err
	}

	return investment, nil
}

func (s *InvestmentService) GetInvestmentDetails(
	ctx context.Context,
	username string,
	investmentTargetId primitive.ObjectID,
	investmentSourceBankId primitive.ObjectID,
) (*models.InvestmentDetailsResponse, error) {
	// 1. Look up the investment target by ID - first try asset types, then banks
	assetName := ""
	assetType, err := s.assetTypeRepo.FindByID(ctx, investmentTargetId)
	if err != nil {
		// Asset type not found, try to find it as a bank
		bank, err := s.bankRepo.FindByID(ctx, investmentTargetId)
		if err != nil {
			return nil, ErrTargetAssetNotFound
		}
		assetName = bank.BankName
	} else {
		assetName = assetType.Name
	}

	// 2. Look up the bank by bank ID and validate ownership
	err = s.bankService.ValidateBankOwnership(ctx, username, investmentSourceBankId)
	if err != nil {
		return nil, err
	}

	// 3. Find out whether this bank has invested in this asset
	var investedAmount int64
	investment, err := s.investmentRepo.FindBySourceIdAndTargetId(ctx, investmentSourceBankId, investmentTargetId)
	if err != nil {
		// If no asset found, invested amount is 0
		investedAmount = 0
	} else {
		investedAmount = investment.Amount
	}

	// 4. Find out whether this bank has any pending transactions for this asset
	pendingAmount, err := s.pendingTransactionRepo.SumPendingAmountBySourceBankIdAndTargetAssetId(
		ctx,
		investmentSourceBankId,
		investmentTargetId,
	)
	if err != nil {
		return nil, err
	}

	historicalData, err := s.historicalPerformanceService.GetAssetHistoricalPerformance(ctx, investmentTargetId, 8)
	if err != nil {
		return nil, err
	}

	return &models.InvestmentDetailsResponse{
		TargetAssetId:  investmentTargetId.Hex(),
		Name:           assetName,
		InvestedAmount: investedAmount,
		PendingAmount:  pendingAmount,
		HistoricalData: historicalData,
	}, nil
}
