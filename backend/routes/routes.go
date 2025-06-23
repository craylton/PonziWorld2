package routes

import (
	"net/http"
	"ponziworld/backend/handlers"
	"ponziworld/backend/middleware"
)

func RegisterRoutes(mux *http.ServeMux) {
	// Public routes (no authentication required)
	mux.HandleFunc("POST /api/user", handlers.CreateUserHandler)
	mux.HandleFunc("POST /api/login", handlers.LoginHandler)

	// Protected routes (authentication required)
	mux.HandleFunc("GET /api/user", middleware.JWTMiddleware(handlers.GetUserHandler))
}
