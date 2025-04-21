package interfaces

import (
	"context"
	"encoding"
	"time"
)

const KeepTTL = -1

type Cache interface {
	Get(ctx context.Context, key string, value encoding.BinaryUnmarshaler) (bool, error)
	Set(ctx context.Context, key string, value encoding.BinaryMarshaler, ttl time.Duration) error
	SetIfNotExists(ctx context.Context, key string, setValue encoding.BinaryMarshaler, ttl time.Duration) (bool, error)
	Delete(ctx context.Context, key string) error
	Persist(ctx context.Context, key string) error
}

type CacheFactory interface {
	Create(name string) (Cache, error)
}
