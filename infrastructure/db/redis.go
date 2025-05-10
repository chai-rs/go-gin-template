package db

import (
	"context"
	"fmt"

	"github.com/0xanonydxck/simple-bookstore/config"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func Redis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.REDIS_HOST, config.REDIS_PORT),
		Password: config.REDIS_PASSWORD,
		DB:       config.REDIS_DB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal().Err(err).Msg("💣 failed to connect to redis")
	}

	return rdb
}
