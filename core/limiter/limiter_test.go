// Package limiter @Description  TODO
// @Author  	 jiangyang
// @Created  	 2023/1/3 15:48
package limiter_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/geekeryy/api-hub/core/limiter"
	"github.com/stretchr/testify/assert"
)

func TestLimiter(t *testing.T) {
	loop := 4
	perCount := 5
	size := time.Millisecond
	width := size * time.Duration(loop)
	l := limiter.NewLimiter(width, size, int64(perCount), int64(loop*perCount))
	for i := 0; i < loop; i++ {
		var counter int64
		var wg sync.WaitGroup
		for j := 0; j < 1000; j++ {
			wg.Add(1)
			go func() {
				if l.Validate() {
					atomic.AddInt64(&counter, 1)
				}
				wg.Done()
			}()
		}
		wg.Wait()
		assert.Equal(t, int64(perCount), counter)
		time.Sleep(size)
	}
}
