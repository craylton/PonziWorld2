package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	client, ctx, cancel := ConnectDB()
	defer cancel()
	defer client.Disconnect(ctx)

	if err := EnsureUserIndexes(client); err != nil {
		log.Fatalf("Failed to ensure user indexes: %v", err)
	}

	mux := http.NewServeMux()
	RegisterRoutes(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Backend listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsMiddleware(mux)))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
