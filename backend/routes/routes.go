package routes

import (
	"net/http"
	"ponziworld/backend/handlers"
	"ponziworld/backend/middleware"
)

func RegisterRoutes(mux *http.ServeMux) {	
	mux.HandleFunc("POST /api/newPlayer", handlers.CreateNewPlayerHandler)
	mux.HandleFunc("/api/bank", middleware.JWTMiddleware(handlers.GetBankHandler))
	mux.HandleFunc("/api/login", handlers.LoginHandler)
	mux.HandleFunc(
		"/api/performanceHistory/ownbank/{bankId}",
		middleware.JWTMiddleware(handlers.GetPerformanceHistoryHandler),
	)
}
