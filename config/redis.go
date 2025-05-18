package config

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return nil
}
