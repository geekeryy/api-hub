package ai

import (
	"context"
	"errors"
	"time"

	"github.com/geekeryy/api-hub/api/gateway/internal/svc"
	"github.com/geekeryy/api-hub/api/gateway/internal/types"
	"github.com/geekeryy/api-hub/library/xerror"
	"golang.org/x/time/rate"

	"github.com/sashabaranov/go-openai"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/logx"
)

type DailySentenceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 每日一句
func NewDailySentenceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DailySentenceLogic {
	return &DailySentenceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var cache, _ = collection.NewCache(10 * time.Minute)
var limiter = rate.NewLimiter(rate.Every(1*time.Minute), 10)

func (l *DailySentenceLogic) DailySentence(req *types.DailySentenceReq) (resp *types.DailySentenceResp, err error) {
	if content, ok := cache.Get(req.Use); ok && len(content.(string)) > 0 {
		return &types.DailySentenceResp{Sentence: content.(string)}, nil
	}
	if !limiter.Allow() {
		return nil, xerror.RequestRateLimitErr
	}
	config := openai.DefaultConfig(l.svcCtx.Config.Deepseek.ApiKey)
	config.BaseURL = "https://api.deepseek.com"
	client := openai.NewClientWithConfig(config)

	response, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: "deepseek-chat",
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: "你是一个才华横溢诗人，对于用户的问题，你每次只回答一句10字以上的优美的诗句，不要回答任何其他内容。"},
			{Role: "user", Content: "提示词为：" + req.Use + "，请用" + req.Lang + "语言回答。"},
		},
		Stream: false,
	})
	if err != nil {
		return nil, xerror.New(err, xerror.InternalServerErr)
	}
	if len(response.Choices) == 0 {
		return nil, xerror.New(errors.New("no response"), xerror.NotFoundErr)
	}
	cache.Set(req.Use, response.Choices[0].Message.Content)

	return &types.DailySentenceResp{Sentence: response.Choices[0].Message.Content}, nil
}
