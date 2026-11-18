package redis

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

// NewHelper - Хелпер для работы с редисом
func NewHelper[E any](ctx context.Context, redisClient *redis.Client) *Helper[E] {
	return &Helper[E]{
		redisClient: redisClient,
		context:     ctx,
	}
}

// Helper - Хелпер для работы с редисом
type Helper[E any] struct {
	redisClient *redis.Client
	context     context.Context
}

// GetStruct Возвращает значение по ключу; контекст — из NewHelper (свой контекст вызова — GetStructContext)
func (h *Helper[E]) GetStruct(key string) (*E, error) {
	return h.GetStructContext(h.context, key)
}

// SetStruct Записывает значение по ключу; контекст — из NewHelper (свой контекст вызова — SetStructContext)
func (h *Helper[E]) SetStruct(key string, model E, ttl time.Duration) error {
	return h.SetStructContext(h.context, key, model, ttl)
}

// GetString Возвращает значение по ключу; контекст — из NewHelper (свой контекст вызова — GetStringContext)
func (h *Helper[E]) GetString(key string) (string, error) {
	return h.GetStringContext(h.context, key)
}

// SetString Записывает значение по ключу; контекст — из NewHelper (свой контекст вызова — SetStringContext)
func (h *Helper[E]) SetString(key string, value string, ttl time.Duration) error {
	return h.SetStringContext(h.context, key, value, ttl)
}

// GetArray Возвращает массив по ключу; контекст — из NewHelper (свой контекст вызова — GetArrayContext)
func (h *Helper[E]) GetArray(key string) ([]E, error) {
	return h.GetArrayContext(h.context, key)
}

// SetArray Записывает массив по ключу; контекст — из NewHelper (свой контекст вызова — SetArrayContext)
func (h *Helper[E]) SetArray(key string, models []E, ttl time.Duration) error {
	return h.SetArrayContext(h.context, key, models, ttl)
}

// MSetStruct сохраняет несколько структур в Redis; контекст — из NewHelper (свой контекст вызова — MSetStructContext)
func (h *Helper[E]) MSetStruct(data map[string]E, exp time.Duration) error {
	return h.MSetStructContext(h.context, data, exp)
}

// MGetStruct Возвращает несколько структур по ключам из Redis; контекст — из NewHelper (свой контекст вызова — MGetStructContext)
func (h *Helper[E]) MGetStruct(keys []string) (map[string]E, error) {
	return h.MGetStructContext(h.context, keys)
}

// GetStructContext GetStruct с контекстом вызова (дедлайн, отмена, отладка) вместо контекста NewHelper
func (h *Helper[E]) GetStructContext(ctx context.Context, key string) (*E, error) {
	start := time.Now()
	val, err := h.redisClient.Get(ctx, key).Result()
	execTime := time.Since(start)
	h.setDebugInfo(ctx, "GET", key, execTime)

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

	logger.Debug(ctx, fmt.Sprintf("CACHE HIT key=%s", key))

	return &result, nil
}

// SetStructContext SetStruct с контекстом вызова (дедлайн, отмена, отладка) вместо контекста NewHelper
func (h *Helper[E]) SetStructContext(ctx context.Context, key string, model E, ttl time.Duration) error {
	jsonModel, err := json.Marshal(model)
	if err != nil {
		return fmt.Errorf("failed to marshal model for key %s: %w", key, err)
	}

	start := time.Now()
	err = h.redisClient.Set(ctx, key, jsonModel, ttl).Err()
	execTime := time.Since(start)
	h.setDebugInfo(ctx, "SET", key, execTime)

	if err != nil {
		return fmt.Errorf("failed to set model for key %s: %w", key, err)
	}

	logger.Debug(ctx, fmt.Sprintf("CACHE SET key=%s", key))

	return nil
}

// GetStringContext GetString с контекстом вызова (дедлайн, отмена, отладка) вместо контекста NewHelper
func (h *Helper[E]) GetStringContext(ctx context.Context, key string) (string, error) {
	start := time.Now()
	val, err := h.redisClient.Get(ctx, key).Result()
	execTime := time.Since(start)
	h.setDebugInfo(ctx, "GET", key, execTime)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}

		return "", fmt.Errorf("failed to get string by key %s: %w", key, err)
	}

	if val == "" {
		return "", nil
	}

	logger.Debug(ctx, fmt.Sprintf("CACHE HIT key=%s", key))

	return val, nil
}

// SetStringContext SetString с контекстом вызова (дедлайн, отмена, отладка) вместо контекста NewHelper
func (h *Helper[E]) SetStringContext(ctx context.Context, key string, value string, ttl time.Duration) error {
	start := time.Now()
	err := h.redisClient.Set(ctx, key, value, ttl).Err()
	execTime := time.Since(start)
	h.setDebugInfo(ctx, "SET", key, execTime)

	if err != nil {
		return fmt.Errorf("failed to set string for key %s: %w", key, err)
	}

	logger.Debug(ctx, fmt.Sprintf("CACHE SET key=%s", key))

	return nil
}

// GetArrayContext GetArray с контекстом вызова (дедлайн, отмена, отладка) вместо контекста NewHelper
func (h *Helper[E]) GetArrayContext(ctx context.Context, key string) ([]E, error) {
	start := time.Now()
	val, err := h.redisClient.Get(ctx, key).Result()
	execTime := time.Since(start)
	h.setDebugInfo(ctx, "GET", key, execTime)

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

	logger.Debug(ctx, fmt.Sprintf("CACHE HIT key=%s", key))

	return resultArr, nil
}

// SetArrayContext SetArray с контекстом вызова (дедлайн, отмена, отладка) вместо контекста NewHelper
func (h *Helper[E]) SetArrayContext(ctx context.Context, key string, models []E, ttl time.Duration) error {
	if models == nil {
		models = []E{} // пустой срез — «[]» (закешированный пустой результат), не «null»
	}

	payload, err := json.Marshal(models)
	if err != nil {
		return fmt.Errorf("failed to marshal array for key %s: %w", key, err)
	}

	start := time.Now()
	err = h.redisClient.Set(ctx, key, payload, ttl).Err()
	execTime := time.Since(start)
	h.setDebugInfo(ctx, "SET", key, execTime)

	if err != nil {
		return fmt.Errorf("failed to set array for key %s: %w", key, err)
	}

	logger.Debug(ctx, fmt.Sprintf("CACHE SET key=%s", key))

	return nil
}

// MSetStructContext MSetStruct с контекстом вызова (дедлайн, отмена, отладка) вместо контекста NewHelper
func (h *Helper[E]) MSetStructContext(ctx context.Context, data map[string]E, exp time.Duration) error {
	if len(data) == 0 {
		return nil
	}

	pipe := h.redisClient.Pipeline()
	keys := make([]string, 0, len(data))

	for key, model := range data {
		jsonModel, err := json.Marshal(model)
		if err != nil {
			return fmt.Errorf("failed to marshal model for key %s: %w", key, err)
		}

		keys = append(keys, key)
		pipe.Set(ctx, key, jsonModel, exp)
	}

	start := time.Now()
	_, err := pipe.Exec(ctx)
	execTime := time.Since(start)
	h.setDebugInfo(ctx, "MSET", keys, execTime)

	if err != nil {
		return fmt.Errorf("failed to set multiple keys in Redis: %w", err)
	}

	return nil
}

// MGetStructContext MGetStruct с контекстом вызова (дедлайн, отмена, отладка) вместо контекста NewHelper
func (h *Helper[E]) MGetStructContext(ctx context.Context, keys []string) (map[string]E, error) {
	if len(keys) == 0 {
		return map[string]E{}, nil
	}

	start := time.Now()
	values, err := h.redisClient.MGet(ctx, keys...).Result()
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

	h.setDebugInfo(ctx, "MGET", keys, execTime)

	return result, nil
}

func (h *Helper[E]) setDebugInfo(ctx context.Context, operation string, keys interface{}, execTime time.Duration) {
	if ctx == nil {
		return
	}

	d := debug.GetDebugFromContext(ctx)

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
