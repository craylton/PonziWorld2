package services

import (
	"context"
	"errors"
	"strings"
)

type PlayerService struct {
	authService        *AuthService
	bankService        *BankService
	assetService       *AssetService
	performanceService *PerformanceService
}

func NewPlayerService(
	authService *AuthService,
	bankService *BankService,
	assetService *AssetService,
	performanceService *PerformanceService,
) *PlayerService {
	return &PlayerService{
		authService:        authService,
		bankService:        bankService,
		assetService:       assetService,
		performanceService: performanceService,
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
	player, err := s.authService.CreatePlayer(ctx, username, password)
	if err != nil {
		return err
	}

	// Create the bank for this player
	initialCapital := int64(1000)
	bank, err := s.bankService.CreateBank(ctx, player.Id, bankName, initialCapital)
	if err != nil {
		return err
	}

	// Create initial cash asset
	_, err = s.assetService.CreateInitialAsset(ctx, bank.Id, initialCapital)
	if err != nil {
		return err
	}

	// Create initial performance history
	err = s.performanceService.CreateInitialPerformanceHistory(ctx, bank.Id, 0, initialCapital)
	if err != nil {
		return err
	}

	return nil
}
