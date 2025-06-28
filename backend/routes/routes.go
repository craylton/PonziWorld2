package routes

import (
	"net/http"
	"ponziworld/backend/handlers"
	"ponziworld/backend/middleware"
)

func RegisterRoutes(mux *http.ServeMux) {	
	mux.HandleFunc("POST /api/user", handlers.CreateUserHandler)
	mux.HandleFunc("/api/bank", middleware.JWTMiddleware(handlers.GetBankHandler))
	mux.HandleFunc("/api/login", handlers.LoginHandler)
	mux.HandleFunc(
		"/api/performanceHistory/ownbank/{bankId}",
		middleware.JWTMiddleware(handlers.GetPerformanceHistoryHandler),
	)
}
