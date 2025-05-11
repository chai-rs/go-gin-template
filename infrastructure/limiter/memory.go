package limiter

import (
	"github.com/rs/zerolog/log"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func NewMemoryLimiter(format string) *limiter.Limiter {
	rate, err := limiter.NewRateFromFormatted(format)
	if err != nil {
		log.Fatal().Err(err).Msg("💣 failed to create rate")
	}

	store := memory.NewStore()
	return limiter.New(store, rate)
}
