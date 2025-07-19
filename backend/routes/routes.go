package routes

import (
	"net/http"
	"ponziworld/backend/config"
	"ponziworld/backend/handlers"
	"ponziworld/backend/middleware"
)

func RegisterRoutes(mux *http.ServeMux, container *config.Container) {
	// Initialize handlers with dependencies
	bankHandler := handlers.NewBankHandler(container)
	gameHandler := handlers.NewGameHandler(container)
	playerHandler := handlers.NewPlayerHandler(container)
	loginHandler := handlers.NewLoginHandler(container)
	historicalPerformanceHandler := handlers.NewHistoricalPerformanceHandler(container)
	assetTypeHandler := handlers.NewAssetTypeHandler(container)
	assetHandler := handlers.NewInvestmentHandler(container)
	pendingTransactionHandler := handlers.NewPendingTransactionHandler(container)

	// Register routes
	mux.HandleFunc("/api/newPlayer", playerHandler.CreateNewPlayer)
	mux.HandleFunc("/api/banks", middleware.JwtMiddleware(bankHandler.HandleBanks))
	mux.HandleFunc("/api/login", loginHandler.LogIn)
	mux.HandleFunc("/api/currentDay", gameHandler.GetCurrentDay)
	mux.HandleFunc("/api/player", middleware.JwtMiddleware(playerHandler.GetPlayer))
	mux.HandleFunc("/api/assetTypes", middleware.JwtMiddleware(assetTypeHandler.GetAllAssetTypes))
	mux.HandleFunc(
		"/api/investment/{targetAssetId}/{sourceBankId}",
		middleware.JwtMiddleware(assetHandler.GetInvestmentDetails),
	)
	mux.HandleFunc("/api/buy", middleware.JwtMiddleware(pendingTransactionHandler.BuyAsset))
	mux.HandleFunc("/api/sell", middleware.JwtMiddleware(pendingTransactionHandler.SellAsset))
	mux.HandleFunc(
		"/api/pendingTransactions/{bankId}",
		middleware.JwtMiddleware(pendingTransactionHandler.GetPendingTransactions),
	)
	mux.HandleFunc(
		"/api/nextDay",
		middleware.AdminJwtMiddleware(gameHandler.AdvanceToNextDay, container.ServiceContainer.Auth),
	)
	mux.HandleFunc(
		"/api/historicalPerformance/ownbank/{bankId}",
		middleware.JwtMiddleware(historicalPerformanceHandler.GetHistoricalPerformance),
	)
}
