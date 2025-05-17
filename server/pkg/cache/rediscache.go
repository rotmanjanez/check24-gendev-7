package cache

import (
	"context"
	"encoding"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
)

type RedisCacheFactory struct {
	client *redis.Client
}

func NewRedisCacheFactory(options *redis.Options) (*RedisCacheFactory, error) {
	if options == nil {
		return nil, fmt.Errorf("redis options cannot be nil")
	}
	if options.Addr == "" {
		return nil, fmt.Errorf("redis address is required")
	}
	if options.ClientName == "" {
		return nil, fmt.Errorf("redis client name is required")
	}

	client := redis.NewClient(options)
	if client == nil {
		return nil, fmt.Errorf("failed to create redis client")
	}

	return &RedisCacheFactory{
		client: client,
	}, nil
}

func (f *RedisCacheFactory) Create(name string) (interfaces.Cache, error) {
	if f.client == nil {
		return nil, fmt.Errorf("redis client is not initialized")
	}
	if name == "" {
		return nil, fmt.Errorf("redis cache name is empty")
	}

	return NewRedisCache(name, f.client), nil
}

type RedisCache struct {
	prefix string
	client *redis.Client
}

func NewRedisCache(prefix string, client *redis.Client) *RedisCache {
	return &RedisCache{
		prefix: prefix + ":",
		client: client,
	}
}

func (r *RedisCache) Get(ctx context.Context, key string, value encoding.BinaryUnmarshaler) (bool, error) {
	key = r.prefix + key
	slog.Debug("Getting value from Redis", "key", key)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		slog.Debug("Key not found in Redis", "key", key)
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("error getting value from redis: %w", err)
	}
	slog.Debug("Value found in Redis", "key", key, "value", val)
	err = value.UnmarshalBinary([]byte(val))
	if err != nil {
		return false, fmt.Errorf("error unmarshalling value from redis: %w", err)
	}
	slog.Debug("Successfully retrieved value from Redis", "key", key, "value", val)
	return true, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value encoding.BinaryMarshaler, ttl time.Duration) error {
	key = r.prefix + key
	slog.Debug("Setting value in Redis", "key", key, "value", value, "ttl", ttl)
	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		slog.Error("Error setting value in Redis", "key", key, "error", err)
		return err
	}
	return nil
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	key = r.prefix + key
	slog.Debug("Deleting value from Redis", "key", key)
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		slog.Error("Error deleting value from Redis", "key", key, "error", err)
		return err
	}
	return nil
}

func (r *RedisCache) Persist(ctx context.Context, key string) error {
	// persist an existing key
	key = r.prefix + key
	slog.Debug("Persisting value in Redis", "key", key)
	err := r.client.Persist(ctx, key).Err()
	if err != nil {
		slog.Error("Error persisting value in Redis", "key", key, "error", err)
		return err
	}
	slog.Debug("Successfully persisted value in Redis", "key", key)
	return nil
}

func (r *RedisCache) SetIfNotExists(ctx context.Context, key string, setValue encoding.BinaryMarshaler, ttl time.Duration) (bool, error) {
	key = r.prefix + key
	slog.Debug("Setting value in Redis if not exists", "key", key, "value", setValue, "ttl", ttl)
	val, err := r.client.SetNX(ctx, key, setValue, ttl).Result()
	if err != nil {
		slog.Error("Error setting value in Redis if not exists", "key", key, "error", err)
		return false, err
	}
	slog.Debug("Successfully set value in Redis if not exists", "key", key, "value", val)
	return val, nil
}

var _ interfaces.Cache = (*RedisCache)(nil)
