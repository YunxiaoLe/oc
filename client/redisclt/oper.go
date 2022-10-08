package redisclt

import (
	"github.com/go-redis/redis"
	"time"
)

var RedisClient *redis.Client

func Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return RedisClient.Set(key, value, expiration)
}
func Get(key string) *redis.StringCmd {
	return RedisClient.Get(key)
}
func Del(keys string) *redis.IntCmd {
	return RedisClient.Del(keys)
}
func Exists(keys string) *redis.IntCmd {
	return RedisClient.Exists(keys)
}
