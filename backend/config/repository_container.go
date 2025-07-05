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

func NewRepositoryContainer(db *mongo.Database) *RepositoryContainer {
	playerRepo := repositories.NewPlayerRepository(db)
	bankRepo := repositories.NewBankRepository(db)
	assetRepo := repositories.NewAssetRepository(db)
	historyRepo := repositories.NewHistoricalPerformanceRepository(db)
	gameRepo := repositories.NewGameRepository(db)

	return &RepositoryContainer{
		Bank:                  bankRepo,
		Asset:                 assetRepo,
		HistoricalPerformance: historyRepo,
		Player:                playerRepo,
		Game:                  gameRepo,
	}
}
