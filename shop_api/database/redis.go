package database

import (
	"context"
	"fmt"
	"log"
	"shop_api/config"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func InitRedis(cfg *config.RedisConfig) error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}

	log.Println("Redis connection established")
	return nil
}

func GetRedis() *redis.Client {
	return RDB
}
