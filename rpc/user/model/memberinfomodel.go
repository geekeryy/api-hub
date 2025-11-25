package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ MemberInfoModel = (*customMemberInfoModel)(nil)

type (
	// MemberInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMemberInfoModel.
	MemberInfoModel interface {
		memberInfoModel
		withSession(session sqlx.Session) MemberInfoModel
	}

	customMemberInfoModel struct {
		*defaultMemberInfoModel
	}
)

// NewMemberInfoModel returns a model for the database table.
func NewMemberInfoModel(conn sqlx.SqlConn) MemberInfoModel {
	return &customMemberInfoModel{
		defaultMemberInfoModel: newMemberInfoModel(conn),
	}
}

func (m *customMemberInfoModel) withSession(session sqlx.Session) MemberInfoModel {
	return NewMemberInfoModel(sqlx.NewSqlConnFromSession(session))
}
