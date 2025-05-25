package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// override local variables with build flags
var buildDate = time.Now()

var commitHash = ""

type Config struct {
	// Version is the version of the application.
	// default: dev
	Version string `json:"version"`

	// Address is the address on which the application will listen.
	// It can be an IP address or a hostname.
	// default: localhost
	Address string `json:"address"`

	// Port is the port on which the application will listen.
	// default: 8080
	Port uint `json:"port"`

	MaxConcurrentRequests uint `json:"maxConcurrentRequests"`

	RequestBufferSize uint `json:"requestBufferSize"`

	// BuildDate is the date when the application was built.
	// It is set at build time using the -X flag.
	// default: empty
	BuildDate time.Time

	// CommitHash is the commit hash of the application.
	// It is set at build time using the -X flag.
	// default: empty
	CommitHash string

	Redis *redis.Options `json:"redis"`

	Backends map[string]BackendConfig `json:"backends"`
}

type BackendConfig struct {
	Enabled       bool                   `json:"enabled"`
	Retries       int                    `json:"retries"`
	Timeout       time.Duration          `json:"timeout"`
	MaxConcurrent int                    `json:"maxConcurrent"`
	Backoff       time.Duration          `json:"backoff"`
	Options       map[string]interface{} `json:"options"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	for key, backend := range config.Backends {
		backend.Timeout = backend.Timeout * time.Millisecond
		backend.Backoff = backend.Backoff * time.Millisecond
		config.Backends[key] = backend
	}

	config.BuildDate = buildDate
	config.CommitHash = commitHash

	return &config, nil
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}
