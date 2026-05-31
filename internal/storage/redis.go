package storage

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"

	"cosmetologybotliza/internal/config"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedis(cfg *config.Config) *RedisClient {

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
		DB:   0,
	})

	ctx := context.Background()

	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Redis connected")

	return &RedisClient{
		Client: rdb,
		Ctx:    ctx,
	}
}
