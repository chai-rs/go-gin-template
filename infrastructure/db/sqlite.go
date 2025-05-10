package db

import (
	driver "github.com/glebarez/sqlite"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var sqlite *gorm.DB

func SQLiteConnect() {
	var err error
	sqlite, err = gorm.Open(driver.Open("data.sqlite"), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("💣 Failed to connect to SQLite")
	}

	log.Debug().Msg("🔌 Connected to SQLite")
}

func SQLite() *gorm.DB {
	return sqlite
}
