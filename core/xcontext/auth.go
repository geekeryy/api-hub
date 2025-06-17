package xcontext

import (
	"context"

	"github.com/geekeryy/api-hub/core/consts"
	"google.golang.org/grpc/metadata"
)

type memberid struct{}

func WithMemberID(ctx context.Context, memberID string) context.Context {
	return context.WithValue(ctx, memberid{}, memberID)
}

func GetMemberID(ctx context.Context) string {
	memberID, ok := ctx.Value(memberid{}).(string)
	if !ok {
		return ""
	}
	return memberID
}

type adminid struct{}

func WithAdminID(ctx context.Context, adminID int64) context.Context {
	return context.WithValue(ctx, adminid{}, adminID)
}

func GetAdminID(ctx context.Context) int64 {
	adminID, ok := ctx.Value(adminid{}).(int64)
	if !ok {
		return 0
	}
	return adminID
}

type roleid struct{}

func WithRoleID(ctx context.Context, roleID int64) context.Context {
	return context.WithValue(ctx, roleid{}, roleID)
}

func GetRoleID(ctx context.Context) int64 {
	roleID, ok := ctx.Value(roleid{}).(int64)
	if !ok {
		return 0
	}
	return roleID
}

type clientip struct{}

func WithClientIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, clientip{}, ip)
}

func GetClientIp(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		if value := ctx.Value(clientip{}); value != nil {
			return value.(string)
		}
		return ""
	}
	ip, ok := md[consts.CONTEXT_CLIENT_IP]
	if !ok || len(ip) == 0 || len(ip[0]) == 0 {
		return ""
	}
	return ip[0]
}


