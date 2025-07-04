package routes

import (
	"net/http"
	"ponziworld/backend/config"
	"ponziworld/backend/handlers"
	"ponziworld/backend/middleware"
)

func RegisterRoutes(mux *http.ServeMux, deps *config.HandlerDependencies) {
	// Initialize handlers with dependencies
	bankHandler := handlers.NewBankHandler(deps)
	gameHandler := handlers.NewGameHandler(deps)
	playerHandler := handlers.NewPlayerHandler(deps)
	loginHandler := handlers.NewLoginHandler(deps)
	
	// Register routes
	mux.HandleFunc("/api/newPlayer", playerHandler.CreateNewPlayer)
	mux.HandleFunc("/api/bank", middleware.JwtMiddleware(bankHandler.GetBank))
	mux.HandleFunc("/api/login", loginHandler.LogIn)
	mux.HandleFunc("/api/currentDay", gameHandler.GetCurrentDay)
	mux.HandleFunc("/api/player", middleware.JwtMiddleware(playerHandler.GetPlayer))
	mux.HandleFunc("/api/nextDay", middleware.AdminJwtMiddleware(gameHandler.AdvanceToNextDay))
	mux.HandleFunc(
		"/api/performanceHistory/ownbank/{bankId}",
		middleware.JwtMiddleware(handlers.GetPerformanceHistoryHandler), // TODO: Convert this handler
	)
}
