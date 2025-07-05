package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"ponziworld/backend/auth"
	"ponziworld/backend/config"
	"ponziworld/backend/services"
)

// LoginHandler handles login-related requests
type LoginHandler struct {
	deps *config.Container
}

// NewLoginHandler creates a new LoginHandler
func NewLoginHandler(deps *config.Container) *LoginHandler {
	return &LoginHandler{deps: deps}
}

// LoginHandler handles POST /api/login
func (h *LoginHandler) LogIn(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}
	if req.Username == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Username and password required"})
		return
	}

	// Use the injected service manager and create a new context for this request
	ctx := context.Background() // Create a fresh context for this request
	serviceManager := h.deps.ServiceManager

	// Attempt login
	_, err := serviceManager.Auth.Login(ctx, req.Username, req.Password)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid username or password"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	// Generate JWT
	token, err := auth.GenerateToken(req.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to generate token"})
		return
	}

	// Return token
	response := map[string]string{"token": token}
	json.NewEncoder(w).Encode(response)
}
