package handlers

import (
	"encoding/json"
	"net/http"

	"ponziworld/backend/config"
	"ponziworld/backend/services"

	"github.com/rs/zerolog"
)

// BankHandler handles bank-related requests
type GameHandler struct {
	gameService *services.GameService
	logger      zerolog.Logger
}

// NewBankHandler creates a new BankHandler
func NewGameHandler(container *config.Container) *GameHandler {
	return &GameHandler{
		gameService: container.ServiceContainer.Game,
		logger:      container.Logger,
	}
}

func (h *GameHandler) GetCurrentDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Use the request context for proper cancellation handling
	ctx := r.Context()

	currentDay, err := h.gameService.GetCurrentDay(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get current day"})
		h.logger.Error().Err(err).Msg("Failed to get current day")
		return
	}

	response := map[string]int{"currentDay": currentDay}
	json.NewEncoder(w).Encode(response)
}

func (h *GameHandler) AdvanceToNextDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Use the request context for proper cancellation handling
	ctx := r.Context()

	newDay, err := h.gameService.AdvanceToNextDay(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to increment day"})
		h.logger.Error().Err(err).Msg("Failed to increment day")
		return
	}

	response := map[string]int{"currentDay": newDay}
	json.NewEncoder(w).Encode(response)
}
