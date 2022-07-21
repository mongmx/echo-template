package infra

import (
	"context"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

// RedisConfig is a configuration for redis.
type RedisConfig struct {
	Host string
	Port string
	Password string
	DB int
}

// NewRedis creates new connection to redis and return the connection
func NewRedis(cfg RedisConfig) (*redis.Client, error) {
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return redisClient, nil
}

func redisDB() int {
	d, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return 1
	}
	return d
}
