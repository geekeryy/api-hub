package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql Mysql
}

type Mysql struct {
	Username string `json:",env=MYSQL_USERNAME"`
	Password string `json:",env=MYSQL_PASSWORD"`
	Host     string `json:",env=MYSQL_HOST"`
	Port     int    `json:",env=MYSQL_PORT,default=3306"`
	Dbname   string `json:",env=MYSQL_DBNAME"`
}
