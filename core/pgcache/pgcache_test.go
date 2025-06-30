// Package pgcache @Description  TODO
// @Author  	 jiangyang
// @Created  	 2024/8/5 下午5:44
package pgcache_test

import (
	"fmt"
	"github.com/SpectatorNan/gorm-zero/gormc/config/pg"
	"github.com/geekeryy/api-hub/core/pgcache"
	"log"
	"testing"
	"time"
)

var conf = pg.PgSql{
	Username: "dev",
	Dbname:   "postgres",
	Path:     "localhost",
	Port:     5432,
	Password: "123456",
	SslMode:  "disable",
}

func TestLock(t *testing.T) {
	var cache, err = pgcache.NewCache(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer cache.Close()

	t.Run("bench lock", func(t *testing.T) {
		benchLockID := "mylock-bench"
		for i := 0; i < 10; i++ {
			go func() {
				for j := 0; j < 1000; j++ {
					if err = cache.TryLock(benchLockID, 3*time.Second); err != nil {
						if err == pgcache.ErrTryLockFailed {
							continue
						} else {
							t.Error(err)
						}
					}
					fmt.Println(benchLockID, i, j, "Lock acquired", time.Now().Unix())
				}
			}()
		}
		time.Sleep(time.Second * 6)

	})

	t.Run("lock unlock", func(t *testing.T) {
		lockID := "mylock"
		counter := 0
		for j := 0; j < 10; j++ {
			err = cache.TryLock(lockID, 10*time.Second)
			if err != nil {
				log.Fatal("Failed to acquire lock:", err)
			}
			counter++
			err = cache.Unlock(lockID)
			if err != nil {
				log.Fatal("Failed to release lock:", err)
			}
		}
		if counter != 10 {
			t.Error("Failed to release lock")
		}
	})

}

func TestPgCache(t *testing.T) {
	var cache, err = pgcache.NewCache(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer cache.Close()

	err = cache.Set("mykey", "myvalue", 3*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	value, err := cache.Get("mykey")
	if err != nil {
		log.Fatal(err)
	}
	if value != "myvalue" {
		t.Error("Failed to get value")
	}

	time.Sleep(time.Second * 3)

	value, err = cache.Get("mykey")
	if err != nil {
		log.Fatal(err)
	}
	if value == "myvalue" {
		t.Error("Failed to get value")
	}

}
