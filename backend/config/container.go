package config

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Container struct {
	DatabaseConfig      *DatabaseConfig
	ServiceContainer    *ServiceContainer
	RepositoryContainer *RepositoryContainer
}

func NewContainer(
	client *mongo.Client,
	cancel context.CancelFunc,
	databaseName string,
) *Container {
	dbConfig := &DatabaseConfig{
		DatabaseName:     databaseName,
		Client:           client,
		connectionCancel: cancel,
	}

	repositoryContainer := NewRepositoryContainer(dbConfig.GetDatabase())
	serviceContainer := NewServiceContainer(repositoryContainer)

	return &Container{
		DatabaseConfig:      dbConfig,
		ServiceContainer:    serviceContainer,
		RepositoryContainer: repositoryContainer,
	}
}

// Close properly closes the database connection
func (d *Container) Close() {
	d.DatabaseConfig.Close()
}
