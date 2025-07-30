package routes

import (
	"net/http"
	"ponziworld/backend/config"
	"ponziworld/backend/handlers"
	"ponziworld/backend/middleware"
)

func RegisterRoutes(mux *http.ServeMux, container *config.Container) {
	bankHandler := handlers.NewBankHandler(container)
	gameHandler := handlers.NewGameHandler(container)
	playerHandler := handlers.NewPlayerHandler(container)
	loginHandler := handlers.NewLoginHandler(container)
	historicalPerformanceHandler := handlers.NewHistoricalPerformanceHandler(container)
	assetTypeHandler := handlers.NewAssetTypeHandler(container)
	assetHandler := handlers.NewInvestmentHandler(container)
	pendingTransactionHandler := handlers.NewPendingTransactionHandler(container)

	// Register routes
	mux.HandleFunc(
		"/api/newPlayer", 
		playerHandler.CreateNewPlayer,
	)
	mux.HandleFunc(
		"/api/banks",
		 middleware.JwtMiddleware(bankHandler.HandleBanks, container.Logger),
	)
	mux.HandleFunc(
		"/api/login",
		loginHandler.LogIn,
	)
	mux.HandleFunc(
		"/api/currentDay",
		gameHandler.GetCurrentDay,
	)
	mux.HandleFunc(
		"/api/player",
		middleware.JwtMiddleware(playerHandler.GetPlayer, container.Logger),
	)
	mux.HandleFunc(
		"/api/assetTypes",
		middleware.JwtMiddleware(assetTypeHandler.GetAllAssetTypes, container.Logger),
	)
	mux.HandleFunc(
		"/api/investment/{targetAssetId}/{sourceBankId}",
		middleware.JwtMiddleware(assetHandler.GetInvestmentDetails, container.Logger),
	)
	mux.HandleFunc(
		"/api/buy",
		middleware.JwtMiddleware(pendingTransactionHandler.BuyAsset, container.Logger),
	)
	mux.HandleFunc(
		"/api/sell",
		middleware.JwtMiddleware(pendingTransactionHandler.SellAsset, container.Logger),
	)
	mux.HandleFunc(
		"/api/pendingTransactions/{bankId}",
		middleware.JwtMiddleware(pendingTransactionHandler.GetPendingTransactions, container.Logger),
	)
	mux.HandleFunc(
		"/api/nextDay",
		middleware.AdminJwtMiddleware(
			gameHandler.AdvanceToNextDay,
			container.ServiceContainer.Auth,
			container.Logger,
		),
	)
	mux.HandleFunc(
		"/api/historicalPerformance/ownbank/{bankId}",
		middleware.JwtMiddleware(historicalPerformanceHandler.GetHistoricalPerformance, container.Logger),
	)
	mux.HandleFunc(
		"/api/historicalPerformance/asset/{targetAssetId}/{sourceBankId}",
		middleware.JwtMiddleware(historicalPerformanceHandler.GetAssetHistoricalPerformance, container.Logger),
	)
}
