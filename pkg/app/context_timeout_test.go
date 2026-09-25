package app

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/exgamer/gosdk-redis-core/pkg/config"
	"github.com/redis/go-redis/v9"
)

func TestNewRedisOptions_ContextTimeoutEnabled(t *testing.T) {
	if NewRedisOptions(&config.RedisConfig{}).ContextTimeoutEnabled {
		t.Error("unset REDIS_CONTEXT_TIMEOUT_ENABLED must keep go-redis' default (false)")
	}

	if !NewRedisOptions(&config.RedisConfig{ContextTimeoutEnabled: true}).ContextTimeoutEnabled {
		t.Error("REDIS_CONTEXT_TIMEOUT_ENABLED=true not passed to redis.Options")
	}
}

// A hung Redis, a 2s read timeout and a 100ms context deadline: only with
// the option does the deadline end the call.
func TestContextTimeoutEnabled_DeadlineEndsTheCall(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")

	if err != nil {
		t.Fatal(err)
	}

	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()

			if err != nil {
				return
			}

			defer conn.Close()
		}
	}()

	call := func(enabled bool) time.Duration {
		opts := NewRedisOptions(&config.RedisConfig{
			RedisHost:             ln.Addr().String(),
			ReadTimeoutMs:         2000,
			MaxRetries:            -1,
			ContextTimeoutEnabled: enabled,
		})
		opts.Protocol = 2
		opts.DisableIdentity = true
		client := redis.NewClient(opts)

		defer client.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		start := time.Now()
		_ = client.Get(ctx, "k").Err()

		return time.Since(start)
	}

	if took := call(true); took > time.Second {
		t.Errorf("enabled: call took %v, want it ended by the 100ms deadline", took)
	}

	if took := call(false); took < time.Second {
		t.Errorf("disabled: call took %v, want go-redis to wait for the 2s read timeout (the default behavior)", took)
	}
}
