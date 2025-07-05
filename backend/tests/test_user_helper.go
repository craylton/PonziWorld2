package tests

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"ponziworld/backend/auth"
	"ponziworld/backend/config"
	"ponziworld/backend/models"
)

// CreateAdminUserForTest creates an admin user for testing purposes
func CreateAdminUserForTest(container *config.Container, username, password, bankName string) (string, error) {
	ctx := context.Background()
	database := container.DatabaseConfig.Client.Database(container.DatabaseConfig.DatabaseName)
	
	// Create admin player manually (bypass the normal service to set IsAdmin = true)
	playersCollection := database.Collection("players")
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	
	playerID := primitive.NewObjectID()
	player := models.Player{
		Id:       playerID,
		Username: username,
		Password: string(hashedPassword),
		IsAdmin:  true, // Set as admin
	}
	
	_, err = playersCollection.InsertOne(ctx, player)
	if err != nil {
		return "", err
	}
	
	// Create bank for the admin user
	_, err = container.ServiceContainer.Bank.CreateBank(ctx, playerID, bankName, 1000)
	if err != nil {
		return "", err
	}
	
	// Generate JWT token for the admin user
	token, err := auth.GenerateToken(username)
	if err != nil {
		return "", err
	}
	
	return token, nil
}

// CreateRegularUserForTest creates a regular (non-admin) user for testing purposes
func CreateRegularUserForTest(container *config.Container, username, password, bankName string) (string, error) {
	ctx := context.Background()
	database := container.DatabaseConfig.Client.Database(container.DatabaseConfig.DatabaseName)
	
	// Create regular player manually (bypass the normal service to set IsAdmin = false)
	playersCollection := database.Collection("players")
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	
	playerID := primitive.NewObjectID()
	player := models.Player{
		Id:       playerID,
		Username: username,
		Password: string(hashedPassword),
		IsAdmin:  false, // Set as regular user
	}
	
	_, err = playersCollection.InsertOne(ctx, player)
	if err != nil {
		return "", err
	}
	
	// Create bank for the user
	_, err = container.ServiceContainer.Bank.CreateBank(ctx, playerID, bankName, 1000)
	if err != nil {
		return "", err
	}
	
	// Generate JWT token for the user
	token, err := auth.GenerateToken(username)
	if err != nil {
		return "", err
	}
	
	return token, nil
}
