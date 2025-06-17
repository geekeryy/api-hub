package usermodel

import (
	"gorm.io/gorm"
)

var _ MemberInfoModel = (*customMemberInfoModel)(nil)

type (
	// MemberInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMemberInfoModel.
	MemberInfoModel interface {
		memberInfoModel
		customMemberInfoLogicModel
	}

	customMemberInfoModel struct {
		*defaultMemberInfoModel
	}

	customMemberInfoLogicModel interface {
	}
)

// NewMemberInfoModel returns a model for the database table.
func NewMemberInfoModel(conn *gorm.DB) MemberInfoModel {
	return &customMemberInfoModel{
		defaultMemberInfoModel: newMemberInfoModel(conn),
	}
}
