package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"ponziworld/backend/config"
	"ponziworld/backend/services"
)

// AssetTypeHandler handles asset type-related requests
type AssetTypeHandler struct {
	assetTypeService *services.AssetTypeService
}

// NewAssetTypeHandler creates a new AssetTypeHandler
func NewAssetTypeHandler(container *config.Container) *AssetTypeHandler {
	return &AssetTypeHandler{
		assetTypeService: container.ServiceContainer.AssetType,
	}
}

// GetAllAssetTypes handles GET /api/assetTypes
func (h *AssetTypeHandler) GetAllAssetTypes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Use the request context for proper cancellation handling
	ctx := r.Context()

	// Get all asset types
	assetTypes, err := h.assetTypeService.GetAllAssetTypes(ctx)
	if err != nil {
		log.Printf("Error getting asset types: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve asset types"})
		return
	}

	// Return the asset types
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(assetTypes)
}
