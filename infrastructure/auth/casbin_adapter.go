package auth

import (
	"github.com/casbin/casbin/v2/persist"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func GormAdapter(db *gorm.DB) persist.Adapter {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatal().Err(err).Msg("💣 failed to create gorm adapter")
	}

	return adapter
}
