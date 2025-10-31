package xgrpc

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// LazyClient 懒加载grpc客户端连接
// 1. 在用户发起第一次grpc请求时才创建连接
// 2. 初次连接如果下游服务不可用，则创建一个fake连接，避免服务panic
// 3. 后续每隔interval秒尝试连接一次，如果成功则更新连接，退出循环
type LazyClient struct {
	ctx           context.Context
	config        zrpc.RpcClientConf
	clientOptions []zrpc.ClientOption
	interval      int64 // sync interval in seconds
	conn          *grpc.ClientConn
	once          sync.Once
}

func NewLazyClient(config zrpc.RpcClientConf, interval int64, clientOptions ...zrpc.ClientOption) *LazyClient {
	return NewLazyClientCtx(context.Background(), config, interval, clientOptions...)
}

func NewLazyClientCtx(ctx context.Context, config zrpc.RpcClientConf, interval int64, clientOptions ...zrpc.ClientOption) *LazyClient {
	if interval < 0 {
		interval = 3
	}
	return &LazyClient{
		ctx:           ctx,
		config:        config,
		clientOptions: clientOptions,
		interval:      interval,
	}
}

func (c *LazyClient) syncConn() {
	timer := time.NewTimer(time.Duration(c.interval) * time.Second)
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-timer.C:
			client, err := zrpc.NewClient(c.config, c.clientOptions...)
			if err != nil {
				logx.Errorf("Failed to create conn. Error: %s", err)
			} else {
				logx.Infof("Successfully created conn.")
				atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&c.conn)), unsafe.Pointer(client.Conn()))
				return
			}
		}
	}
}

func (c *LazyClient) Conn() *grpc.ClientConn {
	c.once.Do(func() {
		client, err := zrpc.NewClient(c.config, c.clientOptions...)
		if err != nil {
			logx.Errorf("Failed to create lazy conn. Error: %s", err)
		} else {
			logx.Infof("Successfully created lazy conn.")
			c.conn = client.Conn()
			return
		}

		opt := grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			logx.Infof("fake conn used, method: %s, req: %v, reply: %v", method, req, reply)
			return nil
		})
		conn, err := grpc.NewClient("", opt, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logx.Errorf("Failed to create fake conn. Error: %s", err)
			return
		}
		c.conn = conn
		go c.syncConn()
	})
	return (*grpc.ClientConn)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&c.conn))))
}
