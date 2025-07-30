package main

import (
	"net/http"
	"os"

	"ponziworld/backend/config"
	"ponziworld/backend/database"
	"ponziworld/backend/logging"
	"ponziworld/backend/middleware"
	"ponziworld/backend/routes"
)

func main() {
	logger := logging.NewLogger()

	client, err := database.InitializeDatabaseConnection()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}

	databaseName := "ponziworld"

	container := config.NewContainer(client, databaseName, logger)
	defer container.Close() // Ensure proper cleanup on exit

	err = database.EnsureDatabaseStructure(container.DatabaseConfig, container.ServiceContainer.AssetType)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to ensure database structure")
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
