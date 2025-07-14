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

// AssetTypeRepository defines the interface for asset type database operations
type AssetTypeRepository interface {
	Create(ctx context.Context, assetType *models.AssetType) error
	FindAll(ctx context.Context) ([]models.AssetType, error)
	FindByName(ctx context.Context, name string) (*models.AssetType, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.AssetType, error)
	UpsertByName(ctx context.Context, assetType *models.AssetType) error
}

// HistoricalPerformanceRepository defines the interface for historical performance database operations
type HistoricalPerformanceRepository interface {
	Create(ctx context.Context, performance *models.HistoricalPerformance) error
	CreateMany(ctx context.Context, performances []models.HistoricalPerformance) error
	FindByBankIDAndDateRange(ctx context.Context, bankID primitive.ObjectID, startDay, endDay int) ([]models.HistoricalPerformance, error)
	FindClaimedByBankIDAndDateRange(ctx context.Context, bankID primitive.ObjectID, startDay, endDay int) ([]models.HistoricalPerformance, error)
	FindActualByBankIDAndDateRange(ctx context.Context, bankID primitive.ObjectID, startDay, endDay int) ([]models.HistoricalPerformance, error)
}

// GameRepository defines the interface for game database operations
type GameRepository interface {
	GetCurrentDay(ctx context.Context) (int, error)
	IncrementDay(ctx context.Context) (int, error)
	CreateInitialGame(ctx context.Context, initialDay int) error
}

// PendingTransactionRepository defines the interface for pending transaction database operations
type PendingTransactionRepository interface {
	Create(ctx context.Context, transaction *models.PendingTransaction) error
	FindByBuyerBankID(ctx context.Context, buyerBankID primitive.ObjectID) ([]models.PendingTransaction, error)
	FindByBuyerBankIDAndAssetID(ctx context.Context, buyerBankID, assetID primitive.ObjectID) ([]models.PendingTransaction, error)
	UpdateAmount(ctx context.Context, id primitive.ObjectID, newAmount int64) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
