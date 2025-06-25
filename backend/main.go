package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	db "ponziworld/backend/db"
	"ponziworld/backend/middleware"
	routes "ponziworld/backend/routes"
)

func main() {
	client, ctx, cancel := db.ConnectDB()
	defer cancel()
	defer client.Disconnect(ctx)

	if err := db.EnsureUserIndexes(client); err != nil {
		log.Fatalf("Failed to ensure user indexes: %v", err)
	}
	
	if err := db.EnsureBankIndexes(client); err != nil {
		log.Fatalf("Failed to ensure bank indexes: %v", err)
	}
	
	if err := db.EnsureAssetIndexes(client); err != nil {
		log.Fatalf("Failed to ensure asset indexes: %v", err)
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
