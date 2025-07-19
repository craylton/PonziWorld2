package config

import (
	"ponziworld/backend/services"
)

type ServiceContainer struct {
	Auth                  *services.AuthService
	Bank                  *services.BankService
	Asset                 *services.AssetService
	AssetType             *services.AssetTypeService
	HistoricalPerformance *services.HistoricalPerformanceService
	Player                *services.PlayerService
	Game                  *services.GameService
	PendingTransaction    *services.PendingTransactionService
}

func NewServiceContainer(repositoryContainer *RepositoryContainer) *ServiceContainer {
	authService := services.NewAuthService(repositoryContainer.Player)
	bankService := services.NewBankService(
		repositoryContainer.Player,
		repositoryContainer.Bank,
		repositoryContainer.Asset,
		repositoryContainer.AssetType,
		repositoryContainer.PendingTransaction,
	)
	gameService := services.NewGameService(repositoryContainer.Game)
	historicalPerformanceService := services.NewHistoricalPerformanceService(
		bankService,
		gameService,
		repositoryContainer.HistoricalPerformance,
	)
	assetService := services.NewAssetService(
		repositoryContainer.Asset,
		repositoryContainer.AssetType,
		repositoryContainer.Bank,
		bankService,
		repositoryContainer.PendingTransaction,
		historicalPerformanceService,
	)
	assetTypeService := services.NewAssetTypeService(repositoryContainer.AssetType)
	playerService := services.NewPlayerService(
		authService,
		bankService,
		assetService,
		historicalPerformanceService,
	)
	pendingTransactionService := services.NewPendingTransactionService(
		repositoryContainer.PendingTransaction,
		repositoryContainer.Bank,
		repositoryContainer.AssetType,
		repositoryContainer.Player,
	)

	return &ServiceContainer{
		Auth:                  authService,
		Bank:                  bankService,
		Asset:                 assetService,
		AssetType:             assetTypeService,
		HistoricalPerformance: historicalPerformanceService,
		Player:                playerService,
		Game:                  gameService,
		PendingTransaction:    pendingTransactionService,
	}
}
