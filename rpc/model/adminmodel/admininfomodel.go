package adminmodel

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AdminInfoModel = (*customAdminInfoModel)(nil)

type (
	// AdminInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminInfoModel.
	AdminInfoModel interface {
		adminInfoModel
		withSession(session sqlx.Session) AdminInfoModel
	}

	customAdminInfoModel struct {
		*defaultAdminInfoModel
	}
)

// NewAdminInfoModel returns a model for the database table.
func NewAdminInfoModel(conn sqlx.SqlConn) AdminInfoModel {
	return &customAdminInfoModel{
		defaultAdminInfoModel: newAdminInfoModel(conn),
	}
}

func (m *customAdminInfoModel) withSession(session sqlx.Session) AdminInfoModel {
	return NewAdminInfoModel(sqlx.NewSqlConnFromSession(session))
}
