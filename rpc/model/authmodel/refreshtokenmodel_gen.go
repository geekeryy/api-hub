// Code generated by goctl. DO NOT EDIT!

package authmodel

import (
	"context"
	"database/sql"
	"github.com/SpectatorNan/gorm-zero/gormc"
	"time"

	"github.com/SpectatorNan/gorm-zero/gormc/pagex"
	"gorm.io/gorm"
)

type (
	refreshTokenModel interface {
		Insert(ctx context.Context, tx *gorm.DB, data *RefreshToken) error
		BatchInsert(ctx context.Context, tx *gorm.DB, news []RefreshToken) error
		FindOne(ctx context.Context, id int64) (*RefreshToken, error)
		FindPageList(ctx context.Context, page *pagex.ListReq, orderBy pagex.OrderBy,
			orderKeys map[string]string, whereClause func(db *gorm.DB) *gorm.DB) ([]RefreshToken, int64, error)
		FindOneByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*RefreshToken, error)
		Update(ctx context.Context, tx *gorm.DB, data *RefreshToken) error
		BatchUpdate(ctx context.Context, tx *gorm.DB, olds, news []RefreshToken) error
		BatchDelete(ctx context.Context, tx *gorm.DB, datas []RefreshToken) error

		Delete(ctx context.Context, tx *gorm.DB, id int64) error
	}

	defaultRefreshTokenModel struct {
		conn  *gorm.DB
		table string
	}

	RefreshToken struct {
		Id               int64          `gorm:"column:id;primary_key"`
		RefreshTokenHash string         `gorm:"column:refresh_token_hash"`
		MemberId         string         `gorm:"column:member_id"`
		Status           int64          `gorm:"column:status"`
		ExpiredAt        time.Time      `gorm:"column:expired_at"`
		CreatedAt        sql.NullTime   `gorm:"column:created_at"`
		UpdatedAt        sql.NullTime   `gorm:"column:updated_at"`
		DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;index"`
	}
)

func (RefreshToken) TableName() string {
	return `"public"."refresh_token"`
}

func newRefreshTokenModel(db *gorm.DB) *defaultRefreshTokenModel {
	return &defaultRefreshTokenModel{
		conn:  db,
		table: `"public"."refresh_token"`,
	}
}

func (m *defaultRefreshTokenModel) Insert(ctx context.Context, tx *gorm.DB, data *RefreshToken) error {
	db := m.conn
	if tx != nil {
		db = tx
	}
	err := db.WithContext(ctx).Create(&data).Error
	return err
}
func (m *defaultRefreshTokenModel) BatchInsert(ctx context.Context, tx *gorm.DB, news []RefreshToken) error {
	db := m.conn
	if tx != nil {
		db = tx
	}
	err := db.Create(&news).Error
	return err
}

func (m *defaultRefreshTokenModel) FindOne(ctx context.Context, id int64) (*RefreshToken, error) {
	var resp RefreshToken
	err := m.conn.WithContext(ctx).Model(&RefreshToken{}).Where("id = ?", id).Take(&resp).Error
	switch err {
	case nil:
		return &resp, nil
	case gormc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (m *defaultRefreshTokenModel) FindPageList(ctx context.Context, page *pagex.ListReq, orderBy pagex.OrderBy,
	orderKeys map[string]string, whereClause func(db *gorm.DB) *gorm.DB) ([]RefreshToken, int64, error) {
	conn := m.conn
	formatDB := func() (*gorm.DB, *gorm.DB) {
		db := conn.Model(&RefreshToken{})
		if whereClause != nil {
			db = whereClause(db)
		}
		return db, nil
	}

	res, total, err := pagex.FindPageListWithCount[RefreshToken](ctx, page, orderBy, orderKeys, formatDB)
	return res, total, err
}

func (m *defaultRefreshTokenModel) FindOneByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*RefreshToken, error) {
	var resp RefreshToken
	err := m.conn.WithContext(ctx).Model(&RefreshToken{}).Where("refresh_token_hash = $1", refreshTokenHash).Take(&resp).Error
	switch err {
	case nil:
		return &resp, nil
	case gormc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultRefreshTokenModel) Update(ctx context.Context, tx *gorm.DB, data *RefreshToken) error {
	db := m.conn
	if tx != nil {
		db = tx
	}
	err := db.WithContext(ctx).Save(data).Error
	return err
}
func (m *defaultRefreshTokenModel) BatchUpdate(ctx context.Context, tx *gorm.DB, olds, news []RefreshToken) error {
	db := m.conn
	if tx != nil {
		db = tx
	}
	err := db.WithContext(ctx).Save(&news).Error

	return err
}

func (m *defaultRefreshTokenModel) Delete(ctx context.Context, tx *gorm.DB, id int64) error {
	db := m.conn
	if tx != nil {
		db = tx
	}
	err := db.WithContext(ctx).Delete(&RefreshToken{}, id).Error

	return err
}

func (m *defaultRefreshTokenModel) BatchDelete(ctx context.Context, tx *gorm.DB, datas []RefreshToken) error {
	db := m.conn
	if tx != nil {
		db = tx
	}
	err := db.Delete(&datas).Error
	return err
}
