package app

import (
	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Addr string
	DB   int
}

func NewRedisClient(config RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: "", // no password set
		DB:       config.DB,
	})

	return rdb
}
