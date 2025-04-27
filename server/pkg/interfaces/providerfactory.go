package interfaces

import (
	"fmt"
	"log/slog"

	"github.com/rotmanjanez/check24-gendev-7/config"
)

type ProviderFactory func(options map[string]interface{}, logger *slog.Logger) (ProviderAdapter, error)

// Global registry of backend factories
var registry = make(map[string]ProviderFactory)

// Register a new backend type
func RegisterProvider(name string, factory ProviderFactory) {
	registry[name] = factory
}

// Create backends from config
func CreateProviders(cfg *config.Config) ([]ProviderAdapter, error) {
	var providers []ProviderAdapter

	for name, backendCfg := range cfg.Backends {
		if !backendCfg.Enabled {
			continue
		}

		provider, err := CreateProvider(name, backendCfg.Options)
		if err != nil {
			return nil, fmt.Errorf("failed to create provider %s: %w", name, err)
		}
		providers = append(providers, provider)
	}

	return providers, nil
}

func CreateProvider(name string, options map[string]interface{}) (ProviderAdapter, error) {
	factory, exists := registry[name]
	if !exists {
		return nil, fmt.Errorf("unknown provider: %s", name)
	}

	logger := slog.With("provider", name)

	provider, err := factory(options, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create backend %s: %w", name, err)
	}
	slog.Info("Created provider", "name", name)
	return provider, nil
}
