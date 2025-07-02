package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"ponziworld/backend/auth"
	"ponziworld/backend/db"
	"ponziworld/backend/services"
)

// validateJWT extracts and validates JWT token from request, returns username
func validateJWT(w http.ResponseWriter, r *http.Request) (string, bool) {
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

// JWTMiddleware validates JWT for protected routes
func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := validateJWT(w, r)
		if !ok {
			return
		}

		// Add username to request context for use in handlers
		r.Header.Set("X-Username", username)
		next(w, r)
	}
}

// AdminMiddleware validates that the user is an admin
func AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := validateJWT(w, r)
		if !ok {
			return
		}

		// Check if user is admin
		client, ctx, cancel := db.ConnectDB()
		defer cancel()
		defer client.Disconnect(ctx)

		serviceManager := services.NewServiceManager(client.Database("ponziworld"))
		player, err := serviceManager.Auth.GetPlayerByUsername(ctx, username)
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

		// Add username to request context for use in handlers
		r.Header.Set("X-Username", username)
		next(w, r)
	}
}
