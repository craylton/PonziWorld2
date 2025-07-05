package config

import (
	"ponziworld/backend/services"
)

type ServiceContainer struct {
	Auth        *services.AuthService
	Bank        *services.BankService
	Asset       *services.AssetService
	Performance *services.PerformanceService
	Player      *services.PlayerService
	Game        *services.GameService
}

func NewServiceContainer(repositoryContainer *RepositoryContainer) *ServiceContainer {
	authService := services.NewAuthService(repositoryContainer.Player)
	bankService := services.NewBankService(
		repositoryContainer.Player,
		repositoryContainer.Bank,
		repositoryContainer.Asset,
	)
	assetService := services.NewAssetService(repositoryContainer.Asset)
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

	return &ServiceContainer{
		Auth:        authService,
		Bank:        bankService,
		Asset:       assetService,
		Performance: performanceService,
		Player:      playerService,
		Game:        gameService,
	}
}
