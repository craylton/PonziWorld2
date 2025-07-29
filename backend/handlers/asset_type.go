package handlers

import (
	"encoding/json"
	"net/http"

	"ponziworld/backend/config"
	"ponziworld/backend/services"

	"github.com/rs/zerolog"
)

// AssetTypeHandler handles asset type-related requests

type AssetTypeHandler struct {
	assetTypeService *services.AssetTypeService
	logger           zerolog.Logger
}

// NewAssetTypeHandler creates a new AssetTypeHandler
func NewAssetTypeHandler(container *config.Container) *AssetTypeHandler {
	return &AssetTypeHandler{
		assetTypeService: container.ServiceContainer.AssetType,
		logger:           container.Logger,
	}
}

// GetAllAssetTypes handles GET /api/assetTypes
func (h *AssetTypeHandler) GetAllAssetTypes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		h.logger.Error().Msg("Invalid method for GetAllAssetTypes")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	assetTypes, err := h.assetTypeService.GetAllAssetTypes(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve asset types"})
		h.logger.Error().Err(err).Msg("Failed to retrieve asset types")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(assetTypes)
}
