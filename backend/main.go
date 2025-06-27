package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ponziworld/backend/db"
	"ponziworld/backend/middleware"
	"ponziworld/backend/routes"
)

func main() {
	if err := db.EnsureAllIndexes(); err != nil {
		log.Fatalf("Failed to ensure database indexes: %v", err)
	}

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Backend listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, middleware.CorsMiddleware(mux)))
}
