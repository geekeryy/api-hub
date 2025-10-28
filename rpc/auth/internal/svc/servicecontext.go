package svc

import (
	"fmt"
	"log"

	"github.com/geekeryy/api-hub/rpc/auth/internal/config"
	"github.com/geekeryy/api-hub/rpc/model/authmodel"
	"github.com/geekeryy/api-hub/rpc/model/membermodel"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                  config.Config
	DB                      sqlx.SqlConn
	RedisClient             *redis.Client
	JwksModel               authmodel.JwksModel
	MemberIdentityModel     membermodel.MemberIdentityModel
	MemberInfoModel         membermodel.MemberInfoModel
	TokenRefreshRecordModel authmodel.TokenRefreshRecordModel
	RefreshTokenModel       authmodel.RefreshTokenModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlClient, err := sqlx.NewConn(sqlx.SqlConf{
		DataSource: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Mysql.Username, c.Mysql.Password, c.Mysql.Host, c.Mysql.Dbname),
		DriverName: "mysql",
		Replicas:   nil,
		Policy:     "",
	})
	if err != nil {
		log.Fatalf("failed to open mysql: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       c.Redis.Db,
	})
	if err != nil {
		log.Fatalf("failed to open redis: %v", err)
	}

	return &ServiceContext{
		Config:                  c,
		DB:                      mysqlClient,
		RedisClient:             redisClient,
		JwksModel:               authmodel.NewJwksModel(mysqlClient),
		MemberIdentityModel:     membermodel.NewMemberIdentityModel(mysqlClient),
		MemberInfoModel:         membermodel.NewMemberInfoModel(mysqlClient),
		TokenRefreshRecordModel: authmodel.NewTokenRefreshRecordModel(mysqlClient),
		RefreshTokenModel:       authmodel.NewRefreshTokenModel(mysqlClient),
	}
}
