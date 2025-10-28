package authmodel

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RolePermissionsModel = (*customRolePermissionsModel)(nil)

type (
	// RolePermissionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRolePermissionsModel.
	RolePermissionsModel interface {
		rolePermissionsModel
		withSession(session sqlx.Session) RolePermissionsModel
	}

	customRolePermissionsModel struct {
		*defaultRolePermissionsModel
	}
)

// NewRolePermissionsModel returns a model for the database table.
func NewRolePermissionsModel(conn sqlx.SqlConn) RolePermissionsModel {
	return &customRolePermissionsModel{
		defaultRolePermissionsModel: newRolePermissionsModel(conn),
	}
}

func (m *customRolePermissionsModel) withSession(session sqlx.Session) RolePermissionsModel {
	return NewRolePermissionsModel(sqlx.NewSqlConnFromSession(session))
}
