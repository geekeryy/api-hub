package authmodel

import (
	"context"

	"gorm.io/gorm"
)

var _ MemberIdentityModel = (*customMemberIdentityModel)(nil)

type (
	// MemberIdentityModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMemberIdentityModel.
	MemberIdentityModel interface {
		memberIdentityModel
		customMemberIdentityLogicModel
	}

	customMemberIdentityModel struct {
		*defaultMemberIdentityModel
	}

	customMemberIdentityLogicModel interface {
		FindByIdentity(ctx context.Context, identityType int64, identifier string) ([]MemberIdentity, error)
		FindByMemberId(ctx context.Context, memberId string) ([]MemberIdentity, error)
	}
)

// NewMemberIdentityModel returns a model for the database table.
func NewMemberIdentityModel(conn *gorm.DB) MemberIdentityModel {
	return &customMemberIdentityModel{
		defaultMemberIdentityModel: newMemberIdentityModel(conn),
	}
}

// FindByIdentity 根据身份类型查询
func (m *customMemberIdentityModel) FindByIdentity(ctx context.Context, identityType int64, identifier string) ([]MemberIdentity, error) {
	var record []MemberIdentity
	if err := m.conn.WithContext(ctx).Where("identity_type = ? AND identifier = ?", identityType, identifier).Find(&record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

// FindByMemberId 根据memberId查询
func (m *customMemberIdentityModel) FindByMemberId(ctx context.Context, memberId string) ([]MemberIdentity, error) {
	var record []MemberIdentity
	if err := m.conn.WithContext(ctx).Where("member_id = ?", memberId).Find(&record).Error; err != nil {
		return nil, err
	}
	return record, nil
}
