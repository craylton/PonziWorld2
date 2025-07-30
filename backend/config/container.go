package config

import (
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Container struct {
	DatabaseConfig      *DatabaseConfig
	ServiceContainer    *ServiceContainer
	RepositoryContainer *RepositoryContainer
	Logger              zerolog.Logger
}

func NewContainer(
	client *mongo.Client,
	databaseName string,
	logger zerolog.Logger,
) *Container {
	dbConfig := &DatabaseConfig{
		DatabaseName: databaseName,
		Client:       client,
	}

	repositoryContainer := NewRepositoryContainer(dbConfig.GetDatabase())
	serviceContainer := NewServiceContainer(repositoryContainer, logger)

	return &Container{
		DatabaseConfig:      dbConfig,
		ServiceContainer:    serviceContainer,
		RepositoryContainer: repositoryContainer,
		Logger:              logger,
	}
}

// Close properly closes the database connection
func (d *Container) Close() {
	d.DatabaseConfig.Close()
}
