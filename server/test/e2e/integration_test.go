package e2e

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/rotmanjanez/check24-gendev-7/internal/requestmanager"
	"github.com/rotmanjanez/check24-gendev-7/pkg/cache"
	i "github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
	p "github.com/rotmanjanez/check24-gendev-7/pkg/provider"
	"github.com/rotmanjanez/check24-gendev-7/providers/byteme"
	"github.com/rotmanjanez/check24-gendev-7/providers/pingperfect"
	"github.com/rotmanjanez/check24-gendev-7/providers/servusspeed"
	"github.com/rotmanjanez/check24-gendev-7/providers/verbyndich"
)

func TestE2E_ProviderIntegration(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Create mock servers for all providers
	byteMeServer := createByteMeTestServer()
	verbynDichServer := createVerbynDichTestServer()
	servusSpeedServer := createServusSpeedTestServer()
	pingPerfectServer := createPingPerfectTestServer()

	defer func() {
		byteMeServer.Close()
		verbynDichServer.Close()
		servusSpeedServer.Close()
		pingPerfectServer.Close()
	}()

	// Create test cache
	testCache := cache.NewInstanceCache("integration-test")

	// Create provider configs using the actual NewProviderConfig constructor
	providers := []*p.ProviderConfig{
		p.NewProviderConfig(
			byteme.NewByteMeAdapter(byteMeServer.URL, "test-api-key", logger),
			3,                    // retries
			5*time.Second,        // timeout
			1,                    // maxConcurrent
			500*time.Millisecond, // backoff
		),
		p.NewProviderConfig(
			verbyndich.NewVerbynDichAdapter(verbynDichServer.URL+"/check24/data", "test-api-key", 1, logger),
			3,                    // retries
			5*time.Second,        // timeout
			1,                    // maxConcurrent
			500*time.Millisecond, // backoff
		),
		p.NewProviderConfig(
			servusspeed.NewServusSpeedAdapter("test-username", "test-password", servusSpeedServer.URL, testCache, 5*time.Minute, logger),
			3,                    // retries
			5*time.Second,        // timeout
			1,                    // maxConcurrent
			500*time.Millisecond, // backoff
		),
		p.NewProviderConfig(
			pingperfect.NewPingPerfectAdapter(pingPerfectServer.URL, "test-client-id", "test-signature-secret", logger),
			3,                    // retries
			5*time.Second,        // timeout
			1,                    // maxConcurrent
			500*time.Millisecond, // backoff
		),
	}

	// Create request coordinator
	coordinator := requestmanager.NewRequestCoordinator(providers)

	// Test address
	address := m.Address{
		Street:      "Marienplatz",
		HouseNumber: "1",
		City:        "München",
		PostalCode:  "80331",
		CountryCode: "DE",
	}

	// Test the coordination
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := i.Request{Address: address}
	// Execute coordination
	responsesChan, errorsChan := coordinator.Run(ctx, request, 100, 100)

	// Collect results
	var products []m.InternetProduct
	for product := range responsesChan {
		products = append(products, product)
	}

	var errors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}

	// Verify we got products from multiple providers
	if len(products) == 0 {
		t.Error("Expected to receive products from providers")
	}

	// Print results for debugging
	logger.Info("Integration test results", "products", len(products), "errors", len(errors))
	for _, product := range products {
		logger.Info("Product", "provider", product.Provider, "name", product.Name, "price", product.Pricing.MonthlyCostInCent)
	}
	for _, err := range errors {
		logger.Warn("Error", "error", err)
	}
}

func TestE2E_ProviderErrorHandling(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Create mock servers - one working, one failing
	workingServer := createByteMeTestServer()
	failingServer := createFailingTestServer()

	defer func() {
		workingServer.Close()
		failingServer.Close()
	}()

	// Create provider configs with one working and one failing provider
	providers := []*p.ProviderConfig{
		p.NewProviderConfig(
			byteme.NewByteMeAdapter(workingServer.URL, "test-api-key", logger),
			3,                    // retries
			5*time.Second,        // timeout
			1,                    // maxConcurrent
			500*time.Millisecond, // backoff
		),
		p.NewProviderConfig(
			byteme.NewByteMeAdapter(failingServer.URL, "test-api-key", logger),
			3,                    // retries
			5*time.Second,        // timeout
			1,                    // maxConcurrent
			500*time.Millisecond, // backoff
		),
	}

	// Create request coordinator
	coordinator := requestmanager.NewRequestCoordinator(providers)

	// Test address
	address := m.Address{
		Street:      "Marienplatz",
		HouseNumber: "1",
		City:        "München",
		PostalCode:  "80331",
		CountryCode: "DE",
	}

	// Test the coordination
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request := i.Request{Address: address}
	// Execute coordination
	responsesChan, errorsChan := coordinator.Run(ctx, request, 100, 100)

	// Collect results
	var products []m.InternetProduct
	for product := range responsesChan {
		products = append(products, product)
	}

	var errors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}

	// Verify we got products from working provider and errors from failing provider
	if len(products) == 0 {
		t.Error("Expected to receive products from working provider")
	}

	if len(errors) == 0 {
		t.Error("Expected to receive errors from failing provider")
	}

	// Print results for debugging
	logger.Info("Error handling test results", "products", len(products), "errors", len(errors))
}

func TestE2E_ProviderTimeouts(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Create slow mock server
	slowServer := createSlowTestServer()
	defer slowServer.Close()

	// Create provider config with very short timeout
	providers := []*p.ProviderConfig{
		p.NewProviderConfig(
			byteme.NewByteMeAdapter(slowServer.URL, "test-api-key", logger),
			1,                    // retries
			1*time.Second,        // short timeout
			1,                    // maxConcurrent
			100*time.Millisecond, // backoff
		),
	}

	// Create request coordinator
	coordinator := requestmanager.NewRequestCoordinator(providers)

	// Test address
	address := m.Address{
		Street:      "Marienplatz",
		HouseNumber: "1",
		City:        "München",
		PostalCode:  "80331",
		CountryCode: "DE",
	}

	// Test the coordination with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	request := i.Request{Address: address}
	// Execute coordination
	responsesChan, errorsChan := coordinator.Run(ctx, request, 100, 100)

	// Collect results
	var products []m.InternetProduct
	for product := range responsesChan {
		products = append(products, product)
	}

	var errors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}

	// Verify timeout handling - should get timeout errors
	if len(errors) == 0 {
		t.Error("Expected timeout errors from slow provider")
	}

	// Print results for debugging
	logger.Info("Timeout test results", "products", len(products), "errors", len(errors))
}

// createFailingTestServer creates a server that always returns 500 errors
func createFailingTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
}

// createSlowTestServer creates a server that responds very slowly to test timeouts
func createSlowTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second) // Sleep for 5 seconds to trigger timeout
		w.Header().Set("Content-Type", "text/csv")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("productId,providerName,speed,monthlyCostInCent\n1,SlowProvider,100,2999\n"))
	}))
}
