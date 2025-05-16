package main

import (
	"github.com/chai-rs/simple-bookstore/config"
	_ "github.com/chai-rs/simple-bookstore/infrastructure/logger"
	"github.com/chai-rs/simple-bookstore/pkg/migration"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

const FILE = "db/migrations"

var migrator migration.DatabaseMigration

func init() {
	config.Init()
}

func main() {
	input := &migration.PostgresMigrationInput{
		Username: config.POSTGRES_USER,
		Password: config.POSTGRES_PASSWORD,
		Host:     config.POSTGRES_HOST,
		Port:     config.POSTGRES_PORT,
		DBName:   config.POSTGRES_DB,
		File:     FILE,
	}

	migrator = migration.NewPostgresMigration(input)
	if err := migrator.Up(); err != nil {
		log.Fatal().Err(err).Str("file", input.MigrationFile()).Str("url", input.URL()).Msg("🚨 failed to migrate")
	}

	log.Info().Str("file", input.MigrationFile()).Str("url", input.URL()).Msg("🚀 migrated successfully")
}
