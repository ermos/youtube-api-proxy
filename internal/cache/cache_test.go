package cache

import (
	"context"
	"testing"
	"time"
)

func TestCacheHitMissExpiry(t *testing.T) {
	ctx := context.Background()
	c := New(1)

	var got []string

	// miss on empty
	if found, _ := c.Get(ctx, "k", &got); found {
		t.Fatal("expected miss on empty cache")
	}

	// hit after set
	if err := c.Set(ctx, "k", []string{"a", "b"}); err != nil {
		t.Fatalf("set: %v", err)
	}
	found, err := c.Get(ctx, "k", &got)
	if err != nil || !found {
		t.Fatalf("expected hit, got found=%v err=%v", found, err)
	}
	if len(got) != 2 || got[0] != "a" {
		t.Fatalf("unexpected value: %v", got)
	}

	// expiry: force the entry into the past
	c.mu.Lock()
	c.items["k"] = entry{data: []byte(`["a"]`), expiresAt: time.Now().Add(-time.Second)}
	c.mu.Unlock()
	if found, _ := c.Get(ctx, "k", &got); found {
		t.Fatal("expected miss after expiry")
	}
}
