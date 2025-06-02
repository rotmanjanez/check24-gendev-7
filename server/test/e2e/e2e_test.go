package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/rotmanjanez/check24-gendev-7/config"
	"github.com/rotmanjanez/check24-gendev-7/internal/api"
	"github.com/rotmanjanez/check24-gendev-7/pkg/cache"
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
	p "github.com/rotmanjanez/check24-gendev-7/pkg/provider"
	"github.com/rotmanjanez/check24-gendev-7/providers/byteme"
	"github.com/rotmanjanez/check24-gendev-7/providers/pingperfect"
	"github.com/rotmanjanez/check24-gendev-7/providers/servusspeed"
	"github.com/rotmanjanez/check24-gendev-7/providers/verbyndich"
	"github.com/rotmanjanez/check24-gendev-7/providers/webwunder"
)

// E2ETestSuite contains all dependencies for end-to-end testing
type E2ETestSuite struct {
	server      *httptest.Server
	client      *http.Client
	redisClient *redis.Client
	cfg         *config.Config
}

// setupE2ETest creates a complete test environment with real providers and Redis
func setupE2ETest(t *testing.T) *E2ETestSuite {
	// Create test configuration (only use valid fields)
	cfg := &config.Config{
		Version:           "test-v1.0.0",
		BuildDate:         time.Now(),
		CommitHash:        "test-commit",
		Address:           "localhost",
		Port:              8081, // Use a test port
		UseInProcessCache: true,
	}

	// Create cache instances (pass a name and the client)
	mainCache := cache.NewInstanceCache("e2e-main")
	queueCache := cache.NewInstanceCache("e2e-queue")

	// Setup logger for testing (silent)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // Only show errors in tests
	}))

	// Create provider configurations with mock servers for controlled testing
	providers := createTestProviders(t, logger)

	// Create API services
	internetService := api.NewInternetProductsAPIService(cfg, mainCache, queueCache, providers)
	systemService := api.NewSystemAPIService(cfg)
	healthService := api.NewHealthAPIService()

	// Create controllers
	internetController := api.NewInternetProductsAPIController(internetService)
	systemController := api.NewSystemAPIController(systemService)
	healthController := api.NewHealthAPIController(healthService)

	// Create router
	router := api.NewRouter(internetController, systemController, healthController)

	// Create test server
	server := httptest.NewServer(router)

	return &E2ETestSuite{
		server: server,
		client: &http.Client{Timeout: 30 * time.Second},
		cfg:    cfg,
	}
}

// tearDownE2ETest cleans up the test environment
func (suite *E2ETestSuite) tearDown(t *testing.T) {
	if suite.server != nil {
		suite.server.Close()
	}
	if suite.redisClient != nil {
		ctx := context.Background()
		suite.redisClient.FlushDB(ctx)
		suite.redisClient.Close()
	}
}

// createTestProviders creates provider configurations for testing
func createTestProviders(t *testing.T, logger *slog.Logger) []*p.ProviderConfig {
	byteMeServer := createByteMeTestServer()
	webWunderServer := createWebWunderTestServer()
	t.Logf("WebWunder mock server URL: %s", webWunderServer.URL)
	verbynDichServer := createVerbynDichTestServer()
	servusSpeedServer := createServusSpeedTestServer()
	pingPerfectServer := createPingPerfectTestServer()

	t.Cleanup(func() {
		byteMeServer.Close()
		webWunderServer.Close()
		verbynDichServer.Close()
		servusSpeedServer.Close()
		pingPerfectServer.Close()
	})

	return []*p.ProviderConfig{
		p.NewProviderConfig(
			byteme.NewByteMeAdapter(byteMeServer.URL, "test-api-key", logger),
			3, 5*time.Second, 1, 500*time.Millisecond,
		),
		p.NewProviderConfig(
			webwunder.NewWebWunderAdapter(
				"test-api-key",
				webWunderServer.URL+"/endpunkte/soap/ws",
				"http://spring.io/guides/gs-producing-web-service/legacyGetInternetOffers",
				"http://webwunder.gendev7.check24.fun/offerservice",
				"http://schemas.xmlsoap.org/soap/envelope/",
				logger,
			),
			3, 5*time.Second, 1, 500*time.Millisecond,
		),
		p.NewProviderConfig(
			verbyndich.NewVerbynDichAdapter(verbynDichServer.URL+"/check24/data", "dummy-api-key", 1, logger),
			3, 5*time.Second, 1, 500*time.Millisecond,
		),
		p.NewProviderConfig(
			servusspeed.NewServusSpeedAdapter("test-username", "test-password", servusSpeedServer.URL+"/api/external/product-details", cache.NewInstanceCache("servusspeed-e2e"), 5*time.Minute, logger),
			3, 5*time.Second, 1, 500*time.Millisecond,
		),
		p.NewProviderConfig(
			pingperfect.NewPingPerfectAdapter(pingPerfectServer.URL, "test-client-id", "test-signature-secret", logger),
			3, 5*time.Second, 1, 500*time.Millisecond,
		),
	}
}

// TestE2E_FullWorkflow tests the complete end-to-end workflow
func TestE2E_FullWorkflow(t *testing.T) {
	suite := setupE2ETest(t)
	defer suite.tearDown(t)

	address := m.Address{
		Street:      "Marienplatz",
		HouseNumber: "1",
		City:        "München",
		PostalCode:  "80331",
		CountryCode: "DE",
	}

	// Step 1: Initiate query
	t.Log("Step 1: Initiating internet products query...")
	cursor := suite.initiateQuery(t, address)
	if cursor == "" {
		t.Fatal("Failed to initiate query - no cursor returned")
	}
	t.Logf("Received cursor: %s", cursor)

	// Step 2: Poll for results
	t.Log("Step 2: Polling for results...")
	products := suite.pollForResults(t, cursor, 30*time.Second)
	if len(products) == 0 {
		t.Fatal("No products returned from polling")
	}
	t.Logf("Retrieved %d products", len(products))

	// Validate products have required fields
	suite.validateProducts(t, products)

	// Step 3: Share results
	t.Log("Step 3: Sharing results...")
	suite.shareResults(t, cursor)

	// Step 4: Retrieve shared results
	t.Log("Step 4: Retrieving shared results...")
	sharedProducts := suite.getSharedResults(t, cursor)
	if len(sharedProducts.Products) != len(products) {
		t.Errorf("Shared products count mismatch: expected %d, got %d",
			len(products), len(sharedProducts.Products))
	}

	t.Log("✅ Full workflow completed successfully!")
}

// initiateQuery starts a new internet products query
func (suite *E2ETestSuite) initiateQuery(t *testing.T, address m.Address) string {
	requestBody, err := json.Marshal(address)
	if err != nil {
		t.Fatalf("Failed to marshal address: %v", err)
	}

	resp, err := suite.client.Post(
		suite.server.URL+"/internet-products",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		t.Fatalf("Failed to initiate query: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Initiate query failed with status %d", resp.StatusCode)
	}

	var cursorResp m.InternetProductsCursor
	if err := json.NewDecoder(resp.Body).Decode(&cursorResp); err != nil {
		t.Fatalf("Failed to decode cursor response: %v", err)
	}

	return cursorResp.NextCursor
}

// pollForResults polls the continue endpoint until all results are available
func (suite *E2ETestSuite) pollForResults(t *testing.T, cursor string, timeout time.Duration) []m.InternetProduct {
	var allProducts []m.InternetProduct
	startTime := time.Now()
	currentCursor := cursor

	for time.Since(startTime) < timeout {
		resp, err := suite.client.Get(
			suite.server.URL + "/internet-products/continue?cursor=" + currentCursor,
		)
		if err != nil {
			t.Fatalf("Failed to poll results: %v", err)
		}

		if resp.StatusCode == http.StatusAccepted {
			// No more results yet, wait and try again
			resp.Body.Close()
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			t.Fatalf("Poll failed with status %d", resp.StatusCode)
		}

		var productsResp m.InternetProductsResponse
		if err := json.NewDecoder(resp.Body).Decode(&productsResp); err != nil {
			resp.Body.Close()
			t.Fatalf("Failed to decode products response: %v", err)
		}
		resp.Body.Close()

		allProducts = append(allProducts, productsResp.Products...)

		// Check if we've reached the end
		if productsResp.NextCursor == "" {
			break
		}

		currentCursor = productsResp.NextCursor
		time.Sleep(100 * time.Millisecond) // Small delay between polls
	}

	return allProducts
}

// validateProducts ensures all products have required fields
func (suite *E2ETestSuite) validateProducts(t *testing.T, products []m.InternetProduct) {
	for i, product := range products {
		if product.Id == "" {
			t.Errorf("Product %d missing Id", i)
		}
		if product.Name == "" {
			t.Errorf("Product %d missing Name", i)
		}
		if product.Provider == "" {
			t.Errorf("Product %d missing Provider", i)
		}
		if product.ProductInfo.Speed <= 0 {
			t.Errorf("Product %d has invalid speed: %d", i, product.ProductInfo.Speed)
		}
		if product.Pricing.MonthlyCostInCent <= 0 {
			t.Errorf("Product %d has invalid monthly cost: %d", i, product.Pricing.MonthlyCostInCent)
		}
	}
}

// shareResults shares the results using the cursor
func (suite *E2ETestSuite) shareResults(t *testing.T, cursor string) {
	resp, err := suite.client.Post(
		suite.server.URL+"/internet-products/share/"+cursor,
		"application/json",
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to share results: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Share results failed with status %d", resp.StatusCode)
	}
}

// getSharedResults retrieves shared results
func (suite *E2ETestSuite) getSharedResults(t *testing.T, cursor string) m.SharedInternetProductsResponse {
	resp, err := suite.client.Get(
		suite.server.URL + "/internet-products/share/" + cursor,
	)
	if err != nil {
		t.Fatalf("Failed to get shared results: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Get shared results failed with status %d", resp.StatusCode)
	}

	var sharedResp m.SharedInternetProductsResponse
	if err := json.NewDecoder(resp.Body).Decode(&sharedResp); err != nil {
		t.Fatalf("Failed to decode shared response: %v", err)
	}

	return sharedResp
}

// TestE2E_ErrorScenarios tests various error conditions
func TestE2E_ErrorScenarios(t *testing.T) {
	suite := setupE2ETest(t)
	defer suite.tearDown(t)

	t.Run("InvalidAddress", func(t *testing.T) {
		invalidAddress := m.Address{
			Street:      "", // Missing required field
			HouseNumber: "1",
			City:        "München",
			PostalCode:  "80331",
			CountryCode: "DE",
		}

		requestBody, _ := json.Marshal(invalidAddress)
		resp, err := suite.client.Post(
			suite.server.URL+"/internet-products",
			"application/json",
			bytes.NewBuffer(requestBody),
		)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("Expected 422, got %d", resp.StatusCode)
		}
	})

	t.Run("InvalidCursor", func(t *testing.T) {
		resp, err := suite.client.Get(
			suite.server.URL + "/internet-products/continue?cursor=invalid-cursor",
		)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("NonExistentCursor", func(t *testing.T) {
		// Use a valid UUID format but non-existent cursor
		fakeCursor := "550e8400-e29b-41d4-a716-446655440000"
		resp, err := suite.client.Get(
			suite.server.URL + "/internet-products/continue?cursor=" + fakeCursor,
		)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", resp.StatusCode)
		}
	})

	t.Run("ShareNonExistentCursor", func(t *testing.T) {
		fakeCursor := "550e8400-e29b-41d4-a716-446655440000"
		resp, err := suite.client.Post(
			suite.server.URL+"/internet-products/share/"+fakeCursor,
			"application/json",
			nil,
		)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		// Should succeed (creates placeholder for sharing)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GetNonSharedResults", func(t *testing.T) {
		fakeCursor := "550e8400-e29b-41d4-a716-446655440000"
		resp, err := suite.client.Get(
			suite.server.URL + "/internet-products/share/" + fakeCursor,
		)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", resp.StatusCode)
		}
	})
}

// TestE2E_ConcurrentRequests tests handling of concurrent requests
func TestE2E_ConcurrentRequests(t *testing.T) {
	suite := setupE2ETest(t)
	defer suite.tearDown(t)

	address := m.Address{
		Street:      "Marienplatz",
		HouseNumber: "1",
		City:        "München",
		PostalCode:  "80331",
		CountryCode: "DE",
	}

	const numConcurrentRequests = 5
	results := make(chan string, numConcurrentRequests)
	errors := make(chan error, numConcurrentRequests)

	// Start multiple concurrent requests
	for i := 0; i < numConcurrentRequests; i++ {
		go func(requestID int) {
			cursor := suite.initiateQuery(t, address)
			if cursor == "" {
				errors <- fmt.Errorf("request %d: no cursor returned", requestID)
				return
			}
			results <- cursor
		}(i)
	}

	// Collect results
	var cursors []string
	for i := 0; i < numConcurrentRequests; i++ {
		select {
		case cursor := <-results:
			cursors = append(cursors, cursor)
		case err := <-errors:
			t.Errorf("Concurrent request failed: %v", err)
		case <-time.After(10 * time.Second):
			t.Fatalf("Timeout waiting for concurrent requests")
		}
	}

	// Verify all cursors are unique
	cursorSet := make(map[string]bool)
	for _, cursor := range cursors {
		if cursorSet[cursor] {
			t.Errorf("Duplicate cursor found: %s", cursor)
		}
		cursorSet[cursor] = true
	}

	t.Logf("Successfully handled %d concurrent requests", len(cursors))
}

// TestE2E_ProviderTimeout tests handling of provider timeouts
func TestE2E_ProviderTimeout(t *testing.T) {
	suite := setupE2ETest(t)
	defer suite.tearDown(t)

	// This test would require modifying provider configurations to include
	// timeouts or slow responses. For now, we test normal workflow.
	address := m.Address{
		Street:      "Timeout Street",
		HouseNumber: "1",
		City:        "München",
		PostalCode:  "80331",
		CountryCode: "DE",
	}

	cursor := suite.initiateQuery(t, address)
	if cursor == "" {
		t.Fatal("Failed to initiate query")
	}

	// Even with potential timeouts, the API should return some results
	// or handle timeouts gracefully
	products := suite.pollForResults(t, cursor, 45*time.Second)
	t.Logf("Retrieved %d products (some providers may have timed out)", len(products))
}

// TestE2E_SystemEndpoints tests system and health endpoints
func TestE2E_SystemEndpoints(t *testing.T) {
	suite := setupE2ETest(t)
	defer suite.tearDown(t)

	t.Run("HealthCheck", func(t *testing.T) {
		resp, err := suite.client.Get(suite.server.URL + "/health")
		if err != nil {
			t.Fatalf("Health check failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		var health m.Health
		if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
			t.Fatalf("Failed to decode health response: %v", err)
		}

		if health.Status != "ok" {
			t.Errorf("Expected status ok, got %s", health.Status)
		}
	})

	t.Run("Version", func(t *testing.T) {
		resp, err := suite.client.Get(suite.server.URL + "/version")
		if err != nil {
			t.Fatalf("Version check failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		var version m.Version
		if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
			t.Fatalf("Failed to decode version response: %v", err)
		}

		if version.Version != suite.cfg.Version {
			t.Errorf("Expected version %s, got %s", suite.cfg.Version, version.Version)
		}
	})
}
