package svc

import (
	"fmt"
	"log"

	"github.com/geekeryy/api-hub/rpc/model/membermodel"
	"github.com/geekeryy/api-hub/rpc/user/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config          config.Config
	MemberInfoModel membermodel.MemberInfoModel
	DB              sqlx.SqlConn
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

	svc := &ServiceContext{
		Config:          c,
		MemberInfoModel: membermodel.NewMemberInfoModel(mysqlClient),
		DB:              mysqlClient,
	}

	return svc
}
