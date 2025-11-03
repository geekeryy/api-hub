package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Mysql          Mysql
	Redis          Redis
	Jwks           Jwks
	MailGun        MailGun
	Secret         Secret
	UserService    zrpc.RpcClientConf
	AuthService    zrpc.RpcClientConf
	MonitorService zrpc.RpcClientConf
}

type Monitor struct {
	ReportInterval int64 `json:",env=MONITOR_REPORTINTERVAL"`
}

type Mysql struct {
	Username string `json:",env=MYSQL_USERNAME"`
	Password string `json:",env=MYSQL_PASSWORD"`
	Host     string `json:",env=MYSQL_HOST"`
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

type Jwks struct {
	RefreshInterval int `json:",env=JWKS_REFRESH_INTERVAL"`
}

type MailGun struct {
	Domain string `json:",env=MAILGUN_DOMAIN"`
	ApiKey string `json:",env=MAILGUN_API_KEY"`
	Sender string `json:",env=MAILGUN_SENDER"`
}
