package redis

import (
	"context"
	"fmt"
	"testing"
)

func TestRedis(t *testing.T) {
	ctx := context.Background()
	rdb := New(ctx, "127.0.0.1", "123456")
	fmt.Println("Client:", rdb.client)
	rdb.Set("key", "value")

	val, _ := rdb.Get("key")
	fmt.Println(val)
}
