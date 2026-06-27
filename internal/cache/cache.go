package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Cache is a tiny in-process TTL cache. It lives in the Go binary's memory —
// no external store. Expired entries are dropped lazily on read.
// ponytail: lazy eviction on Get; add a sweeper goroutine only if memory grows unbounded.
type Cache struct {
	mu    sync.RWMutex
	items map[string]entry
	ttl   time.Duration
}

type entry struct {
	data      []byte
	expiresAt time.Time
}

func New(ttlSeconds int) *Cache {
	return &Cache{
		items: make(map[string]entry),
		ttl:   time.Duration(ttlSeconds) * time.Second,
	}
}

func (c *Cache) Get(_ context.Context, key string, dest interface{}) (bool, error) {
	c.mu.RLock()
	e, ok := c.items[key]
	c.mu.RUnlock()

	if !ok || time.Now().After(e.expiresAt) {
		return false, nil
	}

	if err := json.Unmarshal(e.data, dest); err != nil {
		return false, fmt.Errorf("json unmarshal error: %w", err)
	}
	return true, nil
}

func (c *Cache) Set(_ context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	c.mu.Lock()
	c.items[key] = entry{data: data, expiresAt: time.Now().Add(c.ttl)}
	c.mu.Unlock()
	return nil
}

func VideosKey(channelID string) string {
	return fmt.Sprintf("youtube:videos:%s", channelID)
}

func PlaylistsKey(channelID string) string {
	return fmt.Sprintf("youtube:playlists:%s", channelID)
}
