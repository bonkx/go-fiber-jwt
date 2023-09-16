package configs

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         context.Context
)

func ConnectRedis(config *Config) {
	ctx = context.TODO()

	RedisClient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})

	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	fmt.Println("✅ Redis client connected successfully...")
}
