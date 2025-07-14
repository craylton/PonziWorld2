package config

import (
	"ponziworld/backend/services"
)

type ServiceContainer struct {
	Auth               *services.AuthService
	Bank               *services.BankService
	Asset              *services.AssetService
	AssetType          *services.AssetTypeService
	Performance        *services.PerformanceService
	Player             *services.PlayerService
	Game               *services.GameService
	PendingTransaction *services.PendingTransactionService
}

func NewServiceContainer(repositoryContainer *RepositoryContainer) *ServiceContainer {
	authService := services.NewAuthService(repositoryContainer.Player)
	bankService := services.NewBankService(
		repositoryContainer.Player,
		repositoryContainer.Bank,
		repositoryContainer.Asset,
		repositoryContainer.AssetType,
	)
	assetService := services.NewAssetService(repositoryContainer.Asset, repositoryContainer.AssetType)
	assetTypeService := services.NewAssetTypeService(repositoryContainer.AssetType)
	gameService := services.NewGameService(repositoryContainer.Game)
	performanceService := services.NewPerformanceService(
		bankService,
		gameService,
		repositoryContainer.HistoricalPerformance,
	)
	playerService := services.NewPlayerService(
		authService,
		bankService,
		assetService,
		performanceService,
	)
	pendingTransactionService := services.NewPendingTransactionService(
		repositoryContainer.PendingTransaction,
		repositoryContainer.Bank,
		repositoryContainer.AssetType,
		repositoryContainer.Player,
	)

	return &ServiceContainer{
		Auth:               authService,
		Bank:               bankService,
		Asset:              assetService,
		AssetType:          assetTypeService,
		Performance:        performanceService,
		Player:             playerService,
		Game:               gameService,
		PendingTransaction: pendingTransactionService,
	}
}
