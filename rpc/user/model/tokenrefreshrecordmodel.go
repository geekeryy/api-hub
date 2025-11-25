package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TokenRefreshRecordModel = (*customTokenRefreshRecordModel)(nil)

type (
	// TokenRefreshRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTokenRefreshRecordModel.
	TokenRefreshRecordModel interface {
		tokenRefreshRecordModel
		withSession(session sqlx.Session) TokenRefreshRecordModel
		FindByKid(ctx context.Context, kid string) ([]*TokenRefreshRecord, error)
	}

	customTokenRefreshRecordModel struct {
		*defaultTokenRefreshRecordModel
	}
)

// NewTokenRefreshRecordModel returns a model for the database table.
func NewTokenRefreshRecordModel(conn sqlx.SqlConn) TokenRefreshRecordModel {
	return &customTokenRefreshRecordModel{
		defaultTokenRefreshRecordModel: newTokenRefreshRecordModel(conn),
	}
}

func (m *customTokenRefreshRecordModel) withSession(session sqlx.Session) TokenRefreshRecordModel {
	return NewTokenRefreshRecordModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customTokenRefreshRecordModel) FindByKid(ctx context.Context, kid string) ([]*TokenRefreshRecord, error) {
	query := fmt.Sprintf("select %s from %s where `kid` = ?", tokenRefreshRecordRows, m.table)
	var resp []*TokenRefreshRecord
	err := m.conn.QueryRowsCtx(ctx, &resp, query, kid)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
