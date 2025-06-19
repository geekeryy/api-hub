package config

import (
	"github.com/SpectatorNan/gorm-zero/gormc/config/pg"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessExpire int64
	}
	PgSql pg.PgSql
	Jwks    Jwks
}

type Jwks struct {
	ServerURL string
	RefreshInterval int64
}