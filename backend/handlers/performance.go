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

// GetPerformanceHistoryHandler handles GET /api/performanceHistory/bank/{bankID}
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

	// Check if the user owns the bank
	isOwnBank := bank.UserID == user.ID

	// Get the current day (for now, we'll use 0 as the current day - this can be made dynamic later)
	currentDay := 0
	startDay := currentDay - 29 // Get past 30 days (including today)

	// Get all performance history in a single query
	claimedHistory, actualHistory, err := getOrCreatePerformanceHistory(bankID, startDay, currentDay, isOwnBank)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	response := models.PerformanceHistoryResponse{
		ClaimedHistory: convertToResponse(claimedHistory),
	}

	// If user owns the bank, also return actual history
	if isOwnBank {
		response.ActualHistory = convertToResponse(actualHistory)
	}

	json.NewEncoder(w).Encode(response)
}

// getOrCreatePerformanceHistory retrieves existing performance history or creates it if missing
func getOrCreatePerformanceHistory(bankID primitive.ObjectID, startDay, endDay int, includeActual bool) ([]models.HistoricalPerformance, []models.HistoricalPerformance, error) {
	client, ctx, cancel := db.ConnectDB()
	defer cancel()
	defer client.Disconnect(ctx)

	historyCollection := client.Database("ponziworld").Collection("historicalPerformance")

	// Build query to get both claimed and actual in one call if needed
	var filter bson.M
	if includeActual {
		filter = bson.M{
			"bankId": bankID,
			"day":    bson.M{"$gte": startDay, "$lte": endDay},
		}
	} else {
		filter = bson.M{
			"bankId":    bankID,
			"isClaimed": true,
			"day":       bson.M{"$gte": startDay, "$lte": endDay},
		}
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

	// Check if we need to create missing data and store it in database
	missingData := false
	expectedDays := endDay - startDay + 1
	
	if len(claimedHistory) < expectedDays || (includeActual && len(actualHistory) < expectedDays) {
		missingData = true
	}

	if missingData {
		// Create and store missing performance history
		claimedHistory, actualHistory, err = createAndStoreMissingHistory(bankID, startDay, endDay, claimedHistory, actualHistory, includeActual)
		if err != nil {
			return nil, nil, err
		}
	}

	return claimedHistory, actualHistory, nil
}

// createAndStoreMissingHistory creates missing performance history entries and stores them in the database
func createAndStoreMissingHistory(bankID primitive.ObjectID, startDay, endDay int, existingClaimed, existingActual []models.HistoricalPerformance, includeActual bool) ([]models.HistoricalPerformance, []models.HistoricalPerformance, error) {
	client, ctx, cancel := db.ConnectDB()
	defer cancel()
	defer client.Disconnect(ctx)

	historyCollection := client.Database("ponziworld").Collection("historicalPerformance")

	// Create maps of existing days for quick lookup
	existingClaimedDays := make(map[int]models.HistoricalPerformance)
	existingActualDays := make(map[int]models.HistoricalPerformance)

	for _, entry := range existingClaimed {
		existingClaimedDays[entry.Day] = entry
	}
	for _, entry := range existingActual {
		existingActualDays[entry.Day] = entry
	}

	// Prepare new entries to insert
	var newEntries []interface{}
	var finalClaimedHistory []models.HistoricalPerformance
	var finalActualHistory []models.HistoricalPerformance

	for day := startDay; day <= endDay; day++ {
		// Handle claimed history
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

		// Handle actual history if needed
		if includeActual {
			if actualEntry, exists := existingActualDays[day]; exists {
				finalActualHistory = append(finalActualHistory, actualEntry)
			} else {
				// Create new actual entry (same value as claimed for consistency)
				newActualEntry := models.HistoricalPerformance{
					ID:        primitive.NewObjectID(),
					Day:       day,
					BankID:    bankID,
					Value:     1000, // Same dummy value as claimed
					IsClaimed: false,
				}
				newEntries = append(newEntries, newActualEntry)
				finalActualHistory = append(finalActualHistory, newActualEntry)
			}
		}
	}

	// Insert new entries if any
	if len(newEntries) > 0 {
		_, err := historyCollection.InsertMany(ctx, newEntries)
		if err != nil {
			return nil, nil, err
		}
	}

	return finalClaimedHistory, finalActualHistory, nil
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

// CreateInitialPerformanceHistory creates 30 days of dummy performance history for a new bank
func CreateInitialPerformanceHistory(bankID primitive.ObjectID, currentDay int) error {
	startDay := currentDay - 29
	// Create empty slices to simulate no existing data, and include actual history
	_, _, err := createAndStoreMissingHistory(bankID, startDay, currentDay, []models.HistoricalPerformance{}, []models.HistoricalPerformance{}, true)
	return err
}
