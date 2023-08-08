package redis

import (
	"context"
	"time"

	rdb "github.com/redis/go-redis/v9"
)

type Field struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type Conn struct {
	client *rdb.Client
	ctx    context.Context
}

func New(ctx context.Context, addr string, password string) *Conn {
	c := &Conn{}
	c.client = rdb.NewClient(&rdb.Options{
		Addr:       addr,
		MaxRetries: 3,
		Password:   password,
	})
	c.ctx = ctx
	return c
}

func (rc *Conn) Ping() (string, error) {
	return rc.client.Ping(rc.ctx).Result()
}

func (rc *Conn) Set(key string, value interface{}) error {
	return rc.client.Set(rc.ctx, key, value, 0).Err()
}

func (rc *Conn) SetAndExpire(key string, value interface{}, expire time.Duration) error {
	return rc.client.Set(rc.ctx, key, value, expire).Err()
}

func (rc *Conn) Get(key string) (string, error) {
	val, err := rc.client.Get(rc.ctx, key).Result()
	if err == rdb.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (rc *Conn) Del(key string) error {
	return rc.client.Del(rc.ctx, key).Err()
}

func (rc *Conn) LLen(key string) (int64, error) {
	val, err := rc.client.LLen(rc.ctx, key).Result()
	if err == rdb.Nil {
		return 0, nil
	}
	return val, err
}

func (rc *Conn) LRange(key string, start, stop int64) ([]string, error) {
	val, err := rc.client.LRange(rc.ctx, key, start, stop).Result()
	if err == rdb.Nil {
		return nil, nil
	}
	return val, err
}
