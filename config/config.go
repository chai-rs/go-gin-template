package config

import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)

var (
	PORT int

	REDIS_HOST     string
	REDIS_PORT     string
	REDIS_PASSWORD string
	REDIS_DB       int

	ACCESS_SECRET  string
	REFRESH_SECRET string
)

func Init() {
	var err error

	PORT, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal().Err(err).Msg("🚨 failed to convert PORT to int")
	}

	REDIS_HOST = os.Getenv("REDIS_HOST")
	REDIS_PORT = os.Getenv("REDIS_PORT")
	REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")

	REDIS_DB, err = strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatal().Err(err).Msg("🚨 failed to convert REDIS_DB to int")
	}

	ACCESS_SECRET = os.Getenv("ACCESS_SECRET")
	REFRESH_SECRET = os.Getenv("REFRESH_SECRET")
}
