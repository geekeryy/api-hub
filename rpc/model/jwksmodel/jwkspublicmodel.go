package jwksmodel

import (
	"gorm.io/gorm"
)

var _ JwksPublicModel = (*customJwksPublicModel)(nil)

type (
	// JwksPublicModel is an interface to be customized, add more methods here,
	// and implement the added methods in customJwksPublicModel.
	JwksPublicModel interface {
		jwksPublicModel
		customJwksPublicLogicModel
	}

	customJwksPublicModel struct {
		*defaultJwksPublicModel
	}

	customJwksPublicLogicModel interface {
		FindAll() ([]*JwksPublic, error)
	}
)

// NewJwksPublicModel returns a model for the database table.
func NewJwksPublicModel(conn *gorm.DB) JwksPublicModel {
	return &customJwksPublicModel{
		defaultJwksPublicModel: newJwksPublicModel(conn),
	}
}

// FindAll
func (c *customJwksPublicModel) FindAll() ([]*JwksPublic, error) {
	var jwksPublics []*JwksPublic
	err := c.conn.Model(&JwksPublic{}).Order("id desc").Find(&jwksPublics).Error
	if err != nil {
		return nil, err
	}
	return jwksPublics, nil
}
