package app

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/exgamer/gosdk-redis-core/pkg/config"
	"github.com/redis/go-redis/v9"
)

func TestNewRedisOptions_MapsConfig(t *testing.T) {
	opts := NewRedisOptions(&config.RedisConfig{
		RedisHost:       "redis:6379",
		RedisUser:       "app",
		RedisPassword:   "secret",
		RedisDb:         2,
		PoolSize:        20,
		IdleTimeout:     60,
		MaxConnLifetime: 600,
		DialTimeoutMs:   1000,
		ReadTimeoutMs:   500,
		WriteTimeoutMs:  400,
		MaxRetries:      1,
	})

	checks := []struct {
		name      string
		got, want any
	}{
		{"Addr", opts.Addr, "redis:6379"},
		{"Username", opts.Username, "app"},
		{"Password", opts.Password, "secret"},
		{"DB", opts.DB, 2},
		{"PoolSize", opts.PoolSize, 20},
		{"ConnMaxIdleTime", opts.ConnMaxIdleTime, time.Minute},
		{"ConnMaxLifetime", opts.ConnMaxLifetime, 10 * time.Minute},
		{"DialTimeout", opts.DialTimeout, time.Second},
		{"ReadTimeout", opts.ReadTimeout, 500 * time.Millisecond},
		{"WriteTimeout", opts.WriteTimeout, 400 * time.Millisecond},
		{"MaxRetries", opts.MaxRetries, 1},
	}

	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %v, want %v", c.name, c.got, c.want)
		}
	}
}

func TestNewRedisOptions_UnsetKeepsGoRedisDefaults(t *testing.T) {
	opts := NewRedisOptions(&config.RedisConfig{RedisHost: "redis:6379"})

	// Zero values are what go-redis turns into its own defaults.
	if opts.Username != "" || opts.ReadTimeout != 0 || opts.DialTimeout != 0 || opts.MaxRetries != 0 || opts.PoolSize != 0 {
		t.Errorf("unset config must leave go-redis defaults, got %+v", opts)
	}
}

func TestRedisConfig_MaskedHidesPassword(t *testing.T) {
	cfg := config.RedisConfig{RedisHost: "redis:6379", RedisPassword: "secret"}

	if got := cfg.Masked().RedisPassword; got != "***" {
		t.Errorf("masked password = %q", got)
	}

	if cfg.RedisPassword != "secret" {
		t.Error("Masked must not modify the original config")
	}

	if got := (config.RedisConfig{}).Masked().RedisPassword; got != "" {
		t.Errorf("empty password must stay empty, got %q", got)
	}
}

func TestReadTimeoutAndRetriesAreApplied(t *testing.T) {
	// A server that accepts connections and never answers — a hung Redis.
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

	opts := NewRedisOptions(&config.RedisConfig{RedisHost: ln.Addr().String(), ReadTimeoutMs: 200, MaxRetries: -1})
	opts.Protocol = 2
	opts.DisableIdentity = true
	client := redis.NewClient(opts)

	defer client.Close()

	start := time.Now()
	err = client.Get(context.Background(), "k").Err()
	took := time.Since(start)

	if err == nil || !strings.Contains(err.Error(), "timeout") {
		t.Fatalf("err = %v, want a read timeout", err)
	}

	if took > time.Second {
		t.Errorf("call took %v; with a 200ms read timeout and no retries it must fail fast (defaults would take ~12s)", took)
	}
}
