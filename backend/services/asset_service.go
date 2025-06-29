package services

import (
	"context"
	"ponziworld/backend/models"
	"ponziworld/backend/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssetService struct {
	assetRepo repositories.AssetRepository
}

func NewAssetService(assetRepo repositories.AssetRepository) *AssetService {
	return &AssetService{
		assetRepo: assetRepo,
	}
}

func (s *AssetService) CreateInitialAsset(ctx context.Context, bankID primitive.ObjectID, amount int64) (*models.Asset, error) {
	asset := &models.Asset{
		Id:        primitive.NewObjectID(),
		BankId:    bankID,
		Amount:    amount,
		AssetType: "Cash",
	}

	err := s.assetRepo.Create(ctx, asset)
	if err != nil {
		return nil, err
	}

	return asset, nil
}
