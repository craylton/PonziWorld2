package services

import (
	"context"
	"ponziworld/backend/models"
	"ponziworld/backend/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssetTypeService struct {
	assetTypeRepo repositories.AssetTypeRepository
}

func NewAssetTypeService(assetTypeRepo repositories.AssetTypeRepository) *AssetTypeService {
	return &AssetTypeService{
		assetTypeRepo: assetTypeRepo,
	}
}

func (s *AssetTypeService) GetAllAssetTypes(ctx context.Context) ([]models.AssetType, error) {
	return s.assetTypeRepo.FindAll(ctx)
}

func (s *AssetTypeService) EnsureAssetTypesExist(ctx context.Context) error {
	assetTypes := []string{"Cash", "HYSA", "Bonds", "Stocks", "Crypto"}
	
	for _, typeName := range assetTypes {
		assetType := &models.AssetType{
			Id:   primitive.NewObjectID(),
			Name: typeName,
		}
		
		err := s.assetTypeRepo.UpsertByName(ctx, assetType)
		if err != nil {
			return err
		}
	}
	
	return nil
}
