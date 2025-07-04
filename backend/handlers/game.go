package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"ponziworld/backend/config"
)

// BankHandler handles bank-related requests
type GameHandler struct {
	deps *config.HandlerDependencies
}

// NewBankHandler creates a new BankHandler
func NewGameHandler(deps *config.HandlerDependencies) *GameHandler {
	return &GameHandler{deps: deps}
}

func (h *GameHandler) GetCurrentDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Use the injected service manager and create a new context for this request
	ctx := context.Background() // Create a fresh context for this request
	serviceManager := h.deps.ServiceManager

	currentDay, err := serviceManager.Game.GetCurrentDay(ctx)
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

	// Use the injected service manager and create a new context for this request
	ctx := context.Background() // Create a fresh context for this request
	serviceManager := h.deps.ServiceManager

	newDay, err := serviceManager.Game.NextDay(ctx)
	if err != nil {
		http.Error(w, "Failed to increment day", http.StatusInternalServerError)
		return
	}

	response := map[string]int{"currentDay": newDay}
	json.NewEncoder(w).Encode(response)
}
