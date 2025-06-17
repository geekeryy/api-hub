package authmodel

import (
	"context"

	"gorm.io/gorm"
)

var _ JwksModel = (*customJwksModel)(nil)

type (
	// JwksModel is an interface to be customized, add more methods here,
	// and implement the added methods in customJwksModel.
	JwksModel interface {
		jwksModel
		customJwksLogicModel
	}

	customJwksModel struct {
		*defaultJwksModel
	}

	customJwksLogicModel interface {
		FindAll(ctx context.Context) ([]*Jwks, error)
		FindLatest(ctx context.Context) (*Jwks, error)
		DeleteByKid(ctx context.Context, kid string) error
	}
)

// NewJwksModel returns a model for the database table.
func NewJwksModel(conn *gorm.DB) JwksModel {
	return &customJwksModel{
		defaultJwksModel: newJwksModel(conn),
	}
}

// FindAll
func (c *customJwksModel) FindAll(ctx context.Context) ([]*Jwks, error) {
	var jwkss []*Jwks
	err := c.conn.WithContext(ctx).Model(&Jwks{}).Order("id desc").Find(&jwkss).Error
	if err != nil {
		return nil, err
	}
	return jwkss, nil
}

// FindLatest
func (c *customJwksModel) FindLatest(ctx context.Context) (*Jwks, error) {
	var jwks Jwks
	err := c.conn.WithContext(ctx).Model(&Jwks{}).Order("id desc").First(&jwks).Error
	if err != nil {
		return nil, err
	}
	return &jwks, nil
}

// DeleteByKid
func (c *customJwksModel) DeleteByKid(ctx context.Context, kid string) error {
	return c.conn.WithContext(ctx).Model(&Jwks{}).Where("kid = ?", kid).Delete(&Jwks{}).Error
}
