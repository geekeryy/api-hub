// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package explorer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/library/xerror"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/time/rate"
)

type RssProxyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// RSS代理
func NewRssProxyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RssProxyLogic {
	return &RssProxyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var cache, _ = collection.NewCache(30 * time.Minute)
var limiter = rate.NewLimiter(rate.Every(1*time.Minute), 60)

func (l *RssProxyLogic) RssProxy(req *types.RssProxyReq) (resp *types.RssProxyResp, err error) {
	if content, ok := cache.Get(req.Url); ok && len(content.(string)) > 0 {
		return &types.RssProxyResp{Content: content.(string)}, nil
	}
	if !limiter.Allow() {
		return nil, xerror.RequestRateLimitErr
	}
	response, err := http.Get(req.Url)
	if err != nil {
		return nil, xerror.New(err, xerror.RequestFailedErr)
	}
	if response.StatusCode != http.StatusOK {
		return nil, xerror.New(fmt.Errorf("request failed with status code %d", response.StatusCode), xerror.RequestFailedErr)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, xerror.New(err, xerror.RequestFailedErr)
	}
	if len(body) > 0 {
		cache.Set(req.Url, string(body))
	}
	return &types.RssProxyResp{Content: string(body)}, nil
}
