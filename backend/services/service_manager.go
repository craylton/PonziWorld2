package services

import (
	"ponziworld/backend/repositories"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ServiceManager holds all services and their dependencies
type ServiceManager struct {
	Auth        *AuthService
	Bank        *BankService
	Asset       *AssetService
	Performance *PerformanceService
	Player      *PlayerService
	Game        *GameService
}

// NewServiceManager creates and wires up all services
func NewServiceManager(db *mongo.Database) *ServiceManager {
	// Create repositories
	playerRepo := repositories.NewPlayerRepository(db)
	bankRepo := repositories.NewBankRepository(db)
	assetRepo := repositories.NewAssetRepository(db)
	historyRepo := repositories.NewHistoricalPerformanceRepository(db)
	gameRepo := repositories.NewGameRepository(db)

	// Create services
	authService := NewAuthService(playerRepo)
	bankService := NewBankService(playerRepo, bankRepo, assetRepo)
	assetService := NewAssetService(assetRepo)
	gameService := NewGameService(gameRepo)
	performanceService := NewPerformanceService(bankService, gameService, historyRepo)
	playerService := NewPlayerService(authService, bankService, assetService, performanceService)

	return &ServiceManager{
		Auth:        authService,
		Bank:        bankService,
		Asset:       assetService,
		Performance: performanceService,
		Player:      playerService,
		Game:        gameService,
	}
}
