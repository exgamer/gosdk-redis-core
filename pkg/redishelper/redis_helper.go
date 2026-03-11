package redishelper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/exgamer/gosdk-core/pkg/debug"
	"github.com/exgamer/gosdk-core/pkg/helpers"
	"github.com/exgamer/gosdk-core/pkg/logger"
	"github.com/redis/go-redis/v9"
	"time"
)

// NewRedisHelper - Новый Хелпер для работы с редисом
func NewRedisHelper[E any](ctx context.Context, redisClient *redis.Client) *RedisHelper[E] {
	return &RedisHelper[E]{
		redisClient: redisClient,
		context:     ctx,
	}
}

// RedisHelper - Хелпер для работы с редисом
type RedisHelper[E any] struct {
	redisClient *redis.Client
	context     context.Context
}

// GetModel Возвращает значение по ключу
func (redisHelper *RedisHelper[E]) GetModel(key string) (*E, error) {
	start := time.Now()
	val, err := redisHelper.redisClient.Get(redisHelper.context, key).Result()
	execTime := time.Since(start)
	redisHelper.setDebugInfo("GET", key, execTime)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get model by key %s: %w", key, err)
	}

	if val == "" {
		return nil, nil
	}

	var result E
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal model by key %s: %w", key, err)
	}

	logger.Debug(redisHelper.context, fmt.Sprintf("CACHE HIT key=%s", key))

	return &result, nil
}

// SetModel Записывает значение по ключу
func (redisHelper *RedisHelper[E]) SetModel(key string, model E, ttl time.Duration) error {
	jsonModel, err := json.Marshal(model)
	if err != nil {
		return fmt.Errorf("failed to marshal model for key %s: %w", key, err)
	}

	start := time.Now()
	err = redisHelper.redisClient.Set(redisHelper.context, key, jsonModel, ttl).Err()
	execTime := time.Since(start)
	redisHelper.setDebugInfo("SET", key, execTime)

	if err != nil {
		return fmt.Errorf("failed to set model for key %s: %w", key, err)
	}

	logger.Debug(redisHelper.context, fmt.Sprintf("CACHE SET key=%s", key))

	return nil
}

// GetString Возвращает значение по ключу
func (redisHelper *RedisHelper[E]) GetString(key string) (string, error) {
	start := time.Now()
	val, err := redisHelper.redisClient.Get(redisHelper.context, key).Result()
	execTime := time.Since(start)
	redisHelper.setDebugInfo("GET", key, execTime)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}

		return "", fmt.Errorf("failed to get string by key %s: %w", key, err)
	}

	if val == "" {
		return "", nil
	}

	logger.Debug(redisHelper.context, fmt.Sprintf("CACHE HIT key=%s", key))

	return val, nil
}

// SetString Записывает значение по ключу
func (redisHelper *RedisHelper[E]) SetString(key string, value string, ttl time.Duration) error {
	start := time.Now()
	err := redisHelper.redisClient.Set(redisHelper.context, key, value, ttl).Err()
	execTime := time.Since(start)
	redisHelper.setDebugInfo("SET", key, execTime)

	if err != nil {
		return fmt.Errorf("failed to set string for key %s: %w", key, err)
	}

	logger.Debug(redisHelper.context, fmt.Sprintf("CACHE SET key=%s", key))

	return nil
}

// GetArray Возвращает массив по ключу
func (redisHelper *RedisHelper[E]) GetArray(key string) ([]E, error) {
	start := time.Now()
	val, err := redisHelper.redisClient.Get(redisHelper.context, key).Result()
	execTime := time.Since(start)
	redisHelper.setDebugInfo("GET", key, execTime)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get array by key %s: %w", key, err)
	}

	if val == "" {
		return nil, nil
	}

	var resultArr []E
	if err := json.Unmarshal([]byte(val), &resultArr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal array by key %s: %w", key, err)
	}

	logger.Debug(redisHelper.context, fmt.Sprintf("CACHE HIT key=%s", key))

	return resultArr, nil
}

// SetArray Записывает массив по ключу
func (redisHelper *RedisHelper[E]) SetArray(key string, models []E, ttl time.Duration) error {
	if len(models) == 0 {
		return nil
	}

	payload, err := json.Marshal(models)
	if err != nil {
		return fmt.Errorf("failed to marshal array for key %s: %w", key, err)
	}

	start := time.Now()
	err = redisHelper.redisClient.Set(redisHelper.context, key, payload, ttl).Err()
	execTime := time.Since(start)
	redisHelper.setDebugInfo("SET", key, execTime)

	if err != nil {
		return fmt.Errorf("failed to set array for key %s: %w", key, err)
	}

	logger.Debug(redisHelper.context, fmt.Sprintf("CACHE SET key=%s", key))

	return nil
}

// MSetModels - сохраняет несколько моделей в Redis
func (redisHelper *RedisHelper[E]) MSetModels(data map[string]E, exp time.Duration) error {
	if len(data) == 0 {
		return nil
	}

	pipe := redisHelper.redisClient.Pipeline()
	keys := make([]string, 0, len(data))

	for key, model := range data {
		jsonModel, err := json.Marshal(model)
		if err != nil {
			return fmt.Errorf("failed to marshal model for key %s: %w", key, err)
		}

		keys = append(keys, key)
		pipe.Set(redisHelper.context, key, jsonModel, exp)
	}

	start := time.Now()
	_, err := pipe.Exec(redisHelper.context)
	execTime := time.Since(start)
	redisHelper.setDebugInfo("MSET", keys, execTime)

	if err != nil {
		return fmt.Errorf("failed to set multiple keys in Redis: %w", err)
	}

	return nil
}

// MGetModels - Возвращает несколько моделей по ключам из Redis
func (redisHelper *RedisHelper[E]) MGetModels(keys []string) (map[string]E, error) {
	if len(keys) == 0 {
		return map[string]E{}, nil
	}

	start := time.Now()
	values, err := redisHelper.redisClient.MGet(redisHelper.context, keys...).Result()
	execTime := time.Since(start)

	if err != nil {
		return nil, fmt.Errorf("failed to get multiple keys from Redis: %w", err)
	}

	result := make(map[string]E, len(keys))

	for i, value := range values {
		if value == nil {
			continue
		}

		strValue, ok := value.(string)
		if !ok {
			bytesValue, ok := value.([]byte)
			if !ok {
				return nil, fmt.Errorf("unexpected value type for key %s: %T", keys[i], value)
			}

			strValue = string(bytesValue)
		}

		var model E
		if err := json.Unmarshal([]byte(strValue), &model); err != nil {
			return nil, fmt.Errorf("failed to unmarshal value for key %s: %w", keys[i], err)
		}

		result[keys[i]] = model
	}

	redisHelper.setDebugInfo("MGET", keys, execTime)

	return result, nil
}

func (redisHelper *RedisHelper[E]) setDebugInfo(operation string, keys interface{}, execTime time.Duration) {
	if redisHelper.context == nil {
		return
	}

	d := debug.GetDebugFromContext(redisHelper.context)

	if d == nil {

		return
	}

	statement := RedisStatement{}
	statement.Operation = operation
	statement.Keys = keys
	statement.Duration = execTime
	statement.Time = helpers.GetDurationAsString(execTime)

	d.Cat("cache")

	d.AddStatement("cache", execTime, []RedisStatement{statement})

	d.CalculateTotalTime()
}
