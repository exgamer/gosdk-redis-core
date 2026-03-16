package app

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/exgamer/gosdk-core/pkg/app"
	config2 "github.com/exgamer/gosdk-core/pkg/config"
	"github.com/exgamer/gosdk-core/pkg/di"
	"github.com/exgamer/gosdk-redis-core/pkg/config"
	"github.com/redis/go-redis/v9"
	"strings"
)

const RedisKernelName = "redis"

type RedisKernel struct {
	redisClient *redis.Client
	redisConfig *config.RedisConfig
}

func (m *RedisKernel) Name() string {
	return RedisKernelName
}

func (m *RedisKernel) Init(a *app.App) error {
	err := m.initRedisConfig()

	if err != nil {
		return err
	}

	err = m.initRedisClient()

	if err != nil {
		return err
	}

	di.Register(a.Container, m.redisClient)

	return nil
}

func (m *RedisKernel) Start(a *app.App) error {

	return nil
}

func (m *RedisKernel) Stop(ctx context.Context) error {
	if m.redisClient == nil {
		return nil
	}

	return m.redisClient.Close()
}

func (m *RedisKernel) initRedisClient() error {
	if m.redisConfig.RedisUser == "" {
		m.redisConfig.RedisUser = "default"
	}
	fmt.Printf("redis host=%q db=%d pass_len=%d\n",
		m.redisConfig.RedisHost,
		m.redisConfig.RedisDb,
		len(strings.TrimSpace(m.redisConfig.RedisPassword)),
	)
	m.redisClient = redis.NewClient(&redis.Options{
		Addr:     m.redisConfig.RedisHost,
		DB:       m.redisConfig.RedisDb,
		Password: m.redisConfig.RedisPassword,
	})

	// Проверка подключения
	pong, err := m.redisClient.Ping(context.Background()).Result()

	if err != nil {
		return err
	}

	fmt.Println("Соединение с Redis:", pong)

	return nil
}

// InitRedisConfig Инициализация конфига редиса
func (m *RedisKernel) initRedisConfig() error {
	redisConfig := &config.RedisConfig{}
	err := config2.InitConfig(redisConfig)

	if err != nil {
		return err
	}

	m.redisConfig = redisConfig
	spew.Dump(redisConfig)

	return nil
}
