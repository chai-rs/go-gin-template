package db

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
	driver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	postgres *gorm.DB
	once     sync.Once
)

// PostgreSQLConnect connects to the PostgreSQL database.
func PostgreSQLConnect(host, port, user, password, db string) *gorm.DB {
	once.Do(func() {
		var err error
		postgres, err = gorm.Open(driver.Open(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, db)), &gorm.Config{})
		if err != nil {
			log.Fatal().Err(err).Msg("💣 Failed to connect to PostgreSQL")
		}
		log.Debug().Msg("🔌 Connected to PostgreSQL")
	})

	return postgres
}

// PostgreSQL returns the PostgreSQL database connection.
func PostgreSQL() *gorm.DB {
	return postgres
}
