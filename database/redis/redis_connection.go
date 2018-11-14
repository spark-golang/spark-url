package redis

import (
	"gopkg.in/redis.v5"
	"time"
)

var Demo *Redis

type Redis struct {
	RedisCluster *redis.Client
}

func (rc *Redis) Conn() {
	rc.RedisCluster = redis.NewClient(&redis.Options{
		Addr:        "127.0.0.1:6379",
		PoolSize:    200,
		MaxRetries:  8,
		ReadTimeout: 500 * time.Millisecond,
		IdleTimeout: 15 * time.Second,
	})
}
