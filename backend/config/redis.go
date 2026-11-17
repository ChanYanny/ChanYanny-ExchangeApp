package config

import (
	"exchangeapp/global"
	"log"

	"github.com/go-redis/redis"
)

func initRedis() {
	// 初始化Redis连接
	redisClient := redis.NewClient(&redis.Options{
		Addr:     AppConfig.Redis.Addr,
		Password: AppConfig.Redis.Password,
		DB:       AppConfig.Redis.DB,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to establish Redis connection: %v", err)
	}
	global.RedisClient = redisClient
}
