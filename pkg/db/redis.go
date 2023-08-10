package db

import (
	"context"

	"github.com/CanftIn/gothafoss/pkg/redis"
)

func NewRedis(ctx context.Context, addr string, password string) *redis.Conn {
	return redis.New(ctx, addr, password)
}
