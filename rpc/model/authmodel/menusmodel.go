package authmodel

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ MenusModel = (*customMenusModel)(nil)

type (
	// MenusModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMenusModel.
	MenusModel interface {
		menusModel
		withSession(session sqlx.Session) MenusModel
	}

	customMenusModel struct {
		*defaultMenusModel
	}
)

// NewMenusModel returns a model for the database table.
func NewMenusModel(conn sqlx.SqlConn) MenusModel {
	return &customMenusModel{
		defaultMenusModel: newMenusModel(conn),
	}
}

func (m *customMenusModel) withSession(session sqlx.Session) MenusModel {
	return NewMenusModel(sqlx.NewSqlConnFromSession(session))
}
