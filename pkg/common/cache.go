package common

import (
	"context"
	"sync"
	"time"

	"github.com/CanftIn/gothafoss/pkg/redis"
)

type RedisCache struct {
	conn *redis.Conn
}

func NewRedisCache(ctx context.Context, addr string, password string) *RedisCache {
	r := &RedisCache{}
	r.conn = redis.New(ctx, addr, password)
	return r
}

func (r *RedisCache) Set(key string, value string) error {
	return r.conn.Set(key, value)
}

func (r *RedisCache) Delete(key string) error {
	return r.conn.Del(key)
}

func (r *RedisCache) SetAndExpire(key string, value string, expire time.Duration) error {
	return r.conn.SetAndExpire(key, value, expire)
}

func (r *RedisCache) Get(key string) (string, error) {
	return r.conn.Get(key)
}

func (r *RedisCache) GetRedisConn() *redis.Conn {
	return r.conn
}

type MemoryCache struct {
	cacheMap map[string]string
	sync.RWMutex
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		cacheMap: map[string]string{},
	}
}

func (m *MemoryCache) Set(key string, value string) error {
	m.Lock()
	m.cacheMap[key] = value
	m.Unlock()
	return nil
}

func (m *MemoryCache) SetAndExpire(key string, value string, expire time.Duration) error {
	m.Lock()
	m.cacheMap[key] = value
	m.Unlock()
	return nil
}

func (m *MemoryCache) Get(key string) (string, error) {
	m.RLock()
	defer m.RUnlock()
	return m.cacheMap[key], nil
}

func (m *MemoryCache) Delete(key string) error {
	m.Lock()
	delete(m.cacheMap, key)
	m.Unlock()
	return nil
}
