package services

import (
	"context"
	"errors"
	"strings"
)

type PlayerService struct {
	authService                  *AuthService
	bankService                  *BankService
	assetService                 *InvestmentService
	historicalPerformanceService *HistoricalPerformanceService
}

func NewPlayerService(
	authService *AuthService,
	bankService *BankService,
	assetService *InvestmentService,
	historicalPerformanceService *HistoricalPerformanceService,
) *PlayerService {
	return &PlayerService{
		authService:                  authService,
		bankService:                  bankService,
		assetService:                 assetService,
		historicalPerformanceService: historicalPerformanceService,
	}
}

func (s *PlayerService) CreateNewPlayer(ctx context.Context, username, password, bankName string) error {
	// Validate input
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	bankName = strings.TrimSpace(bankName)

	if username == "" || password == "" || bankName == "" {
		return errors.New("username, password, and bank name required")
	}

	// Create the player
	_, err := s.authService.CreatePlayer(ctx, username, password)
	if err != nil {
		return err
	}

	// Create the bank for this player
	initialCapital := int64(1000)
	bank, err := s.bankService.CreateBankForUsername(ctx, username, bankName, initialCapital)
	if err != nil {
		return err
	}

	// Create initial performance history
	err = s.historicalPerformanceService.CreateInitialHistoricalPerformance(ctx, bank.Id, initialCapital)
	if err != nil {
		return err
	}

	return nil
}
