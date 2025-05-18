package utils

import (
	"barcode-generator-be/config"
	"time"

	"github.com/go-redis/redis/v8"
)

const BlacklistPrefix = "blacklist:"

type BlackList struct{}

var TokenBlacklist = &BlackList{}

func (bl *BlackList) Add(jti string, expirationTime time.Time) error {
	duration := time.Until(expirationTime)
	return config.RedisClient.Set(config.Ctx, BlacklistPrefix+jti, "true", duration).Err()
}

func (bl *BlackList) IsBlacklisted(jti string) (bool, error) {
	_, err := config.RedisClient.Get(config.Ctx, BlacklistPrefix+jti).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
