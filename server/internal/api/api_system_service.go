/*
 * CHECK24 GenDev 7 API
 *
 * API for the 7th CHECK24 GenDev challenge providing product offerings from five different internet providers
 *
 * API version: dev
 */

package api

import (
	"context"

	"github.com/rotmanjanez/check24-gendev-7/config"
	"github.com/rotmanjanez/check24-gendev-7/pkg/models"
)

// SystemAPIService is a service that implements the logic for the SystemAPIServicer
// This service should implement the business logic for every endpoint for the SystemAPI API.
// Include any external packages or services that will be required by this service.
type SystemAPIService struct {
	config *config.Config
}

// NewSystemAPIService creates a default api service
func NewSystemAPIService(cfg *config.Config) *SystemAPIService {
	if cfg == nil {
		panic("config cannot be nil")
	}

	return &SystemAPIService{
		config: cfg,
	}
}

// GetVersion - Version information endpoint
func (s *SystemAPIService) GetVersion(ctx context.Context) (ImplResponse, error) {
	return Response(200, models.Version{
		Version:    s.config.Version,
		BuildDate:  s.config.BuildDate,
		CommitHash: s.config.CommitHash,
	}), nil
}
