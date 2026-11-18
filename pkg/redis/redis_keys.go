package redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrLockTTL Блокировка без срока жизни (ttl <= 0) не берётся: упавший владелец держал бы её вечно
var ErrLockTTL = errors.New("lock ttl must be positive")

// unlockScript Снимает блокировку, только если в ключе токен владельца: чужую (взятую после истечения ttl) не трогает
var unlockScript = redis.NewScript(`if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("DEL", KEYS[1]) end return 0`)

// Delete Удаляет ключи; возвращает число удалённых. Контекст — из NewHelper (свой — DeleteContext)
func (h *Helper[E]) Delete(keys ...string) (int64, error) {
	return h.DeleteContext(h.context, keys...)
}

// DeleteContext Delete с контекстом вызова
func (h *Helper[E]) DeleteContext(ctx context.Context, keys ...string) (int64, error) {
	if len(keys) == 0 {
		return 0, nil
	}

	start := time.Now()
	n, err := h.redisClient.Del(ctx, keys...).Result()
	h.setDebugInfo(ctx, "DEL", keys, time.Since(start))

	if err != nil {
		return 0, fmt.Errorf("failed to delete keys %v: %w", keys, err)
	}

	return n, nil
}

// Exists Есть ли ключ. Контекст — из NewHelper (свой — ExistsContext)
func (h *Helper[E]) Exists(key string) (bool, error) {
	return h.ExistsContext(h.context, key)
}

// ExistsContext Exists с контекстом вызова
func (h *Helper[E]) ExistsContext(ctx context.Context, key string) (bool, error) {
	start := time.Now()
	n, err := h.redisClient.Exists(ctx, key).Result()
	h.setDebugInfo(ctx, "EXISTS", key, time.Since(start))

	if err != nil {
		return false, fmt.Errorf("failed to check key %s: %w", key, err)
	}

	return n > 0, nil
}

// TryLock Берёт блокировку key на ttl (SET NX PX) со случайным токеном владельца. ok=false — блокировка занята.
// Снимать — Unlock(key, token): снимается только своя. Контекст — из NewHelper (свой — TryLockContext)
func (h *Helper[E]) TryLock(key string, ttl time.Duration) (token string, ok bool, err error) {
	return h.TryLockContext(h.context, key, ttl)
}

// TryLockContext TryLock с контекстом вызова
func (h *Helper[E]) TryLockContext(ctx context.Context, key string, ttl time.Duration) (token string, ok bool, err error) {
	if ttl <= 0 {
		return "", false, ErrLockTTL
	}

	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", false, fmt.Errorf("lock token: %w", err)
	}

	token = hex.EncodeToString(buf)

	start := time.Now()
	ok, err = h.redisClient.SetNX(ctx, key, token, ttl).Result()
	h.setDebugInfo(ctx, "LOCK", key, time.Since(start))

	if err != nil {
		return "", false, fmt.Errorf("failed to lock key %s: %w", key, err)
	}

	if !ok {
		return "", false, nil
	}

	return token, true, nil
}

// Unlock Снимает блокировку key, если она ещё принадлежит владельцу token; false — уже истекла или взята другим.
// Контекст — из NewHelper (свой — UnlockContext)
func (h *Helper[E]) Unlock(key, token string) (bool, error) {
	return h.UnlockContext(h.context, key, token)
}

// UnlockContext Unlock с контекстом вызова
func (h *Helper[E]) UnlockContext(ctx context.Context, key, token string) (bool, error) {
	start := time.Now()
	n, err := unlockScript.Run(ctx, h.redisClient, []string{key}, token).Int64()
	h.setDebugInfo(ctx, "UNLOCK", key, time.Since(start))

	if err != nil {
		return false, fmt.Errorf("failed to unlock key %s: %w", key, err)
	}

	return n == 1, nil
}
