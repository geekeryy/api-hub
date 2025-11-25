package config

import (
	"errors"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Jwt         Jwt
	Facebook    Facebook
	Mysql       Mysql
	RedisConf   RedisConf
	Secret      Secret
	Github      Github
	UserService zrpc.RpcClientConf
}

func (c *Config) Validate() error {
	if c.Jwt.AccessExpire <= 0 {
		return errors.New("AUTH_ACCESS_EXPIRE must be greater than 0")
	}
	return nil
}

type Github struct {
	ClientID     string `json:",env=GITHUB_CLIENT_ID"`
	ClientSecret string `json:",env=GITHUB_CLIENT_SECRET"`
	RedirectUri  string `json:",env=GITHUB_REDIRECT_URI"`
}

type Jwt struct {
	AccessExpire  int `json:",env=AUTH_ACCESS_EXPIRE,default=600"`
	RefreshExpire int `json:",env=AUTH_REFRESH_EXPIRE,default=2592000"`
}

type Mysql struct {
	Username string `json:",env=MYSQL_USERNAME"`
	Password string `json:",env=MYSQL_PASSWORD"`
	Host     string `json:",env=MYSQL_HOST"`
	Dbname   string `json:",env=MYSQL_DBNAME"`
}

type RedisConf struct {
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
