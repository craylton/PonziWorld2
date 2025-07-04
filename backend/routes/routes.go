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
	
	// Register routes with dependency-injected handlers
	mux.HandleFunc("/api/newPlayer", handlers.CreateNewPlayerHandler) // TODO: Convert this handler
	mux.HandleFunc("/api/bank", middleware.JwtMiddleware(bankHandler.GetBank))
	mux.HandleFunc("/api/login", handlers.LoginHandler) // TODO: Convert this handler
	mux.HandleFunc("/api/currentDay", gameHandler.GetCurrentDay)
	mux.HandleFunc("/api/player", middleware.JwtMiddleware(handlers.GetPlayerHandler)) // TODO: Convert this handler
	mux.HandleFunc("/api/nextDay", middleware.AdminJwtMiddleware(gameHandler.AdvanceToNextDay))
	mux.HandleFunc(
		"/api/performanceHistory/ownbank/{bankId}",
		middleware.JwtMiddleware(handlers.GetPerformanceHistoryHandler), // TODO: Convert this handler
	)
}
