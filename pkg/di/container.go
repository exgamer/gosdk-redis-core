package di

import (
	"github.com/exgamer/gosdk-core/pkg/di"
	"github.com/redis/go-redis/v9"
)

// GetRedisClient возвращает клиент Redis
func GetRedisClient(c *di.Container) (*redis.Client, error) {
	client, err := di.Resolve[*redis.Client](c)

	if err != nil {
		return nil, err
	}

	return client, nil
}
