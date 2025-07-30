package config

import (
	"ponziworld/backend/services"

	"github.com/rs/zerolog"
)

type ServiceContainer struct {
	Auth                  *services.AuthService
	Bank                  *services.BankService
	Investment            *services.InvestmentService
	AssetType             *services.AssetTypeService
	HistoricalPerformance *services.HistoricalPerformanceService
	Player                *services.PlayerService
	Game                  *services.GameService
	PendingTransaction    *services.PendingTransactionService
}

func NewServiceContainer(repositoryContainer *RepositoryContainer, logger zerolog.Logger) *ServiceContainer {
	authService := services.NewAuthService(repositoryContainer.Player)
	bankService := services.NewBankService(
		repositoryContainer.Player,
		repositoryContainer.Bank,
		repositoryContainer.Investment,
		repositoryContainer.AssetType,
		repositoryContainer.PendingTransaction,
	)
	gameService := services.NewGameService(
		repositoryContainer.Game,
		repositoryContainer.PendingTransaction,
		repositoryContainer.Investment,
		repositoryContainer.Bank,
		repositoryContainer.AssetType,
		repositoryContainer.Player,
		logger,
	)
	historicalPerformanceService := services.NewHistoricalPerformanceService(
		bankService,
		gameService,
		repositoryContainer.HistoricalPerformance,
	)
	investmentService := services.NewInvestmentService(
		repositoryContainer.Investment,
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
		investmentService,
		historicalPerformanceService,
	)
	pendingTransactionService := services.NewPendingTransactionService(
		repositoryContainer.PendingTransaction,
		repositoryContainer.Bank,
		repositoryContainer.AssetType,
		repositoryContainer.Player,
		repositoryContainer.Investment,
	)

	return &ServiceContainer{
		Auth:                  authService,
		Bank:                  bankService,
		Investment:            investmentService,
		AssetType:             assetTypeService,
		HistoricalPerformance: historicalPerformanceService,
		Player:                playerService,
		Game:                  gameService,
		PendingTransaction:    pendingTransactionService,
	}
}
