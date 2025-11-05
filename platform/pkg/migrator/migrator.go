package migrator

// Migrator интерфейс для выполнения миграций базы данных
type Migrator interface {
	// Up выполняет все pending миграции
	Up() error

	// Down откатывает последнюю миграцию
	Down() error
}
