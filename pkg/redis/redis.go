package redis

import (
	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
)

type RedisCache struct {
	*redis.Client
}

type Config struct {
	Address  *string `mapstructure:"address"`
	Port     *string `mapstructure:"port"`
	Password string  `mapstructure:"password"`
}

func NewRedisCache(config *Config) *RedisCache {

	rdb := redis.NewClient(&redis.Options{
		Addr:     *config.Address + ":" + *config.Port,
		Password: config.Password,
	})

	log.Info().Msgf("Create Redis Client Successful!")

	return &RedisCache{
		rdb,
	}
}
