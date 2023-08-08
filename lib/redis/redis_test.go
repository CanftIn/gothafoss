package redis

import (
	"fmt"
	"testing"
)

func TestRedis(t *testing.T) {
	rdb := New("127.0.0.1", "123456")
	fmt.Println("Client:", rdb.client)
	err := rdb.Set("key", "value")
	if err != nil {
		fmt.Println("connect success")
	} else {
		fmt.Println("connect failed")
	}

	val, err := rdb.GetString("key")
	fmt.Println(val)
}
