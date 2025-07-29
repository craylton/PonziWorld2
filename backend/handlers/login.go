package handlers

import (
	"encoding/json"
	"net/http"

	"ponziworld/backend/auth"
	"ponziworld/backend/config"
	"ponziworld/backend/services"

	"github.com/rs/zerolog"
)

// LoginHandler handles login-related requests
type LoginHandler struct {
	authService *services.AuthService
	logger      zerolog.Logger
}

// NewLoginHandler creates a new LoginHandler
func NewLoginHandler(container *config.Container) *LoginHandler {
	return &LoginHandler{
		authService: container.ServiceContainer.Auth,
		logger:      container.Logger,
	}
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
		h.logger.Error().Err(err).Msg("Failed to decode login request body")
		return
	}
	if req.Username == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Username and password required"})
		return
	}

	// Use the request context for proper cancellation handling
	ctx := r.Context()

	// Attempt login
	_, err := h.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid username or password"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		h.logger.Error().Err(err).Str("username", req.Username).Msg("Database error during login")
		return
	}

	// Generate JWT
	token, err := auth.GenerateToken(req.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to generate token"})
		h.logger.Error().Err(err).Str("username", req.Username).Msg("Failed to generate JWT token")
		return
	}

	// Return token
	response := map[string]string{"token": token}
	json.NewEncoder(w).Encode(response)
}
