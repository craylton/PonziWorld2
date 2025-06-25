package routes

import (
	"encoding/json"
	"net/http"
	"ponziworld/backend/handlers"
	"ponziworld/backend/middleware"
)

func RegisterRoutes(mux *http.ServeMux) {
	// User routes: create and fetch
	mux.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateUserHandler(w, r)
		case http.MethodGet:
			middleware.JWTMiddleware(handlers.GetUserHandler)(w, r)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		}
	})
	
	// Bank route: fetch bank details
	mux.HandleFunc("/api/bank", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			middleware.JWTMiddleware(handlers.GetBankHandler)(w, r)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		}
	})
	
	// Login route
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.LoginHandler(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		}
	})
}
