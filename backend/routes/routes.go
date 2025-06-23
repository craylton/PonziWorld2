package routes

import (
	"net/http"
	"ponziworld/backend/handlers"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/user", handlers.CreateUserHandler)
	mux.HandleFunc("POST /api/login", handlers.LoginHandler)
	mux.HandleFunc("GET /api/user", handlers.GetUserHandler)
}
