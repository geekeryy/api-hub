package svc

import (
	"log"

	"github.com/geekeryy/api-hub/core/pgcache"
	"github.com/geekeryy/api-hub/core/xgorm"
	"github.com/geekeryy/api-hub/rpc/model/usermodel"
	"github.com/geekeryy/api-hub/rpc/user/internal/config"
)

type ServiceContext struct {
	Config          config.Config
	MemberInfoModel usermodel.MemberInfoModel
	Cache           *pgcache.Cache
}

func NewServiceContext(c config.Config) *ServiceContext {
	pg, err := xgorm.ConnectPg(c.PgSql)
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %s", err)
	}

	cache, err := pgcache.NewCache(c.PgSql)
	if err != nil {
		log.Fatalf("Failed to init cache. Error: %s", err)
	}

	svc := &ServiceContext{
		Config:          c,
		MemberInfoModel: usermodel.NewMemberInfoModel(pg),
		Cache:           cache,
	}

	return svc
}
