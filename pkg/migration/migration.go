package migration

// DatabaseMigration is the interface for database migrations.
type DatabaseMigration interface {
	Up() error
	Down() error
}

// Input is the interface for database migration input.
type Input interface {
	URL() string
	MigrationFile() string
}
