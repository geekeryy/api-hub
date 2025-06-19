package authmodel

import (
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
		FindAll() ([]*Jwks, error)
		FindLatest() (*Jwks, error)
	}
)

// NewJwksModel returns a model for the database table.
func NewJwksModel(conn *gorm.DB) JwksModel {
	return &customJwksModel{
		defaultJwksModel: newJwksModel(conn),
	}
}


// FindAll
func (c *customJwksModel) FindAll() ([]*Jwks, error) {
	var jwkss []*Jwks
	err := c.conn.Model(&Jwks{}).Order("id desc").Find(&jwkss).Error
	if err != nil {
		return nil, err
	}
	return jwkss, nil
}

// FindLatest
func (c *customJwksModel) FindLatest() (*Jwks, error) {
	var jwks Jwks
	err := c.conn.Model(&Jwks{}).Order("id desc").First(&jwks).Error
	if err != nil {
		return nil, err
	}
	return &jwks, nil
}
