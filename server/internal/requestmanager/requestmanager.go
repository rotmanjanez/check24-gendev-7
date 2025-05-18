package requestmanager

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	i "github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
	p "github.com/rotmanjanez/check24-gendev-7/pkg/provider"
)

// requestContext tracks inflight work per provider
// via its own WaitGroup for follow-up requests.
type requestContext struct {
	config *p.ProviderConfig
	wg     sync.WaitGroup
}

// RequestCoordinator dispatches a Request across providers
// and collects results from fresh channels per call.
type RequestCoordinator struct {
	providers []*p.ProviderConfig
}

// NewRequestCoordinator returns a coordinator over the given provider configs
func NewRequestCoordinator(cfgs []*p.ProviderConfig) *RequestCoordinator {
	return &RequestCoordinator{providers: cfgs}
}

// Run executes req on all providers, returning new channels for responses and errors
func (c *RequestCoordinator) Run(ctx context.Context, req i.Request, respBuf, errBuf int) (<-chan m.InternetProduct, <-chan error) {
	responses := make(chan m.InternetProduct, respBuf)
	errors := make(chan error, errBuf)
	var wg sync.WaitGroup

	// dispatch per provider
	for _, cfg := range c.providers {
		wg.Add(1)
		rctx := &requestContext{config: cfg}
		go func(pc *p.ProviderConfig, rc *requestContext) {
			defer wg.Done()
			// initial preparation and follow-ups
			c.dispatchProvider(ctx, pc, req, responses, errors, &rc.wg)
			// wait for all follow-up requests to finish
			rc.wg.Wait()
		}(cfg, rctx)
	}

	// close channels when all work completes
	go func() {
		wg.Wait()
		close(responses)
		close(errors)
	}()

	return responses, errors
}

// dispatchProvider handles a single provider's preparation
// and issues follow-up requests via the provided wait group.
func (c *RequestCoordinator) dispatchProvider(
	ctx context.Context,
	cfg *p.ProviderConfig,
	initialReq i.Request,
	responses chan<- m.InternetProduct,
	errors chan<- error,
	wg *sync.WaitGroup,
) {
	parsedResp, err := cfg.Adapter.PrepareRequest(ctx, initialReq)
	if err != nil {
		errors <- err
		return
	}
	// handle initial parse and spawn follow-ups
	c.handleParsed(ctx, cfg, parsedResp, initialReq, responses, errors, wg)
}

// handleParsed emits products and schedules follow-up requests
func (c *RequestCoordinator) handleParsed(
	ctx context.Context,
	cfg *p.ProviderConfig,
	parsed i.ParsedResponse,
	orig i.Request,
	responses chan<- m.InternetProduct,
	errors chan<- error,
	wg *sync.WaitGroup,
) {
	// emit parsed products
	for _, p := range parsed.InternetProducts {
		// check if product date offered is zero and set to now
		if p.DateOffered.IsZero() {
			p.DateOffered = time.Now()
		}
		p = m.CanonicalizeInternetProduct(p)

		err := m.AssertInternetProductRequired(p)
		if err != nil {
			slog.Warn("InternetProduct missing fields", "provider", cfg.Adapter.Name(), "product", p, "error", err)
			errors <- fmt.Errorf("invalid product from %s: %w", cfg.Adapter.Name(), err)
			continue
		}

		err = m.AssertInternetProductConstraints(p)
		if err != nil {
			slog.Warn("Invalid InternetProduct constraints", "provider", cfg.Adapter.Name(), "product", p, "error", err)
			errors <- fmt.Errorf("invalid product from %s: %w", cfg.Adapter.Name(), err)
			continue
		}
		// canonicalize the product before sending
		responses <- p
	}

	// issue follow-up requests in parallel
	for _, follow := range parsed.Requests {
		if follow.Request == nil {
			slog.Warn("No follow-up request provided, skipping", "followUp", follow)
			continue // skip if no follow-up request
		}

		if follow.Callback == nil {
			follow.Callback = cfg.Adapter
		}
		follow.Request = follow.Request.WithContext(ctx)

		wg.Add(1)
		// each follow-up decrements on completion
		go func(r i.Response) {
			defer wg.Done()
			c.dispatchRequest(ctx, cfg, r, orig, responses, errors, wg)
		}(i.Response{InitialRequestData: orig, Request: follow})
	}
}

// dispatchRequest executes a single HTTP call with retry and backoff
func (c *RequestCoordinator) dispatchRequest(
	ctx context.Context,
	cfg *p.ProviderConfig,
	respWrapper i.Response,
	orig i.Request,
	responses chan<- m.InternetProduct,
	errors chan<- error,
	wg *sync.WaitGroup,
) {
	for attempt := 0; attempt <= cfg.RetryCount; attempt++ {
		if attempt > 0 {
			slog.Info("Retrying request", "adapter", cfg.Adapter.Name(), "attempt", attempt)
		}
		cfg.Semaphore <- struct{}{}
		resp, err := cfg.Client.Do(respWrapper.Request.Request)
		<-cfg.Semaphore
		if err != nil {
			slog.Debug("Error executing request", "adapter", cfg.Adapter.Name(), "error", err, "attempt", attempt)
			if attempt == cfg.RetryCount {
				slog.Debug("Max retries reached, giving up", "adapter", cfg.Adapter.Name(), "error", err)
				errors <- err
			}
			time.Sleep(cfg.BackoffInterval)
			continue
		}

		switch resp.StatusCode {
		case http.StatusOK:
			respWrapper.HTTPResponse = resp
			parsed, perr := cfg.Adapter.ParseResponse(ctx, respWrapper)
			resp.Body.Close()
			if perr != nil {
				errors <- perr
			} else {
				c.handleParsed(ctx, cfg, parsed, orig, responses, errors, wg)
			}
			return

		case http.StatusTooManyRequests:
			resp.Body.Close()
			time.Sleep(cfg.BackoffInterval)
			continue // retry after backoff

		default:
			slog.Debug("Unexpected response from provider",
				"adapter", cfg.Adapter.Name(),
				"statusCode", resp.StatusCode,
				"request", respWrapper.Request.Request.URL.String(),
			)
			resp.Body.Close()
			if attempt == cfg.RetryCount {
				errors <- fmt.Errorf("unexpected responses after mutliple retries from %s: %s", cfg.Adapter.Name(), resp.Status)
			}
			continue
		}
	}
}
