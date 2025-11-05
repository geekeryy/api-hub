package config

import (
	"errors"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessExpire  int `json:",env=AUTH_ACCESS_EXPIRE,default=600"`
		RefreshExpire int `json:",env=AUTH_REFRESH_EXPIRE,default=2592000"`
	}
	Mysql          Mysql
	Redis          Redis
	MailGun        MailGun
	Secret         Secret
	Oms            Oms
	MemberService  zrpc.RpcClientConf
	MonitorService zrpc.RpcClientConf
}

func (c *Config) Validate() error {
	if c.Auth.AccessExpire <= 0 {
		return errors.New("AUTH_ACCESS_EXPIRE must be greater than 0")
	}
	return nil
}

type Mysql struct {
	Username string `json:",env=MYSQL_USERNAME"`
	Password string `json:",env=MYSQL_PASSWORD"`
	Host     string `json:",env=MYSQL_HOST"`
	Port     int    `json:",env=MYSQL_PORT,default=3306"`
	Dbname   string `json:",env=MYSQL_DBNAME"`
}

type Redis struct {
	Addr     string `json:",env=REDIS_ADDR"`
	Password string `json:",env=REDIS_PASSWORD"`
	Db       int    `json:",env=REDIS_DB,default=0"`
}

type Secret struct {
	RefreshToken string `json:",env=SECRET_REFRESH_TOKEN"`
	PrivateKey   string `json:",env=SECRET_PRIVATE_KEY"`
	PublicKey    string `json:",env=SECRET_PUBLIC_KEY"`
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
