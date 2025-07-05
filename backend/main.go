package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ponziworld/backend/config"
	"ponziworld/backend/database"
	"ponziworld/backend/middleware"
	"ponziworld/backend/routes"
)

func main() {
	// Initialize database connection
	client, err := database.InitializeDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Get database name from environment or use default
	databaseName := "ponziworld"

	// Create handler dependencies
	container := config.NewContainer(client, databaseName)
	defer container.Close() // Ensure proper cleanup on exit

	// Ensure database indexes using the existing connection
	if err := database.EnsureAllIndexes(container.DatabaseConfig); err != nil {
		log.Fatalf("Failed to ensure database indexes: %v", err)
	}

	// Set up routes with dependencies
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, container)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Backend listening on :%s (database: %s)\n", port, databaseName)
	log.Fatal(http.ListenAndServe(":"+port, middleware.CorsMiddleware(mux)))
}
