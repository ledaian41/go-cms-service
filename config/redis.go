package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func InitRedisClient() *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: Env.RedisHost,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Error connecting to Redis", err)
		return nil
	}
	log.Println("Connected to Redis")

	return &RedisClient{
		Client: rdb,
		Ctx:    context.Background(),
	}
}
