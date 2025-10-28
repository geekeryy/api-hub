package authmodel

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RefreshTokenModel = (*customRefreshTokenModel)(nil)

type (
	// RefreshTokenModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRefreshTokenModel.
	RefreshTokenModel interface {
		refreshTokenModel
		withSession(session sqlx.Session) RefreshTokenModel
	}

	customRefreshTokenModel struct {
		*defaultRefreshTokenModel
	}
)

// NewRefreshTokenModel returns a model for the database table.
func NewRefreshTokenModel(conn sqlx.SqlConn) RefreshTokenModel {
	return &customRefreshTokenModel{
		defaultRefreshTokenModel: newRefreshTokenModel(conn),
	}
}

func (m *customRefreshTokenModel) withSession(session sqlx.Session) RefreshTokenModel {
	return NewRefreshTokenModel(sqlx.NewSqlConnFromSession(session))
}
