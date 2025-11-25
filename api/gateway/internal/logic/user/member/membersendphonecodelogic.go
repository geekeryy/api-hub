// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package member

import (
	"context"
	"fmt"
	"time"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/core/limiter"
	"github.com/geekeryy/api-hub/core/xstrings"
	"github.com/geekeryy/api-hub/library/xerror"

	"github.com/zeromicro/go-zero/core/logx"
)

type MemberSendPhoneCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发送手机验证码
func NewMemberSendPhoneCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MemberSendPhoneCodeLogic {
	return &MemberSendPhoneCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberSendPhoneCodeLogic) MemberSendPhoneCode(req *types.MemberSendPhoneCodeReq) error {
	limit, ok := l.svcCtx.CodeLimiter.Get(req.Phone)
	if !ok {
		limit = limiter.NewLimiter(time.Minute*10, time.Minute, 2, 10)
		l.svcCtx.CodeLimiter.Add(req.Phone, limit)
	}
	if !limit.Validate() {
		return xerror.RequestRateLimitError
	}
	code := xstrings.GenerateRandomNumber(6)
	l.Infof("send phone code to %s code: %s", req.Phone, code)
	// TODO: 发送短信

	if err := l.svcCtx.RedisClient.Set(l.ctx, fmt.Sprintf("phone_code_%s", req.Phone), code, time.Minute*10).Err(); err != nil {
		return err
	}
	return nil
}
