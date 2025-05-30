package testing

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
)

// ProviderTestCase represents a single test case for a provider
type ProviderTestCase struct {
	Address m.Address

	// URLResponseMap maps request URLs to their expected response bodies
	URLResponseMap map[string]HTTPResponse

	// Expected results
	ExpectedProducts []m.InternetProduct
	ExpectedError    bool

	// Whether this represents a valid or invalid response scenario
	IsValidResponse bool
}

// HTTPResponse represents an HTTP response for testing
type HTTPResponse struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

// ProviderTestData contains all test cases for a provider
type ProviderTestData struct {
	ProviderName string
	TestCases    []ProviderTestCase
}

// MockHTTPServer creates a test server that responds based on URLResponseMap
func (tc *ProviderTestCase) MockHTTPServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()

		// Try to find a matching response
		for urlPattern, response := range tc.URLResponseMap {
			if strings.Contains(url, urlPattern) || url == urlPattern {
				// Set headers
				for key, value := range response.Headers {
					w.Header().Set(key, value)
				}

				// Set status code
				w.WriteHeader(response.StatusCode)

				// Write body
				if response.Body != "" {
					w.Write([]byte(response.Body))
				}
				return
			}
		}

		// No matching response found - return 404
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Mock: No response configured for URL: " + url))
	}))
}

// RunProviderTestCase runs a single provider test case (for use in per-test-case test functions)
func RunProviderTestCase(t *testing.T, tc ProviderTestCase, createProvider func(baseURL string, logger *slog.Logger) (interfaces.ProviderAdapter, error)) {
	// Create mock server
	server := tc.MockHTTPServer()
	defer server.Close()

	// Create provider with mock server URL
	logger := slog.New(slog.NewTextHandler(log.Writer(), &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	provider, err := createProvider(server.URL, logger)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Test PrepareRequest
	ctx := context.Background()
	request := interfaces.Request{
		Address: tc.Address,
	}

	hasError := false

	parsedResponse, err := provider.PrepareRequest(ctx, request)
	if err != nil {
		if !tc.ExpectedError {
			t.Errorf("Unexpected error in PrepareRequest: %v", err)
			return
		} else {
			hasError = true
		}
	}

	// Test each prepared request by simulating HTTP calls
	allProducts := parsedResponse.InternetProducts

	pendingRequests := parsedResponse.Requests

	for len(pendingRequests) > 0 {
		// Take the first pending request
		preparedReq := pendingRequests[0]
		pendingRequests = pendingRequests[1:]

		if preparedReq.Request == nil {
			t.Errorf("Prepared request is nil")
			continue
		}

		// Create a mock client that uses our test server
		client := &http.Client{}

		// Execute the request
		resp, err := client.Do(preparedReq.Request)
		if err != nil {
			t.Errorf("Failed to execute request: %v", err)
			continue
		}

		// Create response for ParseResponse
		response := interfaces.Response{
			InitialRequestData: request,
			Request:            preparedReq,
			HTTPResponse:       resp,
		}

		// Test ParseResponse
		parsedResp, err := provider.ParseResponse(ctx, response)
		if tc.ExpectedError && err != nil {
			hasError = true
		}
		if !tc.ExpectedError && err != nil {
			t.Errorf("Unexpected error in ParseResponse: %v", err)
		}

		allProducts = append(allProducts, parsedResp.InternetProducts...)
		pendingRequests = append(pendingRequests, parsedResp.Requests...)

		// Clean up
		resp.Body.Close()
	}

	if tc.ExpectedError && !hasError {
		t.Errorf("Expected error in ParseResponse but got none")
	}

	// Validate results
	if tc.IsValidResponse {
		validateProducts(t, tc.ExpectedProducts, allProducts)
	}
}

// validateProducts compares expected vs actual products
func validateProducts(t *testing.T, expected, actual []m.InternetProduct) {
	if len(expected) != len(actual) {
		t.Errorf("Expected %d products, got %d", len(expected), len(actual))
		return
	}

	for i, expectedProduct := range expected {
		if i >= len(actual) {
			t.Errorf("Missing product at index %d", i)
			continue
		}

		actualProduct := actual[i]

		// Basic validation - can be extended
		if expectedProduct.Provider != actualProduct.Provider {
			t.Errorf("Product %d: Expected provider %s, got %s", i, expectedProduct.Provider, actualProduct.Provider)
		}

		if expectedProduct.Name != actualProduct.Name {
			t.Errorf("Product %d: Expected name %s, got %s", i, expectedProduct.Name, actualProduct.Name)
		}

		if expectedProduct.ProductInfo.Speed != actualProduct.ProductInfo.Speed {
			t.Errorf("Product %d: Expected speed %d, got %d", i, expectedProduct.ProductInfo.Speed, actualProduct.ProductInfo.Speed)
		}
	}
}
