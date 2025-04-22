package config

import (
	"fmt"
	"time"
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
}

func NewConfig() *Config {
	return &Config{
		Version:    "dev",
		Address:    "localhost",
		Port:       8080,
		BuildDate:  buildDate,
		CommitHash: commitHash,
	}
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}
