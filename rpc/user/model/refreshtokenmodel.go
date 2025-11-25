package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RefreshTokenModel = (*customRefreshTokenModel)(nil)

type (
	// RefreshTokenModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRefreshTokenModel.
	RefreshTokenModel interface {
		refreshTokenModel
		withSession(session sqlx.Session) RefreshTokenModel
		FindOneByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*RefreshToken, error)
	}

	customRefreshTokenModel struct {
		*defaultRefreshTokenModel
	}
)

// NewRefreshTokenModel returns a model for the database table.
func NewRefreshTokenModel(conn sqlx.SqlConn) RefreshTokenModel {
	return &customRefreshTokenModel{
		defaultRefreshTokenModel: newRefreshTokenModel(conn),
	}
}

func (m *customRefreshTokenModel) withSession(session sqlx.Session) RefreshTokenModel {
	return NewRefreshTokenModel(sqlx.NewSqlConnFromSession(session))
}

// FindOneByRefreshTokenHash
func (m *customRefreshTokenModel) FindOneByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*RefreshToken, error) {
	query := fmt.Sprintf("select %s from %s where `refresh_token_hash` = ? limit 1", refreshTokenRows, m.table)
	var resp RefreshToken
	err := m.conn.QueryRowCtx(ctx, &resp, query, refreshTokenHash)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
