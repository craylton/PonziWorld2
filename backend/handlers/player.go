package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"ponziworld/backend/config"
	"ponziworld/backend/requestcontext"
	"ponziworld/backend/services"

	"github.com/rs/zerolog"
)

type PlayerHandler struct {
	authService   *services.AuthService
	playerService *services.PlayerService
	logger        zerolog.Logger
}

func NewPlayerHandler(container *config.Container) *PlayerHandler {
	return &PlayerHandler{
		authService:   container.ServiceContainer.Auth,
		playerService: container.ServiceContainer.Player,
		logger:        container.Logger,
	}
}

// CreateNewPlayerHandler handles POST /api/newPlayer
func (h *PlayerHandler) CreateNewPlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		h.logger.Error().Msg("Invalid method for CreateNewPlayer")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		BankName string `json:"bankName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		h.logger.Error().Err(err).Msg("Failed to decode request body for CreateNewPlayer")
		return
	}

	// Trim whitespace and validate
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	req.BankName = strings.TrimSpace(req.BankName)
	if req.Username == "" || req.Password == "" || req.BankName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Username, password, and bank name required"})
		h.logger.Error().Msg("Username, password, or bank name is empty")
		return
	}

	ctx := r.Context()

	// Create new player with bank and initial assets
	err := h.playerService.CreateNewPlayer(ctx, req.Username, req.Password, req.BankName)
	if err != nil {
		if err == services.ErrUsernameExists {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Username already exists"})
			h.logger.Error().Err(err).Str("username", req.Username).Msg("Username already exists")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create player"})
		h.logger.Error().Err(err).Msg("Failed to create player")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetPlayerHandler handles GET /api/player - returns current player info
func (h *PlayerHandler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		h.logger.Error().Msg("Invalid method for GetPlayer")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Get username from context (set by JwtMiddleware)
	ctx := r.Context()
	username, ok := requestcontext.UsernameFromContext(ctx)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Username not found in token"})
		h.logger.Error().Msg("Username not found in context")
		return
	}

	player, err := h.authService.GetPlayerByUsername(ctx, username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch player data"})
		h.logger.Error().Err(err).Msg("Failed to fetch player data")
		return
	}

	json.NewEncoder(w).Encode(player)
}
