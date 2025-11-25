package svc

import (
	"fmt"
	"log"

	"github.com/geekeryy/api-hub/rpc/user/internal/config"
	"github.com/geekeryy/api-hub/rpc/user/model"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                  config.Config
	DB                      sqlx.SqlConn
	MemberInfoModel         model.MemberInfoModel
	JwksModel               model.JwksModel
	TokenRefreshRecordModel model.TokenRefreshRecordModel
	RefreshTokenModel       model.RefreshTokenModel
	MemberIdentityModel     model.MemberIdentityModel
	RedisClient             *redis.Client
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
		Addr:     c.RedisConf.Addr,
		Password: c.RedisConf.Password,
		DB:       c.RedisConf.Db,
	})
	if err != nil {
		log.Fatalf("failed to open redis: %v", err)
	}

	svc := &ServiceContext{
		Config:                  c,
		MemberInfoModel:         model.NewMemberInfoModel(mysqlClient),
		DB:                      mysqlClient,
		JwksModel:               model.NewJwksModel(mysqlClient),
		TokenRefreshRecordModel: model.NewTokenRefreshRecordModel(mysqlClient),
		RefreshTokenModel:       model.NewRefreshTokenModel(mysqlClient),
		MemberIdentityModel:     model.NewMemberIdentityModel(mysqlClient),
		RedisClient:             redisClient,
	}

	return svc
}
