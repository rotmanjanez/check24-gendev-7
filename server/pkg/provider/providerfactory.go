package provider

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/rotmanjanez/check24-gendev-7/config"
	i "github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
)

type ProviderFactory func(options map[string]interface{}, cache i.Cache, logger *slog.Logger) (i.ProviderAdapter, error)

// Global registry of backend factories
var registry = make(map[string]ProviderFactory)

// Register a new backend type
func RegisterProvider(name string, factory ProviderFactory) {
	registry[name] = factory
}

// ProviderConfig holds settings and HTTP client for one provider
// and controls concurrency and retry behavior.
type ProviderConfig struct {
	Adapter         i.ProviderAdapter
	Client          *http.Client
	RetryCount      int
	Timeout         time.Duration
	ConcurrentLimit int           // max parallel HTTP requests
	Semaphore       chan struct{} // throttles concurrent calls
	BackoffInterval time.Duration // base wait for retries/throttling
}

// NewProviderConfig constructs a ProviderConfig with concurrency control
func NewProviderConfig(adapter i.ProviderAdapter, retries int, timeout time.Duration, maxConcurrent int, backoff time.Duration) *ProviderConfig {
	if maxConcurrent <= 0 {
		maxConcurrent = 1
	}
	return &ProviderConfig{
		Adapter: adapter,
		Client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 32 {
					return http.ErrUseLastResponse
				}
				return nil
			},
			Timeout: timeout,
		},
		RetryCount:      retries,
		Timeout:         timeout,
		ConcurrentLimit: maxConcurrent,
		Semaphore:       make(chan struct{}, maxConcurrent),
		BackoffInterval: backoff,
	}
}

// Create backends from config
func CreateProviders(cacheFactory i.CacheFactory, cfg *config.Config) ([]*ProviderConfig, error) {
	var providers []*ProviderConfig

	for name, backendCfg := range cfg.Backends {
		if !backendCfg.Enabled {
			continue
		}

		cache, err := cacheFactory.Create(name)
		if err != nil {
			return nil, fmt.Errorf("failed to create cache for provider %s: %w", name, err)
		}

		provider, err := CreateProvider(name, cache, backendCfg.Options)
		if err != nil {
			return nil, fmt.Errorf("failed to create provider %s: %w", name, err)
		}
		providers = append(providers, NewProviderConfig(
			provider,
			backendCfg.Retries,
			backendCfg.Timeout,
			backendCfg.MaxConcurrent,
			backendCfg.Backoff,
		))
	}

	return providers, nil
}

func CreateProvider(name string, cache i.Cache, options map[string]interface{}) (i.ProviderAdapter, error) {
	factory, exists := registry[name]
	if !exists {
		return nil, fmt.Errorf("unknown provider: %s", name)
	}

	logger := slog.With("provider", name)

	provider, err := factory(options, cache, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create backend %s: %w", name, err)
	}
	slog.Info("Created provider", "name", name)
	return provider, nil
}
