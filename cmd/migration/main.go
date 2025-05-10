package main

import (
	_ "github.com/0xanonydxck/simple-bookstore/infrastructure/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

func main() {
	m, err := migrate.New(
		"file://db/migrations",
		"sqlite3://data.sqlite",
	)
	if err != nil {
		log.Fatal().Err(err).Msg("💣 Failed to create migration")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("💣 Failed to apply migration")
	}

	log.Info().Msg("✅ Migration applied successfully.")
}
