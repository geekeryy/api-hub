package config

import (
	"github.com/SpectatorNan/gorm-zero/gormc/config/pg"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessExpire int64
		RefreshExpire int64
	}
	PgSql pg.PgSql
	Jwks    Jwks
	Facebook Facebook
	MailGun MailGun
}

type Jwks struct {
	ServerURL string
	RefreshInterval int64
}

type Facebook struct {
	AppID     string
	AppSecret string
}

type MailGun struct {
	Domain string
	ApiKey string
	Sender string
}