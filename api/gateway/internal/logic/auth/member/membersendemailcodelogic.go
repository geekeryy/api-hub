package member

import (
	"context"
	"fmt"
	"time"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/core/email"
	"github.com/geekeryy/api-hub/core/limiter"
	"github.com/geekeryy/api-hub/core/xstrings"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberSendEmailCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发送邮箱验证码
func NewMemberSendEmailCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberSendEmailCodeLogic {
	return &MemberSendEmailCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberSendEmailCodeLogic) MemberSendEmailCode(req *types.MemberSendEmailCodeReq) error {
	limit, ok := l.svcCtx.CodeLimiter.Get(req.Email)
	if !ok {
		limit = limiter.NewLimiter(time.Minute*10, time.Minute, 2, 10)
		l.svcCtx.CodeLimiter.Add(req.Email, limit)
	}
	if !limit.Validate() {
		return fmt.Errorf("发送频率过高")
	}
	code := xstrings.GenerateRandomNumber(6)
	l.Infof("send email code to %s code: %s", req.Email, code)
	err := email.New(l.svcCtx.Config.MailGun.Domain, l.svcCtx.Config.MailGun.ApiKey).SendMailGun(&email.SendMsg{
		Subject: "[API-HUB] 邮箱验证码",
		Body:    fmt.Sprintf("您的邮箱验证码是：%s", code),
		To:      []string{req.Email},
	}, l.svcCtx.Config.MailGun.Sender)
	if err != nil {
		return err
	}
	if err := l.svcCtx.Cache.Set(fmt.Sprintf("email_code_%s", req.Email), code, time.Minute*10); err != nil {
		return err
	}

	return nil
}
