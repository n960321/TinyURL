package redis

import "github.com/go-redis/redis"

type RedisCache struct {
	*redis.Client
}

func NewRedisCache() *RedisCache {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	
	return &RedisCache{
		rdb,
	}
}
