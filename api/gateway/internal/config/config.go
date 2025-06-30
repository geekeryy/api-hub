package config

import (
	"github.com/SpectatorNan/gorm-zero/gormc/config/pg"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessExpire  int64 `json:",env=AUTH_ACCESS_EXPIRE"`
		RefreshExpire int64 `json:",env=AUTH_REFRESH_EXPIRE"`
	}
	PgSql    pg.PgSql
	Jwks     Jwks
	Facebook Facebook
	MailGun  MailGun
	Secret   Secret
	Oms      Oms
}

type Secret struct {
	RefreshToken string `json:",env=SECRET_REFRESH_TOKEN"`
	PrivateKey   string `json:",env=SECRET_PRIVATE_KEY"`
	PublicKey    string `json:",env=SECRET_PUBLIC_KEY"`
}

type Jwks struct {
	ServerURL       string `json:",env=JWKS_SERVER_URL"`
	RefreshInterval int64  `json:",env=JWKS_REFRESH_INTERVAL"`
}

type Facebook struct {
	AppID     string `json:",env=FACEBOOK_APP_ID"`
	AppSecret string `json:",env=FACEBOOK_APP_SECRET"`
}

type MailGun struct {
	Domain string `json:",env=MAILGUN_DOMAIN"`
	ApiKey string `json:",env=MAILGUN_API_KEY"`
	Sender string `json:",env=MAILGUN_SENDER"`
}

type Oms struct {
	OtpSecret string `json:",env=OMS_OTP_SECRET"`
}
