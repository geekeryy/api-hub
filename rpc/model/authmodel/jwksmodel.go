package authmodel

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ JwksModel = (*customJwksModel)(nil)

type (
	// JwksModel is an interface to be customized, add more methods here,
	// and implement the added methods in customJwksModel.
	JwksModel interface {
		jwksModel
		withSession(session sqlx.Session) JwksModel
		FindLatest(ctx context.Context) (*Jwks, error)
		DeleteByKid(ctx context.Context, kid string) error
		FindAll(ctx context.Context) ([]*Jwks, error)
	}

	customJwksModel struct {
		*defaultJwksModel
	}
)

// NewJwksModel returns a model for the database table.
func NewJwksModel(conn sqlx.SqlConn) JwksModel {
	return &customJwksModel{
		defaultJwksModel: newJwksModel(conn),
	}
}

func (m *customJwksModel) withSession(session sqlx.Session) JwksModel {
	return NewJwksModel(sqlx.NewSqlConnFromSession(session))
}

// FindLatest finds the latest jwks.
func (m *customJwksModel) FindLatest(ctx context.Context) (*Jwks, error) {
	query := fmt.Sprintf("select %s from %s order by `created_at` desc limit 1", jwksFieldNames, m.table)
	var resp Jwks
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *customJwksModel) DeleteByKid(ctx context.Context, kid string) error {
	query := fmt.Sprintf("delete from %s where `kid` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, kid)
	return err
}

func (m *customJwksModel) FindAll(ctx context.Context) ([]*Jwks, error) {
	query := fmt.Sprintf("select %s from %s", jwksRows, m.table)
	var resp []*Jwks
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
