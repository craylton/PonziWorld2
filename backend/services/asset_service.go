package services

import (
	"context"
	"ponziworld/backend/models"
	"ponziworld/backend/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssetService struct {
	assetRepo     repositories.AssetRepository
	assetTypeRepo repositories.AssetTypeRepository
}

func NewAssetService(
	assetRepo repositories.AssetRepository,
	assetTypeRepo repositories.AssetTypeRepository,
) *AssetService {
	return &AssetService{
		assetRepo:     assetRepo,
		assetTypeRepo: assetTypeRepo,
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
