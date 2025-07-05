package config

import (
	"ponziworld/backend/repositories"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RepositoryContainer struct {
	Bank                  repositories.BankRepository
	Asset                 repositories.AssetRepository
	HistoricalPerformance repositories.HistoricalPerformanceRepository
	Player                repositories.PlayerRepository
	Game                  repositories.GameRepository
}

func NewRepositoryContainer(database *mongo.Database) *RepositoryContainer {
	playerRepo := repositories.NewPlayerRepository(database)
	bankRepo := repositories.NewBankRepository(database)
	assetRepo := repositories.NewAssetRepository(database)
	historyRepo := repositories.NewHistoricalPerformanceRepository(database)
	gameRepo := repositories.NewGameRepository(database)

	return &RepositoryContainer{
		Bank:                  bankRepo,
		Asset:                 assetRepo,
		HistoricalPerformance: historyRepo,
		Player:                playerRepo,
		Game:                  gameRepo,
	}
}
