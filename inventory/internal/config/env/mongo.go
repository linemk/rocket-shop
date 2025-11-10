package env

import (
	"os"
)

const (
	mongoURIEnv = "INVENTORY_MONGO_URI"
	mongoDBEnv  = "INVENTORY_MONGO_DB"
)

type mongoConfig struct {
	uri          string
	databaseName string
}

// NewMongoConfig создает конфигурацию MongoDB из переменных окружения
func NewMongoConfig() (*mongoConfig, error) {
	uri := os.Getenv(mongoURIEnv)
	if uri == "" {
		uri = "mongodb://inventory_user:inventory_password@localhost:27017"
	}

	dbName := os.Getenv(mongoDBEnv)
	if dbName == "" {
		dbName = "inventory_db"
	}

	return &mongoConfig{
		uri:          uri,
		databaseName: dbName,
	}, nil
}

func (c *mongoConfig) URI() string {
	return c.uri
}

func (c *mongoConfig) DatabaseName() string {
	return c.databaseName
}
