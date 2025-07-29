package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"ponziworld/backend/config"
	"ponziworld/backend/database"
	"ponziworld/backend/logging"
	"ponziworld/backend/middleware"
	"ponziworld/backend/routes"
)

func main() {
	// Initialize the logger first
	logger := logging.NewLogger()

	client, err := database.InitializeDatabaseConnection()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}

	databaseName := "ponziworld"

	container := config.NewContainer(client, databaseName, logger)
	defer container.Close() // Ensure proper cleanup on exit

	if err := database.EnsureAllIndexes(container.DatabaseConfig); err != nil {
		logger.Fatal().Err(err).Msg("Failed to ensure database indexes")
	}

	// Initialize asset types on startup
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := container.ServiceContainer.AssetType.EnsureAssetTypesExist(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize asset types")
	}

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, container)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	logger.Info().Str("port", port).Msg("Backend listening")
	
	if err := http.ListenAndServe(":"+port, middleware.CorsMiddleware(mux)); err != nil {
		logger.Fatal().Err(err).Msg("Server failed to start")
	}
}
