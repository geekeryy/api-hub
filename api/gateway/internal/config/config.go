package config

import (
	"github.com/SpectatorNan/gorm-zero/gormc/config/pg"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	PgSql pg.PgSql
}
