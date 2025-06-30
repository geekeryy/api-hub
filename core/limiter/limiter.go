package limiter

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

/*
 限流器
 控制指定事件间隔最多通过的请求数量
 NewLimiter(窗口大小，事件槽大小，每个事件最大请求数量，最大请求数量)
*/

type timeSlot struct {
	timestamp time.Time
	count     int64
}

// Limiter 内存限流器
type Limiter struct {
	sync.Mutex
	windows  []*timeSlot
	width    time.Duration
	size     time.Duration
	perCount int64
	maxCount int64
}

func NewLimiter(width time.Duration, size time.Duration, perCount int64, maxCount int64) *Limiter {
	return &Limiter{
		Mutex:    sync.Mutex{},
		windows:  make([]*timeSlot, 0),
		width:    width,
		size:     size,
		perCount: perCount,
		maxCount: maxCount,
	}
}

func (m *Limiter) SetMaxCount(maxCount int64) {
	atomic.StoreInt64(&m.maxCount, maxCount)
}

func (m *Limiter) Reset() {
	m.Lock()
	defer m.Unlock()
	m.windows = make([]*timeSlot, 0)
}

func (m *Limiter) Counter() int64 {
	m.Lock()
	defer m.Unlock()
	var sum int64
	for _, v := range m.windows {
		sum += v.count
	}
	return sum
}

func (m *Limiter) Average() int64 {
	m.Lock()
	defer m.Unlock()
	var sum int64
	for _, v := range m.windows {
		sum += v.count
	}
	if len(m.windows) > 0 {
		return sum / int64(len(m.windows))
	}
	return 0
}

func (m *Limiter) Validate() bool {
	m.Lock()
	defer m.Unlock()

	now := time.Now()
	windowsIndex := -1

	// 删除过期时间槽
	for k, v := range m.windows {
		// v.timestamp > now-m.width
		if v.timestamp.Add(m.width).After(now) {
			break
		}
		windowsIndex = k
	}
	if windowsIndex > -1 {
		m.windows = m.windows[windowsIndex+1:]
	}

	// 判断是否超出限制
	if m.maxCount > 0 {
		var sum int64
		for _, v := range m.windows {
			sum += v.count
		}
		if sum >= m.maxCount {
			return false
		}
	}

	// 写入窗口数组
	// timestamp > now-size
	if len(m.windows) > 0 && m.windows[len(m.windows)-1].timestamp.Add(m.size).After(now) {
		if m.perCount > 0 && m.windows[len(m.windows)-1].count >= m.perCount {
			return false
		}
		m.windows[len(m.windows)-1].count++
	} else {
		m.windows = append(m.windows, &timeSlot{
			timestamp: now,
			count:     1,
		})
	}

	return true

}

// RedisLimiter 使用redis存储限流器
// TODO 测试RedisLimiter
type RedisLimiter struct {
	width    time.Duration
	maxCount int64
	cli      *redis.Client
}

func NewRedisLimiter(width time.Duration, maxCount int64, cli *redis.Client) *RedisLimiter {
	return &RedisLimiter{
		width:    width,
		maxCount: maxCount,
		cli:      cli,
	}
}

// Validate redis存储
func (m *RedisLimiter) Validate(ctx context.Context, key string) bool {
	now := time.Now()
	min := strconv.Itoa(int(now.Add(-m.width).Unix()))
	max := strconv.Itoa(int(now.Unix()))
	member := uuid.New().String()

	// 删除过期记录
	m.cli.ZRemRangeByScore(ctx, key, "0", min)

	// 判断是否超限
	if count, _ := m.cli.ZCount(ctx, key, min, max).Uint64(); count >= uint64(m.maxCount) {
		return false
	}

	// 写入请求记录
	m.cli.ZAdd(ctx, key, &redis.Z{Score: float64(now.Unix()), Member: member})

	return true
}

// ValidateScript 使用redis脚本
func (m *RedisLimiter) ValidateScript(ctx context.Context, key string) bool {
	now := time.Now()
	min := strconv.Itoa(int(now.Add(-m.width).Unix()))
	max := strconv.Itoa(int(now.Unix()))
	member := uuid.New().String()
	nowStr := strconv.Itoa(int(now.Unix()))

	script := redis.NewScript(`
local key = KEYS[1]
local min = KEYS[2]
local max = KEYS[3]
local member = KEYS[4]
local maxCount = tonumber(KEYS[5])
local now = tonumber(KEYS[6])+0.0
redis.call('zremrangebyscore',key,'0',min)
local count = redis.call('zcount',key,min,max)
if count >= maxCount then
    return false
end
redis.call('zadd',key,now,member)
return true
`)
	b, err := script.Run(ctx, m.cli, []string{key, min, max, member, strconv.Itoa(int(m.maxCount)), nowStr}).Bool()
	if err != nil {
		return false
	}
	return b
}
