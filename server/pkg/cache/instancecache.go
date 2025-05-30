package cache

import (
	"context"
	"encoding"
	"log/slog"
	"sync"
	"time"

	"github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
)

type InstanceCacheFactory struct{}

func (f *InstanceCacheFactory) Create(name string) (interfaces.Cache, error) {
	return NewInstanceCache(name), nil
}

func NewInstanceCacheFactory() *InstanceCacheFactory {
	return &InstanceCacheFactory{}
}

type cacheItem struct {
	value     []byte
	expiresAt time.Time
}

type InstanceCache struct {
	data   map[string]cacheItem
	mutex  sync.RWMutex
	ticker *time.Ticker
	done   chan bool
	logger *slog.Logger
}

func NewInstanceCache(name string) *InstanceCache {
	cache := &InstanceCache{
		data:   make(map[string]cacheItem),
		done:   make(chan bool),
		logger: slog.Default().With("cache", name),
	}

	// Start cleanup goroutine that runs every 30 seconds
	cache.ticker = time.NewTicker(30 * time.Second)
	go cache.cleanup()

	return cache
}

func (c *InstanceCache) cleanup() {
	for {
		select {
		case <-c.ticker.C:
			c.mutex.Lock()
			now := time.Now()
			for key, item := range c.data {
				if !item.expiresAt.IsZero() && now.After(item.expiresAt) {
					delete(c.data, key)
				}
			}
			c.mutex.Unlock()
		case <-c.done:
			return
		}
	}
}

func (c *InstanceCache) Get(ctx context.Context, key string, value encoding.BinaryUnmarshaler) (bool, error) {
	c.mutex.RLock()
	item, exists := c.data[key]
	c.mutex.RUnlock()

	if !exists {
		return false, nil
	}

	// Check if expired (lazy cleanup)
	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		c.mutex.Lock()
		delete(c.data, key)
		c.mutex.Unlock()
		return false, nil
	}

	if err := value.UnmarshalBinary(item.value); err != nil {
		return false, err
	}

	return true, nil
}

func (c *InstanceCache) Set(ctx context.Context, key string, value encoding.BinaryMarshaler, ttl time.Duration) error {
	data, err := value.MarshalBinary()
	if err != nil {
		return err
	}

	expiresAt := time.Time{} // zero value means never expires
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	c.mutex.Lock()
	c.data[key] = cacheItem{
		value:     data,
		expiresAt: expiresAt,
	}
	c.mutex.Unlock()

	return nil
}

func (c *InstanceCache) SetIfNotExists(ctx context.Context, key string, setValue encoding.BinaryMarshaler, ttl time.Duration) (bool, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.data[key]; exists {
		return false, nil
	}

	data, err := setValue.MarshalBinary()
	if err != nil {
		return false, err
	}

	expiresAt := time.Time{} // zero value means never expires
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	c.data[key] = cacheItem{
		value:     data,
		expiresAt: expiresAt,
	}

	return true, nil
}

func (c *InstanceCache) Delete(ctx context.Context, key string) error {
	c.mutex.Lock()
	delete(c.data, key)
	c.mutex.Unlock()
	return nil
}

func (c *InstanceCache) Persist(ctx context.Context, key string) error {
	slog.Warn("Persist called on in-memory cache - this is a no-op for development",
		"key", key,
		"message", "persistence is not supported in instance cache")
	return nil
}

// Ensure InstanceCache implements the Cache interface
var _ interfaces.Cache = (*InstanceCache)(nil)
