package env

import (
	"os"
)

const (
	postgresDBURIEnv         = "ORDER_DB_URI"
	postgresMigrationsDirEnv = "ORDER_MIGRATIONS_DIR"
)

type postgresConfig struct {
	dsn           string
	migrationsDir string
}

// NewPostgresConfig создает конфигурацию PostgreSQL из переменных окружения
func NewPostgresConfig() (*postgresConfig, error) {
	dsn := os.Getenv(postgresDBURIEnv)
	if dsn == "" {
		dsn = "postgres://order_user:order_password@localhost:5432/order_db?sslmode=disable"
	}

	migrationsDir := os.Getenv(postgresMigrationsDirEnv)
	if migrationsDir == "" {
		migrationsDir = "migrations"
	}

	return &postgresConfig{
		dsn:           dsn,
		migrationsDir: migrationsDir,
	}, nil
}

func (c *postgresConfig) DSN() string {
	return c.dsn
}

func (c *postgresConfig) MigrationsDir() string {
	return c.migrationsDir
}
