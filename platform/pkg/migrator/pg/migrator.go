package pg

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

// Migrator реализует интерфейс migrator.Migrator для PostgreSQL
type Migrator struct {
	db            *sql.DB
	migrationsDir string
}

// NewMigrator создает новый экземпляр мигратора для PostgreSQL
func NewMigrator(db *sql.DB, migrationsDir string) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// Up выполняет все pending миграции
func (m *Migrator) Up() error {
	if err := goose.Up(m.db, m.migrationsDir); err != nil {
		return err
	}

	return nil
}

// Down откатывает последнюю миграцию
func (m *Migrator) Down() error {
	if err := goose.Down(m.db, m.migrationsDir); err != nil {
		return err
	}

	return nil
}
