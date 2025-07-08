package database

import (
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST" validate:"required"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type RedisDB struct {
	Client *redis.Client
	Config *RedisConfig
}

func NewRedisClient(cfg *RedisConfig) *RedisDB {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &RedisDB{
		Client: redisClient,
		Config: cfg,
	}
}
