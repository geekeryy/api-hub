package config

import (
	"errors"

	"github.com/SpectatorNan/gorm-zero/gormc/config/pg"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessExpire  int `json:",env=AUTH_ACCESS_EXPIRE,default=600"`
		RefreshExpire int `json:",env=AUTH_REFRESH_EXPIRE,default=2592000"`
	}
	PgSql    pg.PgSql
	Jwks     Jwks
	Facebook Facebook
	MailGun  MailGun
	Secret   Secret
	Oms      Oms
}

func (c *Config) Validate() error {
	if c.Auth.AccessExpire <= 0 {
		return errors.New("AUTH_ACCESS_EXPIRE must be greater than 0")
	}
	return nil
}

type Secret struct {
	RefreshToken string `json:",env=SECRET_REFRESH_TOKEN"`
	PrivateKey   string `json:",env=SECRET_PRIVATE_KEY"`
	PublicKey    string `json:",env=SECRET_PUBLIC_KEY"`
}

type Jwks struct {
	ServerURL       string `json:",env=JWKS_SERVER_URL"`
	RefreshInterval int    `json:",env=JWKS_REFRESH_INTERVAL"`
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
