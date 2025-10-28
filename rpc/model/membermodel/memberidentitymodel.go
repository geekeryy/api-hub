package membermodel

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MemberIdentityModel = (*customMemberIdentityModel)(nil)

type (
	// MemberIdentityModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMemberIdentityModel.
	MemberIdentityModel interface {
		memberIdentityModel
		withSession(session sqlx.Session) MemberIdentityModel
		FindByMemberUUID(ctx context.Context, memberUUID string) ([]*MemberIdentity, error)
		FindByIdentity(ctx context.Context, identityType int64, identifier string) ([]*MemberIdentity, error)
		UpdateCredential(ctx context.Context, id int64, credential string) error
		BatchInsert(ctx context.Context, session sqlx.Session, data []MemberIdentity) error
	}

	customMemberIdentityModel struct {
		*defaultMemberIdentityModel
	}
)

// NewMemberIdentityModel returns a model for the database table.
func NewMemberIdentityModel(conn sqlx.SqlConn) MemberIdentityModel {
	return &customMemberIdentityModel{
		defaultMemberIdentityModel: newMemberIdentityModel(conn),
	}
}

func (m *customMemberIdentityModel) withSession(session sqlx.Session) MemberIdentityModel {
	return NewMemberIdentityModel(sqlx.NewSqlConnFromSession(session))
}

// FindByMemberUUID finds the member identity by member UUID.
func (m *customMemberIdentityModel) FindByMemberUUID(ctx context.Context, memberUUID string) ([]*MemberIdentity, error) {
	query := fmt.Sprintf("select %s from %s where `member_uuid` = ?", memberIdentityRows, m.table)
	var resp []*MemberIdentity
	err := m.conn.QueryRowsCtx(ctx, &resp, query, memberUUID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// FindByIdentity finds the member identity by identity type and identifier.
func (m *customMemberIdentityModel) FindByIdentity(ctx context.Context, identityType int64, identifier string) ([]*MemberIdentity, error) {
	query := fmt.Sprintf("select %s from %s where `identity_type` = ? and `identifier` = ?", memberIdentityRows, m.table)
	var resp []*MemberIdentity
	err := m.conn.QueryRowsCtx(ctx, &resp, query, identityType, identifier)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customMemberIdentityModel) UpdateCredential(ctx context.Context, id int64, credential string) error {
	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err := session.ExecCtx(ctx, fmt.Sprintf("update %s set `credential` = ? where `id` = ?", m.table), credential, id)
		if err != nil {
			return err
		}
		// TODO 密码更新记录，禁止使用最近设置的密码
		return nil
	})
	return err
}

// BatchInsert
func (m *customMemberIdentityModel) BatchInsert(ctx context.Context, session sqlx.Session, data []MemberIdentity) error {
	if session == nil {
		session = m.conn
	}
	if len(data) > 100 {
		for i := 0; i < len(data); i += 100 {
			if err := m.batchInsert(ctx, session, data[i:min(i+100, len(data))]); err != nil {
				return err
			}
		}
		return nil
	} else {
		return m.batchInsert(ctx, session, data)
	}

}

func (m *customMemberIdentityModel) batchInsert(ctx context.Context, session sqlx.Session, data []MemberIdentity) error {
	args := make([]interface{}, 0)
	querys := make([]string, 0)
	for _, item := range data {
		querys = append(querys, "(?, ?, ?, ?, ?)")
		args = append(args, item.MemberUuid, item.IdentityType, item.Identifier, item.Credential, item.Status)
	}
	query := fmt.Sprintf("insert into %s (%s) values %s", m.table, memberIdentityRowsExpectAutoSet, strings.Join(querys, ", "))
	_, err := session.ExecCtx(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
