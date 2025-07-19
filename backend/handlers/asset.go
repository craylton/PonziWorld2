package handlers

import (
	"encoding/json"
	"net/http"

	"ponziworld/backend/config"
	"ponziworld/backend/requestcontext"
	"ponziworld/backend/services"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssetHandler struct {
	assetService *services.AssetService
}

func NewAssetHandler(container *config.Container) *AssetHandler {
	return &AssetHandler{
		assetService: container.ServiceContainer.Asset,
	}
}

// GetAssetDetails handles GET /api/asset/{assetId}/{bankId}
func (h *AssetHandler) GetAssetDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get username from context
	username, ok := requestcontext.UsernameFromContext(r.Context())
	if !ok {
		http.Error(w, "username not found in context", http.StatusUnauthorized)
		return
	}

	// Get path parameters
	assetIdStr := r.PathValue("assetId")
	bankIdStr := r.PathValue("bankId")

	// Validate and convert assetId
	assetId, err := primitive.ObjectIDFromHex(assetIdStr)
	if err != nil {
		http.Error(w, "invalid asset ID", http.StatusBadRequest)
		return
	}

	// Validate and convert bankId
	bankId, err := primitive.ObjectIDFromHex(bankIdStr)
	if err != nil {
		http.Error(w, "invalid bank ID", http.StatusBadRequest)
		return
	}

	// Get asset details
	response, err := h.assetService.GetAssetDetails(r.Context(), username, assetId, bankId)
	if err != nil {
		switch err {
		case services.ErrUnauthorized:
			http.Error(w, "unauthorized access", http.StatusUnauthorized)
		case services.ErrAssetNotFound:
			http.Error(w, "asset not found", http.StatusNotFound)
		case services.ErrBankNotFound:
			http.Error(w, "bank not found", http.StatusNotFound)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
