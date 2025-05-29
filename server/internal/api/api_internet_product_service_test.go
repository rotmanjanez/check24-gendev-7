package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/rotmanjanez/check24-gendev-7/pkg/cache"
	"github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
	"github.com/rotmanjanez/check24-gendev-7/pkg/models"
	"github.com/rotmanjanez/check24-gendev-7/pkg/provider"
)

// Global test data for reuse across tests
var (
	validAddressDE = models.Address{
		Street:      "Teststrasse",
		HouseNumber: "1",
		City:        "Berlin",
		PostalCode:  "10115",
		CountryCode: "DE",
	}

	validAddressCH = models.Address{
		Street:      "Bahnhofstrasse",
		HouseNumber: "1",
		City:        "Zürich",
		PostalCode:  "8001",
		CountryCode: "CH",
	}

	validAddressAT = models.Address{
		Street:      "Stephansplatz",
		HouseNumber: "1",
		City:        "Wien",
		PostalCode:  "1010",
		CountryCode: "AT",
	}

	// Address with excessively long street name (over reasonable limits)
	invalidAddressLongStreet = models.Address{
		Street:      strings.Repeat("Very-Long-Street-Name-", 50), // 1100+ chars
		HouseNumber: "1",
		City:        "Berlin",
		PostalCode:  "10115",
		CountryCode: "DE",
	}

	sampleProduct = models.InternetProduct{
		Id:          "test-product-1",
		Provider:    "mock-provider",
		Name:        "Test Internet Package",
		DateOffered: time.Now(),
		ProductInfo: models.ProductInfo{
			Speed:          100,
			ConnectionType: models.FIBER,
		},
		Pricing: models.Pricing{
			MonthlyCostInCent:        4999,
			ContractDurationInMonths: &[]int32{24}[0],
		},
	}
)

// mockProviderAdapter is a simple provider that records the address it receives
// and returns a canned product for testing
type mockProviderAdapter struct {
	lastRequest             interfaces.Request
	lastMetadata            interface{}
	returnProductsOnPrepare bool
	returnProductsOnParse   bool
	mockServer              *httptest.Server
	product                 *models.InternetProduct
}

func (m *mockProviderAdapter) PrepareRequest(ctx context.Context, req interfaces.Request) (interfaces.ParsedResponse, error) {
	slog.Error("Preparing request for mock provider", "address", req.Address)
	m.lastRequest = req

	response := interfaces.ParsedResponse{}

	if m.returnProductsOnPrepare {
		if m.product != nil {
			response.InternetProducts = []models.InternetProduct{*m.product}
		} else {
			response.InternetProducts = []models.InternetProduct{sampleProduct}
		}
	}

	if m.mockServer != nil {
		httpReq, _ := http.NewRequest("GET", m.mockServer.URL+"/test", nil)
		response.Requests = []interfaces.PreparedRequest{
			{
				Request:  httpReq,
				Metadata: "test-metadata",
			},
		}
	}

	return response, nil
}

func (m *mockProviderAdapter) ParseResponse(ctx context.Context, resp interfaces.Response) (interfaces.ParsedResponse, error) {
	response := interfaces.ParsedResponse{}

	if m.returnProductsOnParse {
		response.InternetProducts = []models.InternetProduct{sampleProduct}
	}

	return response, nil
}

func (m *mockProviderAdapter) Name() string { return "mock" }

// MockProviderAdapter implements the ProviderAdapter interface
var _ interfaces.ProviderAdapter = &mockProviderAdapter{}

// setupTestService creates a service with cache and mock provider for testing
func setupTestService(mockProvider *mockProviderAdapter) (*mux.Router, *InternetProductsAPIController) {
	cacheInst := cache.NewInstanceCache("test-cache")
	queueInst := cache.NewInstanceCache("test-queue")

	service := NewInternetProductsAPIService(
		nil,
		cacheInst,
		queueInst,
		[]*provider.ProviderConfig{provider.NewProviderConfig(mockProvider, 0, 10*time.Minute /* when debugging is required*/, 1, 0)},
	)

	controller := NewInternetProductsAPIController(service)
	return NewRouter(controller), controller
}

// createRequestFromAddress creates an HTTP request with the given address as JSON body
func createRequestFromAddress(address models.Address) *http.Request {
	body, _ := json.Marshal(address)
	req := httptest.NewRequest(http.MethodPost, "/internet-products", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestInitiateInternetProductsQuery_ValidGermanAddress(t *testing.T) {
	mockProvider := &mockProviderAdapter{}
	_, controller := setupTestService(mockProvider)

	req := createRequestFromAddress(validAddressDE)
	w := httptest.NewRecorder()

	controller.InitiateInternetProductsQuery(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	var result models.InternetProductsCursor
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("invalid response: %v", err)
	}
	if result.NextCursor == "" {
		t.Errorf("expected a cursor in response")
	}
}

func TestInitiateInternetProductsQuery_ValidSwissAddress(t *testing.T) {
	mockProvider := &mockProviderAdapter{}
	_, controller := setupTestService(mockProvider)

	req := createRequestFromAddress(validAddressCH)
	w := httptest.NewRecorder()

	controller.InitiateInternetProductsQuery(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	var result models.InternetProductsCursor
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("invalid response: %v", err)
	}
	if result.NextCursor == "" {
		t.Errorf("expected a cursor in response")
	}
}

func TestInitiateInternetProductsQuery_ValidAustrianAddress(t *testing.T) {
	mockProvider := &mockProviderAdapter{}
	_, controller := setupTestService(mockProvider)

	req := createRequestFromAddress(validAddressAT)
	w := httptest.NewRecorder()

	controller.InitiateInternetProductsQuery(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	var result models.InternetProductsCursor
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("invalid response: %v", err)
	}
	if result.NextCursor == "" {
		t.Errorf("expected a cursor in response")
	}
}

func TestInitiateInternetProductsQuery_InvalidLongStreetAddress(t *testing.T) {
	mockProvider := &mockProviderAdapter{}
	_, controller := setupTestService(mockProvider)

	req := createRequestFromAddress(invalidAddressLongStreet)
	w := httptest.NewRecorder()

	controller.InitiateInternetProductsQuery(w, req)
	resp := w.Result()

	// The service should still process the request, but we can check that the address was piped through
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK even for long street, got %d", resp.StatusCode)
	}

	var result models.InternetProductsCursor
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("invalid response: %v", err)
	}
	if result.NextCursor == "" {
		t.Errorf("expected a cursor in response")
	}
}

func TestInitiateInternetProductsQuery_ProductsFromPrepareRequest(t *testing.T) {
	mockProvider := &mockProviderAdapter{
		returnProductsOnPrepare: true,
	}
	_, controller := setupTestService(mockProvider)

	req := createRequestFromAddress(validAddressDE)
	w := httptest.NewRecorder()

	controller.InitiateInternetProductsQuery(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	var result models.InternetProductsCursor
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("invalid response: %v", err)
	}
	if result.NextCursor == "" {
		t.Errorf("expected a cursor in response")
	}

	time.Sleep(10 * time.Millisecond) // Give goroutine time to run
	// Verify the address was piped through correctly
	if mockProvider.lastRequest.Address.Street != validAddressDE.Street {
		t.Errorf("expected address street %s, got %s", validAddressDE.Street, mockProvider.lastRequest.Address.Street)
	}
}

func TestInitiateInternetProductsQuery_ProductsFromParseResponse(t *testing.T) {
	// Create mock server for HTTP requests
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"products": []}`))
	}))
	defer mockServer.Close()

	mockProvider := &mockProviderAdapter{
		mockServer: mockServer,
	}
	_, controller := setupTestService(mockProvider)

	req := createRequestFromAddress(validAddressDE)
	w := httptest.NewRecorder()

	controller.InitiateInternetProductsQuery(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	var result models.InternetProductsCursor
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("invalid response: %v", err)
	}
	if result.NextCursor == "" {
		t.Errorf("expected a cursor in response")
	}

	time.Sleep(10 * time.Millisecond) // Give goroutine time to run
	// Verify the address was piped through correctly
	if mockProvider.lastRequest.Address.Street != validAddressDE.Street {
		t.Errorf("expected address street %s, got %s", validAddressDE.Street, mockProvider.lastRequest.Address.Street)
	}
}
