package migration

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

var _ Input = (*PostgresMigrationInput)(nil)

// PostgresMigration handles PostgreSQL database migrations.
type PostgresMigration struct {
	migrator *migrate.Migrate
}

// NewPostgresMigration creates a new PostgreSQL migration instance.
func NewPostgresMigration(input *PostgresMigrationInput) *PostgresMigration {
	migrator, err := migrate.New(input.MigrationFile(), input.URL())
	if err != nil {
		log.Fatal().Err(err).Msg("🚨 failed to create migrator")
	}

	return &PostgresMigration{migrator}
}

// Up applies all available migrations.
func (p *PostgresMigration) Up() error {
	if err := p.migrator.Up(); err != nil {
		return err
	}

	return nil
}

// Down rolls back the most recent migration.
func (p *PostgresMigration) Down() error {
	if err := p.migrator.Down(); err != nil {
		return err
	}

	return nil
}

// PostgresMigrationInput contains configuration for PostgreSQL migrations.
type PostgresMigrationInput struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
	File     string
}

// URL constructs the PostgreSQL connection URL.
func (p *PostgresMigrationInput) URL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", p.Username, p.Password, p.Host, p.Port, p.DBName)
}

// MigrationFile returns the file URL for migration files.
func (p *PostgresMigrationInput) MigrationFile() string {
	return fmt.Sprintf("file://%s", p.File)
}
