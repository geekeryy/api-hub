package authmodel

import (
	"context"

	"gorm.io/gorm"
)

var _ TokenRefreshRecordModel = (*customTokenRefreshRecordModel)(nil)

type (
	// TokenRefreshRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTokenRefreshRecordModel.
	TokenRefreshRecordModel interface {
		tokenRefreshRecordModel
		customTokenRefreshRecordLogicModel
	}

	customTokenRefreshRecordModel struct {
		*defaultTokenRefreshRecordModel
	}

	customTokenRefreshRecordLogicModel interface {
		FindByKid(ctx context.Context, kid string) ([]TokenRefreshRecord, error)
	}
)

// NewTokenRefreshRecordModel returns a model for the database table.
func NewTokenRefreshRecordModel(conn *gorm.DB) TokenRefreshRecordModel {
	return &customTokenRefreshRecordModel{
		defaultTokenRefreshRecordModel: newTokenRefreshRecordModel(conn),
	}
}

// FindByKid 根据kid查询
func (m *customTokenRefreshRecordModel) FindByKid(ctx context.Context, kid string) ([]TokenRefreshRecord, error) {
	var records []TokenRefreshRecord
	if err := m.conn.WithContext(ctx).Where("kid = ?", kid).Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}
