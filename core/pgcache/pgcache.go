// Package cache @Description  TODO
// @Author  	 jiangyang
// @Created  	 2024/8/5 下午5:43
package pgcache

import (
	"database/sql"
	"fmt"
	"github.com/SpectatorNan/gorm-zero/gormc/config/pg"
	"time"

	_ "github.com/lib/pq"
)

var ErrTryLockFailed = fmt.Errorf("failed to acquire lock")

// Cache struct to hold the DB connection
// INSERT OR UPDATE ON cache will cleanup expired rows
// 功能
// 1. 带过期时间的分布式锁
// 2. 带过期时间的键值对缓存
// 3. 自动清理过期数据
type Cache struct {
	db *sql.DB
}

// NewCache creates a new Cache instance
func NewCache(conf pg.PgSql) (*Cache, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s host=%s port=%d password=%s sslmode=%s", conf.Username, conf.Dbname, conf.Path, conf.Port, conf.Password, conf.SslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &Cache{db: db}, nil
}

// Set sets a value in the cache with a TTL
func (c *Cache) Set(key string, value string, ttl time.Duration) error {
	ttlInterval := fmt.Sprintf("%d seconds", int(ttl.Seconds()))
	_, err := c.db.Exec("SELECT insert_cache($1, $2, $3)", key, value, ttlInterval)
	return err
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) (string, error) {
	var value string
	err := c.db.QueryRow("SELECT value FROM cache WHERE key = $1 AND created_at + ttl > NOW()", key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Cache miss
		}
		return "", err
	}
	return value, nil
}

// Delete removes a value from the cache
func (c *Cache) Delete(key string) error {
	_, err := c.db.Exec("DELETE FROM cache WHERE key = $1", key)
	return err
}

func (c *Cache) TryLock(lockName string, ttl time.Duration) error {
	ttlInterval := fmt.Sprintf("%d seconds", int(ttl.Seconds()))
	result, err := c.db.Exec("INSERT INTO cache (key, value,ttl) VALUES ($1,$2,$3) ON CONFLICT (key) DO NOTHING", lockName, "", ttlInterval)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrTryLockFailed
	}

	return nil
}

// Unlock TryLock必须与Unlock配对使用，否则会造成解了他人锁的问题
func (c *Cache) Unlock(lockName string) error {
	_, err := c.db.Exec("DELETE FROM cache WHERE key = $1", lockName)
	return err
}

// Cleanup manually triggers cache cleanup
func (c *Cache) Cleanup() error {
	_, err := c.db.Exec("SELECT cleanup_cache()")
	return err
}

// Close closes the database connection
func (c *Cache) Close() error {
	return c.db.Close()
}
