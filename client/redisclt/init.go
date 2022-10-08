package redisclt

import (
	"github.com/OverClockTeam/OC-WordsInSongs-BE/conf"
	"github.com/go-redis/redis"
)

func InitRedis(redisSettings conf.RedisSettings) {
	RedisClient = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     redisSettings.Address + ":" + redisSettings.Port,
		Password: redisSettings.Password,
		DB:       0,
	})
}
