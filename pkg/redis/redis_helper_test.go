package redis

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

type item struct {
	ID int `json:"id"`
}

func newTestClient(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	t.Helper()

	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = client.Close() })

	return mr, client
}

func TestSetArray_EmptyWritesEmptyArray(t *testing.T) {
	mr, client := newTestClient(t)
	h := NewHelper[item](context.Background(), client)

	for _, models := range [][]item{nil, {}} {
		mr.FlushAll()

		if err := h.SetArray("k", models, time.Minute); err != nil {
			t.Fatalf("SetArray: %v", err)
		}

		if got, _ := mr.Get("k"); got != "[]" {
			t.Fatalf("stored %q, want []", got)
		}

		arr, err := h.GetArray("k")
		if err != nil || arr == nil || len(arr) != 0 {
			t.Fatalf("GetArray = %v, %v — want an empty non-nil slice (a cached empty result)", arr, err)
		}
	}

	arr, err := h.GetArray("missing")
	if err != nil || arr != nil {
		t.Fatalf("missing key: %v, %v — want nil (a cache miss)", arr, err)
	}
}

func TestHelper_ContextVariantsUseCallContext(t *testing.T) {
	_, client := newTestClient(t)

	// the helper's own context is already cancelled: the old methods fail, the Context variants use the call context
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()

	h := NewHelper[item](cancelled, client)
	if err := h.SetStruct("k", item{ID: 1}, time.Minute); !errors.Is(err, context.Canceled) {
		t.Fatalf("SetStruct with a cancelled helper context: %v, want context.Canceled", err)
	}

	ctx := context.Background()
	if err := h.SetStructContext(ctx, "k", item{ID: 1}, time.Minute); err != nil {
		t.Fatalf("SetStructContext: %v", err)
	}

	got, err := h.GetStructContext(ctx, "k")
	if err != nil || got == nil || got.ID != 1 {
		t.Fatalf("GetStructContext = %+v, %v", got, err)
	}

	if err := h.SetStringContext(ctx, "s", "v", time.Minute); err != nil {
		t.Fatal(err)
	}

	if s, err := h.GetStringContext(ctx, "s"); err != nil || s != "v" {
		t.Fatalf("GetStringContext = %q, %v", s, err)
	}

	if err := h.SetArrayContext(ctx, "a", []item{{ID: 2}}, time.Minute); err != nil {
		t.Fatal(err)
	}

	if a, err := h.GetArrayContext(ctx, "a"); err != nil || len(a) != 1 || a[0].ID != 2 {
		t.Fatalf("GetArrayContext = %v, %v", a, err)
	}

	if err := h.MSetStructContext(ctx, map[string]item{"m1": {ID: 3}}, time.Minute); err != nil {
		t.Fatal(err)
	}

	if m, err := h.MGetStructContext(ctx, []string{"m1", "none"}); err != nil || len(m) != 1 || m["m1"].ID != 3 {
		t.Fatalf("MGetStructContext = %v, %v", m, err)
	}
}

func TestDeleteExists(t *testing.T) {
	mr, client := newTestClient(t)
	h := NewHelper[item](context.Background(), client)

	_ = mr.Set("a", "1")
	_ = mr.Set("b", "2")

	if ok, err := h.Exists("a"); err != nil || !ok {
		t.Fatalf("Exists(a) = %v, %v", ok, err)
	}

	n, err := h.Delete("a", "b", "c")
	if err != nil || n != 2 {
		t.Fatalf("Delete = %d, %v — want 2", n, err)
	}

	if ok, err := h.ExistsContext(context.Background(), "a"); err != nil || ok {
		t.Fatalf("Exists after Delete = %v, %v", ok, err)
	}

	if n, err := h.Delete(); err != nil || n != 0 {
		t.Fatalf("Delete() = %d, %v", n, err)
	}
}

func TestTryLockUnlock_OwnerToken(t *testing.T) {
	mr, client := newTestClient(t)
	h := NewHelper[any](context.Background(), client)

	token, ok, err := h.TryLock("lock", time.Second)
	if err != nil || !ok || token == "" {
		t.Fatalf("TryLock = %q, %v, %v", token, ok, err)
	}

	if _, ok, err := h.TryLock("lock", time.Second); err != nil || ok {
		t.Fatalf("second TryLock = %v, %v — want busy", ok, err)
	}

	if ok, err := h.Unlock("lock", "foreign"); err != nil || ok {
		t.Fatalf("Unlock with a foreign token = %v, %v — want false", ok, err)
	}

	if !mr.Exists("lock") {
		t.Fatal("a foreign token must not remove the lock")
	}

	// the lock expired and another owner took it: the old owner must not remove it
	mr.FastForward(2 * time.Second)

	token2, ok, err := h.TryLock("lock", time.Second)
	if err != nil || !ok {
		t.Fatalf("TryLock after expiry = %v, %v", ok, err)
	}

	if ok, _ := h.Unlock("lock", token); ok {
		t.Fatal("the expired owner removed another owner's lock")
	}

	if ok, err := h.UnlockContext(context.Background(), "lock", token2); err != nil || !ok {
		t.Fatalf("Unlock by the owner = %v, %v", ok, err)
	}

	if _, _, err := h.TryLock("lock", 0); !errors.Is(err, ErrLockTTL) {
		t.Fatalf("TryLock without ttl: %v, want ErrLockTTL", err)
	}
}

func TestTryLock_OneWinnerConcurrently(t *testing.T) {
	_, client := newTestClient(t)
	h := NewHelper[any](context.Background(), client)

	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		wins int
	)

	for i := 0; i < 16; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if _, ok, err := h.TryLock("lock", time.Minute); err == nil && ok {
				mu.Lock()
				wins++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if wins != 1 {
		t.Fatalf("winners = %d, want 1", wins)
	}
}
