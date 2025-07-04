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
	
	// Register routes with dependency-injected handlers
	mux.HandleFunc("/api/newPlayer", handlers.CreateNewPlayerHandler)
	mux.HandleFunc("/api/bank", middleware.JwtMiddleware(bankHandler.GetBank))
	
	// Add back the login handler (still uses old pattern but needed for authentication)
	mux.HandleFunc("/api/login", handlers.LoginHandler)
	mux.HandleFunc("/api/currentDay", handlers.CurrentDayHandler) // TODO: Convert this handler
	mux.HandleFunc("/api/player", middleware.JwtMiddleware(handlers.GetPlayerHandler)) // TODO: Convert this handler
	mux.HandleFunc("/api/nextDay", middleware.AdminJwtMiddleware(handlers.NextDayHandler)) // TODO: Convert this handler
	mux.HandleFunc(
		"/api/performanceHistory/ownbank/{bankId}",
		middleware.JwtMiddleware(handlers.GetPerformanceHistoryHandler), // TODO: Convert this handler
	)
}
