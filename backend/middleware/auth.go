package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"ponziworld/backend/auth"
	"ponziworld/backend/requestcontext"
	"ponziworld/backend/services"
)

// validateJwt extracts and validates JWT token from request, returns username
func validateJwt(w http.ResponseWriter, r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Authorization header required"})
		return "", false
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bearer token required"})
		return "", false
	}

	username, err := auth.ValidateToken(tokenString)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"})
		return "", false
	}

	return username, true
}

// JwtMiddleware validates JWT for protected routes
func JwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := validateJwt(w, r)
		if !ok {
			return
		}

		// Store username in context for business logic
		ctx := requestcontext.WithUsername(r.Context(), username)
		next(w, r.WithContext(ctx))
	}
}

// AdminJwtMiddleware validates that the user is an admin
func AdminJwtMiddleware(next http.HandlerFunc, authService *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := validateJwt(w, r)
		if !ok {
			return
		}

		// Use the request context for proper cancellation handling
		ctx := r.Context()

		player, err := authService.GetPlayerByUsername(ctx, username)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to verify admin status"})
			return
		}

		if !player.IsAdmin {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "Admin access required"})
			return
		}

		// Store username in context for downstream handlers
		ctx = requestcontext.WithUsername(ctx, username)
		next(w, r.WithContext(ctx))
	}
}
