package cache

import "time"

type Cache interface {
	Set(key string, value string) error
	Delete(key string) error
	SetAndExpire(key string, value string, expire time.Duration) error
	Get(key string) (string, error)
}
