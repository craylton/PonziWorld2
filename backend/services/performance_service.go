package services

import (
	"context"
	"ponziworld/backend/models"
	"ponziworld/backend/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	PerformanceHistoryDays  = 30
	DefaultPerformanceValue = 1000
)

type PerformanceService struct {
	bankService *BankService
	gameService *GameService
	historyRepo repositories.HistoricalPerformanceRepository
}

func NewPerformanceService(
	bankService *BankService,
	gameService *GameService,
	historyRepo repositories.HistoricalPerformanceRepository,
) *PerformanceService {
	return &PerformanceService{
		bankService: bankService,
		gameService: gameService,
		historyRepo: historyRepo,
	}
}

func (s *PerformanceService) GetPerformanceHistory(ctx context.Context, username string, bankID primitive.ObjectID) (*models.PerformanceHistoryResponse, error) {
	// Validate ownership
	if err := s.bankService.ValidateBankOwnership(ctx, username, bankID); err != nil {
		return nil, err
	}

	currentDay, err := s.gameService.GetCurrentDay(ctx)
	if err != nil {
		return nil, err
	}
	startDay := currentDay - PerformanceHistoryDays

	// Get performance data
	claimedHistory, actualHistory, err := s.getPerformanceData(ctx, bankID, startDay, currentDay)
	if err != nil {
		return nil, err
	}

	return &models.PerformanceHistoryResponse{
		ClaimedHistory: convertToResponse(claimedHistory),
		ActualHistory:  convertToResponse(actualHistory),
	}, nil
}

func (s *PerformanceService) getPerformanceData(ctx context.Context, bankID primitive.ObjectID, startDay, endDay int) (
	[]models.HistoricalPerformance,
	[]models.HistoricalPerformance, error,
) {
	// Get all existing history for this bank in the date range
	allHistory, err := s.historyRepo.FindByBankIDAndDateRange(ctx, bankID, startDay, endDay)
	if err != nil {
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

	// Ensure we have claimed history for all days - create missing entries
	claimedHistory, err = s.ensureClaimedHistory(ctx, bankID, startDay, endDay, claimedHistory)
	if err != nil {
		return nil, nil, err
	}

	return claimedHistory, actualHistory, nil
}

func (s *PerformanceService) ensureClaimedHistory(
	ctx context.Context,
	bankID primitive.ObjectID,
	startDay,
	endDay int,
	existingClaimed []models.HistoricalPerformance,
) ([]models.HistoricalPerformance, error) {
	// Create map of existing claimed days for quick lookup
	existingClaimedDays := make(map[int]models.HistoricalPerformance)
	for _, entry := range existingClaimed {
		existingClaimedDays[entry.Day] = entry
	}

	var finalClaimedHistory []models.HistoricalPerformance
	var newEntries []models.HistoricalPerformance

	// Ensure we have claimed history for all days in range
	for day := startDay + 1; day <= endDay; day++ {
		if claimedEntry, exists := existingClaimedDays[day]; exists {
			finalClaimedHistory = append(finalClaimedHistory, claimedEntry)
		} else {
			// Create new claimed entry
			newClaimedEntry := models.HistoricalPerformance{
				Id:        primitive.NewObjectID(),
				Day:       day,
				BankId:    bankID,
				Value:     DefaultPerformanceValue,
				IsClaimed: true,
			}
			newEntries = append(newEntries, newClaimedEntry)
			finalClaimedHistory = append(finalClaimedHistory, newClaimedEntry)
		}
	}

	// Insert new claimed entries if any
	if len(newEntries) > 0 {
		err := s.historyRepo.CreateMany(ctx, newEntries)
		if err != nil {
			return nil, err
		}
	}

	return finalClaimedHistory, nil
}

func (s *PerformanceService) CreateInitialPerformanceHistory(
	ctx context.Context,
	bankID primitive.ObjectID,
	initialCapital int64,
) error {
	currentDay, err := s.gameService.GetCurrentDay(ctx)
	if err != nil {
		return err
	}

	startDay := currentDay - PerformanceHistoryDays
	_, err = s.ensureClaimedHistory(ctx, bankID, startDay, currentDay, []models.HistoricalPerformance{})
	if err != nil {
		return err
	}

	// Create actual performance history for the current day
	actualEntry := &models.HistoricalPerformance{
		Id:        primitive.NewObjectID(),
		Day:       currentDay,
		BankId:    bankID,
		Value:     initialCapital,
		IsClaimed: false,
	}

	return s.historyRepo.Create(ctx, actualEntry)
}

// convertToResponse converts HistoricalPerformance to useful response format
func convertToResponse(history []models.HistoricalPerformance) []models.HistoricalPerformanceResponse {
	result := make([]models.HistoricalPerformanceResponse, len(history))
	for i, entry := range history {
		result[i] = models.HistoricalPerformanceResponse{
			Day:   entry.Day,
			Value: entry.Value,
		}
	}
	return result
}
