package authmodel

import (
	"gorm.io/gorm"
)

var _ RefreshTokenModel = (*customRefreshTokenModel)(nil)

type (
	// RefreshTokenModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRefreshTokenModel.
	RefreshTokenModel interface {
		refreshTokenModel
		customRefreshTokenLogicModel
	}

	customRefreshTokenModel struct {
		*defaultRefreshTokenModel
	}

	customRefreshTokenLogicModel interface {
	}
)

// NewRefreshTokenModel returns a model for the database table.
func NewRefreshTokenModel(conn *gorm.DB) RefreshTokenModel {
	return &customRefreshTokenModel{
		defaultRefreshTokenModel: newRefreshTokenModel(conn),
	}
}
