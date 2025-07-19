package config

import (
	"ponziworld/backend/repositories"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RepositoryContainer struct {
	Bank                  repositories.BankRepository
	Investment            repositories.InvestmentRepository
	AssetType             repositories.AssetTypeRepository
	HistoricalPerformance repositories.HistoricalPerformanceRepository
	Player                repositories.PlayerRepository
	Game                  repositories.GameRepository
	PendingTransaction    repositories.PendingTransactionRepository
}

func NewRepositoryContainer(database *mongo.Database) *RepositoryContainer {
	playerRepo := repositories.NewPlayerRepository(database)
	bankRepo := repositories.NewBankRepository(database)
	investmentRepo := repositories.NewInvestmentRepository(database)
	assetTypeRepo := repositories.NewAssetTypeRepository(database)
	historyRepo := repositories.NewHistoricalPerformanceRepository(database)
	gameRepo := repositories.NewGameRepository(database)
	pendingTransactionRepo := repositories.NewPendingTransactionRepository(database)

	return &RepositoryContainer{
		Bank:                  bankRepo,
		Investment:            investmentRepo,
		AssetType:             assetTypeRepo,
		HistoricalPerformance: historyRepo,
		Player:                playerRepo,
		Game:                  gameRepo,
		PendingTransaction:    pendingTransactionRepo,
	}
}
