package config

import (
	"ponziworld/backend/repositories"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RepositoryContainer struct {
	Bank                  repositories.BankRepository
	Asset                 repositories.AssetRepository
	AssetType             repositories.AssetTypeRepository
	HistoricalPerformance repositories.HistoricalPerformanceRepository
	Player                repositories.PlayerRepository
	Game                  repositories.GameRepository
	PendingTransaction    repositories.PendingTransactionRepository
}

func NewRepositoryContainer(database *mongo.Database) *RepositoryContainer {
	playerRepo := repositories.NewPlayerRepository(database)
	bankRepo := repositories.NewBankRepository(database)
	assetRepo := repositories.NewAssetRepository(database)
	assetTypeRepo := repositories.NewAssetTypeRepository(database)
	historyRepo := repositories.NewHistoricalPerformanceRepository(database)
	gameRepo := repositories.NewGameRepository(database)
	pendingTransactionRepo := repositories.NewPendingTransactionRepository(database)

	return &RepositoryContainer{
		Bank:                  bankRepo,
		Asset:                 assetRepo,
		AssetType:             assetTypeRepo,
		HistoricalPerformance: historyRepo,
		Player:                playerRepo,
		Game:                  gameRepo,
		PendingTransaction:    pendingTransactionRepo,
	}
}
