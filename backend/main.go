package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ponziworld/backend/config"
	"ponziworld/backend/db"
	"ponziworld/backend/middleware"
	"ponziworld/backend/routes"
)

func main() {
	// Initialize database connection
	client, ctx, cancel := db.ConnectDB()
	
	// Get database name from environment or use default
	databaseName := "ponziworld"
	
	// Ensure database indexes using the existing connection
	if err := db.EnsureAllIndexes(); err != nil {
		cancel()
		client.Disconnect(ctx)
		log.Fatalf("Failed to ensure database indexes: %v", err)
	}
	
	// Create handler dependencies
	deps := config.NewHandlerDependencies(client, ctx, cancel, databaseName)
	// Note: We don't defer deps.Close() here because it would close the connection
	// during server startup. The connection will be closed when the server shuts down.

	// Set up routes with dependencies
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, deps)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Backend listening on :%s (database: %s)\n", port, databaseName)
	log.Fatal(http.ListenAndServe(":"+port, middleware.CorsMiddleware(mux)))
}
