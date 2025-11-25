package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ PermissionsModel = (*customPermissionsModel)(nil)

type (
	// PermissionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPermissionsModel.
	PermissionsModel interface {
		permissionsModel
		withSession(session sqlx.Session) PermissionsModel
	}

	customPermissionsModel struct {
		*defaultPermissionsModel
	}
)

// NewPermissionsModel returns a model for the database table.
func NewPermissionsModel(conn sqlx.SqlConn) PermissionsModel {
	return &customPermissionsModel{
		defaultPermissionsModel: newPermissionsModel(conn),
	}
}

func (m *customPermissionsModel) withSession(session sqlx.Session) PermissionsModel {
	return NewPermissionsModel(sqlx.NewSqlConnFromSession(session))
}
