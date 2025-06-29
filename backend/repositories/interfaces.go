package repositories

import (
	"context"
	"ponziworld/backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PlayerRepository defines the interface for player database operations
type PlayerRepository interface {
	Create(ctx context.Context, player *models.Player) error
	FindByUsername(ctx context.Context, username string) (*models.Player, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Player, error)
}

// BankRepository defines the interface for bank database operations
type BankRepository interface {
	Create(ctx context.Context, bank *models.Bank) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.Bank, error)
	FindByPlayerID(ctx context.Context, playerID primitive.ObjectID) (*models.Bank, error)
}

// AssetRepository defines the interface for asset database operations
type AssetRepository interface {
	Create(ctx context.Context, asset *models.Asset) error
	FindByBankID(ctx context.Context, bankID primitive.ObjectID) ([]models.Asset, error)
	CalculateActualCapital(ctx context.Context, bankID primitive.ObjectID) (int64, error)
}

// HistoricalPerformanceRepository defines the interface for historical performance database operations
type HistoricalPerformanceRepository interface {
	Create(ctx context.Context, performance *models.HistoricalPerformance) error
	CreateMany(ctx context.Context, performances []models.HistoricalPerformance) error
	FindByBankIDAndDateRange(ctx context.Context, bankID primitive.ObjectID, startDay, endDay int) ([]models.HistoricalPerformance, error)
	FindClaimedByBankIDAndDateRange(ctx context.Context, bankID primitive.ObjectID, startDay, endDay int) ([]models.HistoricalPerformance, error)
	FindActualByBankIDAndDateRange(ctx context.Context, bankID primitive.ObjectID, startDay, endDay int) ([]models.HistoricalPerformance, error)
}
