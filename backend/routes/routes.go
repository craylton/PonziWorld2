package routes

import (
	"net/http"
	"ponziworld/backend/handlers"
	"ponziworld/backend/middleware"
)

func RegisterRoutes(mux *http.ServeMux) {	
	mux.HandleFunc("/api/newPlayer", handlers.CreateNewPlayerHandler)
	mux.HandleFunc("/api/bank", middleware.JwtMiddleware(handlers.GetBankHandler))
	mux.HandleFunc("/api/login", handlers.LoginHandler)
	mux.HandleFunc("/api/player", middleware.JwtMiddleware(handlers.GetPlayerHandler))
	mux.HandleFunc("/api/currentDay", handlers.CurrentDayHandler)
	mux.HandleFunc("/api/nextDay", middleware.AdminJwtMiddleware(handlers.NextDayHandler))
	mux.HandleFunc(
		"/api/performanceHistory/ownbank/{bankId}",
		middleware.JwtMiddleware(handlers.GetPerformanceHistoryHandler),
	)
}
