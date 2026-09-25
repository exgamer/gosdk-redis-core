package app

import (
	"context"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/exgamer/gosdk-core/pkg/app"
	config2 "github.com/exgamer/gosdk-core/pkg/config"
	"github.com/exgamer/gosdk-core/pkg/di"
	"github.com/exgamer/gosdk-redis-core/pkg/config"
	"github.com/redis/go-redis/v9"
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
	m.redisClient = redis.NewClient(NewRedisOptions(m.redisConfig))

	// Проверка подключения
	pong, err := m.redisClient.Ping(context.Background()).Result()

	if err != nil {
		return err
	}

	fmt.Println("Соединение с Redis:", pong)

	return nil
}

// NewRedisOptions собирает redis.Options из конфига. Нулевые значения не
// передаются, чтобы действовали значения go-redis по умолчанию.
func NewRedisOptions(cfg *config.RedisConfig) *redis.Options {
	return &redis.Options{
		Addr:            cfg.RedisHost,
		Username:        cfg.RedisUser,
		Password:        cfg.RedisPassword,
		DB:              cfg.RedisDb,
		PoolSize:        cfg.PoolSize,
		ConnMaxIdleTime: seconds(cfg.IdleTimeout),
		ConnMaxLifetime: seconds(cfg.MaxConnLifetime),
		DialTimeout:     milliseconds(cfg.DialTimeoutMs),
		ReadTimeout:     milliseconds(cfg.ReadTimeoutMs),
		WriteTimeout:    milliseconds(cfg.WriteTimeoutMs),
		MaxRetries:      cfg.MaxRetries,

		ContextTimeoutEnabled: cfg.ContextTimeoutEnabled,
	}
}

func seconds(v int) time.Duration {
	if v <= 0 {
		return 0
	}

	return time.Duration(v) * time.Second
}

func milliseconds(v int) time.Duration {
	if v <= 0 {
		return 0
	}

	return time.Duration(v) * time.Millisecond
}

// InitRedisConfig Инициализация конфига редиса
func (m *RedisKernel) initRedisConfig() error {
	redisConfig := &config.RedisConfig{}
	err := config2.InitConfig(redisConfig)

	if err != nil {
		return err
	}

	m.redisConfig = redisConfig
	masked := redisConfig.Masked()
	spew.Dump(&masked)

	return nil
}
