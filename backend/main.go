package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"ponziworld/backend/config"
	"ponziworld/backend/database"
	"ponziworld/backend/middleware"
	"ponziworld/backend/routes"
)

func main() {
	client, err := database.InitializeDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	databaseName := "ponziworld"

	container := config.NewContainer(client, databaseName)
	defer container.Close() // Ensure proper cleanup on exit

	if err := database.EnsureAllIndexes(container.DatabaseConfig); err != nil {
		log.Fatalf("Failed to ensure database indexes: %v", err)
	}

	// Initialize asset types on startup
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := container.ServiceContainer.AssetType.EnsureAssetTypesExist(ctx); err != nil {
		log.Fatalf("Failed to initialize asset types: %v", err)
	}

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, container)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Backend listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, middleware.CorsMiddleware(mux)))
}
