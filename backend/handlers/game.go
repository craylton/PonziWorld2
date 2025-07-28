package handlers

import (
	"encoding/json"
	"net/http"

	"ponziworld/backend/config"
	"ponziworld/backend/services"
)

// BankHandler handles bank-related requests
type GameHandler struct {
	gameService *services.GameService
}

// NewBankHandler creates a new BankHandler
func NewGameHandler(container *config.Container) *GameHandler {
	return &GameHandler{
		gameService: container.ServiceContainer.Game,
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
		http.Error(w, "Failed to get current day", http.StatusInternalServerError)
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
		http.Error(w, "Failed to increment day", http.StatusInternalServerError)
		return
	}

	response := map[string]int{"currentDay": newDay}
	json.NewEncoder(w).Encode(response)
}
