package backend

import (
	"context"
	"fmt"
	"log"
	"time"

	env "rip/pkg/settings"

	"github.com/redis/go-redis/v9"
)

// InitializeRedis инициализирует клиент Redis
func InitializeRedis() (*redis.Client, error) {
	addr, password, db, err := env.FromEnvRedis()
	if err != nil {
		return nil, fmt.Errorf("failed to get Redis configuration from environment: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	log.Println("Redis client initialized successfully")
	return client, nil
}
