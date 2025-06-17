// Package limiter @Description  TODO
// @Author  	 jiangyang
// @Created  	 2023/1/3 15:48
package limiter_test

import (
	"github.com/geekeryy/api-hub/core/limiter"
	"log"
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	l := limiter.NewLimiter(time.Second*10, time.Second, 1, 5)
	for i := 0; i < 36; i++ {
		log.Println(l.Validate())
		time.Sleep(time.Millisecond * 500)
	}
	time.Sleep(time.Second * 11)
	if l.Counter() > 0 {
		l.Reset()
		log.Println("reset")
	}
	for i := 0; i < 36; i++ {
		log.Println(l.Validate())
		time.Sleep(time.Millisecond * 500)
	}
}
