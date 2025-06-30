package handlers

import (
	"encoding/json"
	"net/http"
	"ponziworld/backend/db"
	"ponziworld/backend/services"
)

func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	client, ctx, cancel := db.ConnectDB()
	defer cancel()
	defer client.Disconnect(ctx)

	serviceManager := services.NewServiceManager(client.Database("ponziworld"))

	newDay, err := serviceManager.Game.NextDay(ctx)
	if err != nil {
		http.Error(w, "Failed to increment day", http.StatusInternalServerError)
		return
	}

	response := map[string]int{"currentDay": newDay}
	json.NewEncoder(w).Encode(response)
}
