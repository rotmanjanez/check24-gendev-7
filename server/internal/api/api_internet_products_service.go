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
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/rotmanjanez/check24-gendev-7/config"
	"github.com/rotmanjanez/check24-gendev-7/internal/requestmanager"
	i "github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
	p "github.com/rotmanjanez/check24-gendev-7/pkg/provider"
)

// InternetProductsAPIService is a service that implements the logic for the InternetProductsAPIServicer
// This service should implement the business logic for every endpoint for the InternetProductsAPI API.
// Include any external packages or services that will be required by this service.
type InternetProductsAPIService struct {
	config *config.Config
	cache  i.Cache
	queue  i.Cache
	rc     *requestmanager.RequestCoordinator
}

const persistIndicator string = "indicator-persist"
const workInProgressIndicator string = "indicator-work-in-progress"

// NewInternetProductsAPIService creates a default api service
func NewInternetProductsAPIService(cfg *config.Config, cache i.Cache, queue i.Cache, providers []*p.ProviderConfig) *InternetProductsAPIService {
	return &InternetProductsAPIService{
		config: cfg,
		cache:  cache,
		queue:  queue,
		rc:     requestmanager.NewRequestCoordinator(providers),
	}
}

func (s *InternetProductsAPIService) ContinueInternetProductsQuery(ctx context.Context, cursor string) (ImplResponse, error) {
	// check if the cursor is a valid UUID
	if _, err := uuid.Parse(cursor); err != nil {
		return Response(http.StatusBadRequest, nil), errors.New("invalid cursor")
	}

	allProducts := m.SharedInternetProductsResponse{}
	foundAny := false

	for len(cursor) > 0 {
		products := new(m.InternetProductsResponse)
		exists, err := s.queue.Get(ctx, cursor, products)
		if err != nil {
			slog.Error("Error getting products from cache", "error", err)
			return Response(http.StatusInternalServerError, nil), err
		}
		if !exists {
			break
		}

		if products.NextCursor == workInProgressIndicator {
			if !foundAny {
				return Response(http.StatusAccepted, nil, map[string]string{"Retry-After": "3"}), nil // 3 seconds suggested
			}
			break
		}

		foundAny = true
		allProducts.Products = append(allProducts.Products, products.Products...)
		cursor = products.NextCursor
	}

	if !foundAny {
		return Response(http.StatusNotFound, nil), errors.New("products not found")
	}

	return Response(http.StatusOK, &m.InternetProductsResponse{
		Products:   allProducts.Products,
		NextCursor: cursor,
	}), nil
}

func (s *InternetProductsAPIService) processRequest(ctx context.Context, address m.Address, cursor string) {
	prods, errs := s.rc.Run(ctx, i.Request{
		Address: address,
	}, 10, 10)

	go func() {
		for err := range errs {
			slog.Error("Error fetching products", "error", err)
		}
	}()

	var products []m.InternetProduct
	current := cursor
	next := uuid.New().String()

	// set the initial cursor in the queue to indicate work in progress
	err := s.queue.Set(ctx, current, &m.InternetProductsResponse{NextCursor: workInProgressIndicator}, time.Duration(15*time.Minute))
	if err != nil {
		slog.Error("Error setting next cursor in cache", "error", err)
	}

	for prod := range prods {
		slog.Debug("Fetched product", "product", prod.Name, "cursor", current, "next", next)
		products = append(products, prod)

		slog.Debug("Adding product to queue", "cursor", current, "next", next)
		err := s.queue.Set(ctx, current, &m.InternetProductsResponse{
			Products:   []m.InternetProduct{prod},
			NextCursor: next,
		}, time.Duration(1*time.Hour))
		if err != nil {
			slog.Error("Error setting product in cache", "error", err)
			continue
		}

		// set next cursor in the queue to indicate work in progress for /internet-products/continue/{cursor}
		err = s.queue.Set(ctx, next, &m.InternetProductsResponse{NextCursor: workInProgressIndicator}, time.Duration(15*time.Minute))
		if err != nil {
			slog.Error("Error setting next cursor in cache", "error", err)
		}

		current = next
		next = uuid.New().String()
	}

	// add a final entry to the queue with the last cursor
	slog.Debug("Adding final product to queue", "cursor", current, "next", "")
	err = s.queue.Set(ctx, current, &m.InternetProductsResponse{
		Products:   []m.InternetProduct{},
		NextCursor: "",
	}, time.Duration(1*time.Hour))
	if err != nil {
		slog.Error("Error setting final product in cache", "error", err)
	}

	slog.Info("Fetched products", "count", len(products))

	ok, err := s.cache.SetIfNotExists(ctx, cursor, &m.SharedInternetProductsResponse{
		Products: products,
		Address:  address,
		Version:  m.INTERNET_PRODUCTS_RESPONSE_VERSION,
	}, time.Duration(5*time.Minute))

	if err != nil {
		slog.Error("Error setting products in cache", "error", err)
	}

	if !ok {
		// products already exist in the cache, need to persist
		existingProducts := new(m.SharedInternetProductsResponse)
		exists, err := s.cache.Get(ctx, cursor, existingProducts)
		if err != nil {
			slog.Error("Error getting products from cache", "error", err)
			return
		}
		if exists && existingProducts.Version != persistIndicator {
			// unlikely chance of uuid collision, but possible
			slog.Error("Products already exist in cache, persisting", "cursor", cursor)
		}
		// if it doesn't exist, it means that the persisted key expired in the meantime
		err = s.cache.Set(ctx, cursor, &m.SharedInternetProductsResponse{
			Products: products,
			Address:  address,
			Version:  m.INTERNET_PRODUCTS_RESPONSE_VERSION,
		}, i.KeepTTL)

		if err != nil {
			slog.Error("Error setting products in cache", "error", err)
			return
		}
	}
}

func (s *InternetProductsAPIService) InitiateInternetProductsQuery(ctx context.Context, address m.Address, providers []string) (ImplResponse, error) {
	cursor := uuid.New().String()

	go func() {
		bg := context.Background()
		bgWithTimeout, cancel := context.WithTimeout(bg, 60*time.Second)
		defer cancel()

		s.processRequest(bgWithTimeout, address, cursor)
	}()

	return Response(200, m.InternetProductsCursor{
		Version:    m.INTERNET_PRODUCTS_RESPONSE_VERSION,
		NextCursor: cursor,
	}), nil
}

func (s *InternetProductsAPIService) ShareInternetProducts(ctx context.Context, cursor string) (ImplResponse, error) {
	// check if the cursor is a valid UUID
	if _, err := uuid.Parse(cursor); err != nil {
		return Response(http.StatusBadRequest, nil), errors.New("invalid cursor")
	}

	value := m.SharedInternetProductsResponse{
		Version: persistIndicator,
	}
	ok, err := s.cache.SetIfNotExists(ctx, cursor, value, time.Duration(24*time.Hour))
	if err != nil {
		slog.Error("Error setting products in cache", "error", err)
		return Response(http.StatusInternalServerError, nil), err
	}

	if !ok {
		// products already exist in the cache, need to persist
		err := s.cache.Persist(ctx, cursor)
		if err != nil {
			slog.Error("Error persisting products in cache", "error", err)
			return Response(http.StatusInternalServerError, nil), err
		}
	}

	return Response(http.StatusOK, nil), nil
}

// GetSharedInternetProducts -
func (s *InternetProductsAPIService) GetSharedInternetProducts(ctx context.Context, cursor string) (ImplResponse, error) {
	if _, err := uuid.Parse(cursor); err != nil {
		return Response(http.StatusBadRequest, nil), errors.New("invalid cursor")
	}
	products := new(m.SharedInternetProductsResponse)
	exists, err := s.cache.Get(ctx, cursor, products)
	if err != nil {
		slog.Error("Error getting products from cache", "error", err)
		return Response(http.StatusInternalServerError, nil), err
	}

	if !exists || products.Version == persistIndicator {
		return Response(http.StatusNotFound, nil), errors.New("products not found")
	}

	return Response(http.StatusOK, products), nil
}
