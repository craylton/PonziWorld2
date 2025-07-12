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
	performanceHistoryHandler := handlers.NewPerformanceHistoryHandler(container)
	assetTypeHandler := handlers.NewAssetTypeHandler(container)
	pendingTransactionHandler := handlers.NewPendingTransactionHandler(container)

	// Register routes
	mux.HandleFunc("/api/newPlayer", playerHandler.CreateNewPlayer)
	mux.HandleFunc("/api/bank", middleware.JwtMiddleware(bankHandler.GetBank))
	mux.HandleFunc("/api/login", loginHandler.LogIn)
	mux.HandleFunc("/api/currentDay", gameHandler.GetCurrentDay)
	mux.HandleFunc("/api/player", middleware.JwtMiddleware(playerHandler.GetPlayer))
	mux.HandleFunc("/api/assetTypes", middleware.JwtMiddleware(assetTypeHandler.GetAllAssetTypes))
	mux.HandleFunc("/api/buy", middleware.JwtMiddleware(pendingTransactionHandler.BuyAsset))
	mux.HandleFunc("/api/sell", middleware.JwtMiddleware(pendingTransactionHandler.SellAsset))
	mux.HandleFunc(
		"/api/nextDay",
		middleware.AdminJwtMiddleware(gameHandler.AdvanceToNextDay, container.ServiceContainer.Auth),
	)
	mux.HandleFunc(
		"/api/performanceHistory/ownbank/{bankId}",
		middleware.JwtMiddleware(performanceHistoryHandler.GetPerformanceHistory),
	)
}
