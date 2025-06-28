package handlers

import (
	"encoding/json"
	"net/http"

	"ponziworld/backend/db"
	"ponziworld/backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// GetPerformanceHistoryHandler handles GET /api/performanceHistory/ownbank/{bankID}
func GetPerformanceHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get username from the JWT token (set by middleware)
	username := r.Header.Get("X-Username")
	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Authentication required"})
		return
	}

	// Extract bank ID from URL path parameter
	bankIDStr := r.PathValue("bankID")
	if bankIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bank ID required"})
		return
	}

	bankID, err := primitive.ObjectIDFromHex(bankIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid bank ID"})
		return
	}

	client, ctx, cancel := db.ConnectDB()
	defer cancel()
	defer client.Disconnect(ctx)

	// First, find the user to get their ID
	usersCollection := client.Database("ponziworld").Collection("users")
	var user models.User
	err = usersCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	// Check if the bank exists and get bank details
	banksCollection := client.Database("ponziworld").Collection("banks")
	var bank models.Bank
	err = banksCollection.FindOne(ctx, bson.M{"_id": bankID}).Decode(&bank)
	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bank not found"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	// Check if the user owns the bank - reject if they don't
	if bank.UserID != user.ID {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: You can only view your own bank's performance history"})
		return
	}

	// Get the current day (for now, we'll use 0 as the current day - this can be made dynamic later)
	currentDay := 0
	startDay := currentDay - 29 // Get past 30 days (including today)

	// Get performance history
	claimedHistory, actualHistory, err := getPerformanceHistory(bankID, startDay, currentDay)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	response := models.PerformanceHistoryResponse{
		ClaimedHistory: convertToResponse(claimedHistory),
		ActualHistory:  convertToResponse(actualHistory),
	}

	json.NewEncoder(w).Encode(response)
}

// getPerformanceHistory retrieves existing performance history and ensures claimed history exists for 30 days
func getPerformanceHistory(bankID primitive.ObjectID, startDay, endDay int) ([]models.HistoricalPerformance, []models.HistoricalPerformance, error) {
	client, ctx, cancel := db.ConnectDB()
	defer cancel()
	defer client.Disconnect(ctx)

	historyCollection := client.Database("ponziworld").Collection("historicalPerformance")

	// Get all existing history for this bank in the date range
	filter := bson.M{
		"bankId": bankID,
		"day":    bson.M{"$gte": startDay, "$lte": endDay},
	}

	cursor, err := historyCollection.Find(ctx, filter, options.Find().SetSort(bson.M{"day": 1}))
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var allHistory []models.HistoricalPerformance
	if err = cursor.All(ctx, &allHistory); err != nil {
		return nil, nil, err
	}

	// Separate claimed and actual history
	claimedHistory := make([]models.HistoricalPerformance, 0)
	actualHistory := make([]models.HistoricalPerformance, 0)

	for _, entry := range allHistory {
		if entry.IsClaimed {
			claimedHistory = append(claimedHistory, entry)
		} else {
			actualHistory = append(actualHistory, entry)
		}
	}

	// Ensure we have claimed history for all 30 days - create missing entries
	claimedHistory, err = ensureClaimedHistory(bankID, startDay, endDay, claimedHistory)
	if err != nil {
		return nil, nil, err
	}

	return claimedHistory, actualHistory, nil
}

// ensureClaimedHistory creates missing claimed history entries if needed
func ensureClaimedHistory(bankID primitive.ObjectID, startDay, endDay int, existingClaimed []models.HistoricalPerformance) ([]models.HistoricalPerformance, error) {
	// Create map of existing claimed days for quick lookup
	existingClaimedDays := make(map[int]models.HistoricalPerformance)
	for _, entry := range existingClaimed {
		existingClaimedDays[entry.Day] = entry
	}

	var finalClaimedHistory []models.HistoricalPerformance
	var newEntries []interface{}

	// Ensure we have claimed history for all days in range
	for day := startDay; day <= endDay; day++ {
		if claimedEntry, exists := existingClaimedDays[day]; exists {
			finalClaimedHistory = append(finalClaimedHistory, claimedEntry)
		} else {
			// Create new claimed entry
			newClaimedEntry := models.HistoricalPerformance{
				ID:        primitive.NewObjectID(),
				Day:       day,
				BankID:    bankID,
				Value:     1000, // Dummy value
				IsClaimed: true,
			}
			newEntries = append(newEntries, newClaimedEntry)
			finalClaimedHistory = append(finalClaimedHistory, newClaimedEntry)
		}
	}

	// Insert new claimed entries if any
	if len(newEntries) > 0 {
		client, ctx, cancel := db.ConnectDB()
		defer cancel()
		defer client.Disconnect(ctx)

		historyCollection := client.Database("ponziworld").Collection("historicalPerformance")
		_, err := historyCollection.InsertMany(ctx, newEntries)
		if err != nil {
			return nil, err
		}
	}

	return finalClaimedHistory, nil
}

// convertToResponse converts HistoricalPerformance to DayValue response format
func convertToResponse(history []models.HistoricalPerformance) []models.DayValue {
	result := make([]models.DayValue, len(history))
	for i, entry := range history {
		result[i] = models.DayValue{
			Day:   entry.Day,
			Value: entry.Value,
		}
	}
	return result
}

// CreateInitialPerformanceHistory creates 30 days of initial claimed performance history for a new bank
func CreateInitialPerformanceHistory(bankID primitive.ObjectID, currentDay int) error {
	startDay := currentDay - 29
	_, err := ensureClaimedHistory(bankID, startDay, currentDay, []models.HistoricalPerformance{})
	return err
}
