package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserRolesModel = (*customUserRolesModel)(nil)

type (
	// UserRolesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserRolesModel.
	UserRolesModel interface {
		userRolesModel
		withSession(session sqlx.Session) UserRolesModel
	}

	customUserRolesModel struct {
		*defaultUserRolesModel
	}
)

// NewUserRolesModel returns a model for the database table.
func NewUserRolesModel(conn sqlx.SqlConn) UserRolesModel {
	return &customUserRolesModel{
		defaultUserRolesModel: newUserRolesModel(conn),
	}
}

func (m *customUserRolesModel) withSession(session sqlx.Session) UserRolesModel {
	return NewUserRolesModel(sqlx.NewSqlConnFromSession(session))
}
