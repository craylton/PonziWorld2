package routes

import (
	"net/http"
	"ponziworld/backend/handlers"
	"ponziworld/backend/middleware"
)

func RegisterRoutes(mux *http.ServeMux) {	
	mux.HandleFunc("/api/newPlayer", handlers.CreateNewPlayerHandler)
	mux.HandleFunc("/api/bank", middleware.JWTMiddleware(handlers.GetBankHandler))
	mux.HandleFunc("/api/login", handlers.LoginHandler)
	mux.HandleFunc("/api/player", middleware.JWTMiddleware(handlers.GetPlayerHandler))
	mux.HandleFunc("/api/nextDay", middleware.AdminMiddleware(handlers.NextDayHandler))
	mux.HandleFunc(
		"/api/performanceHistory/ownbank/{bankId}",
		middleware.JWTMiddleware(handlers.GetPerformanceHistoryHandler),
	)
}
