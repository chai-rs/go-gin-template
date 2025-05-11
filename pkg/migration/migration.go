package migration

type DatabaseMigration interface {
	Up() error
	Down() error
}

type Input interface {
	URL() string
	MigrationFile() string
}
