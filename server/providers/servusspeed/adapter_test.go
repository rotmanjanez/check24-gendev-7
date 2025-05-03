package servusspeed

import (
	"log"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/rotmanjanez/check24-gendev-7/pkg/cache"
	i "github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
	providertest "github.com/rotmanjanez/check24-gendev-7/pkg/provider/testing"
)

// CreateTestProvider creates a ServusSpeed provider instance for testing
func CreateTestProvider(baseURL string, logger *slog.Logger) (i.ProviderAdapter, error) {
	// Use in-memory instance cache for tests
	return NewServusSpeedAdapter(
		"test-username",
		"test-password",
		baseURL,
		cache.NewInstanceCache("test-servusspeed"),
		5*time.Minute,
		logger,
	), nil
}

// Global test address and headers for reuse
var (
	testGermanAddress = m.Address{
		Street:      "Teststrasse",
		HouseNumber: "42",
		City:        "München",
		PostalCode:  "80331",
		CountryCode: "DE",
	}
	applicationJson = map[string]string{"Content-Type": "application/json"}

	// Global response map for all tests
	servusSpeedTestResponses = map[string]providertest.HTTPResponse{
		"/api/external/product-details/prod1": {
			StatusCode: 200,
			Headers:    applicationJson,
			Body:       `{"servusSpeedProduct": {"providerName": "ServusSpeed 1000", "productInfo": {"speed": 1000, "contractDurationInMonths": 24, "connectionType": "FIBER", "tv": "Premium", "limitFrom": 100, "maxAge": 65}, "pricingDetails": {"monthlyCostInCent": 4999, "installationService": true}, "discount": 500}}`,
		},
		"/api/external/product-details/prod2": {
			StatusCode: 200,
			Headers:    applicationJson,
			Body:       `{"servusSpeedProduct": {"providerName": "ServusSpeed DSL", "productInfo": {"speed": 50, "contractDurationInMonths": 12, "connectionType": "DSL", "tv": "Basic", "limitFrom": 0, "maxAge": 75}, "pricingDetails": {"monthlyCostInCent": 2999, "installationService": false}, "discount": 0}}`,
		},
		"/api/external/product-details/strnum": {
			StatusCode: 200,
			Headers:    applicationJson,
			Body:       `{"servusSpeedProduct":{"providerName":"StrNum","productInfo":{"speed":"fast","contractDurationInMonths":12,"connectionType":"DSL"},"pricingDetails":{"monthlyCostInCent":1999,"installationService":false},"discount":0}}`,
		},
		"/api/external/product-details/numstr": {
			StatusCode: 200,
			Headers:    applicationJson,
			Body:       `{"servusSpeedProduct":{"providerName":"NumStr","productInfo":{"speed":10,"contractDurationInMonths":12,"connectionType":"DSL","tv":123},"pricingDetails":{"monthlyCostInCent":1999,"installationService":false},"discount":0}}`,
		},
	}

	prod1 = m.InternetProduct{
		Id:       "prod1",
		Provider: "ServusSpeed",
		Name:     "ServusSpeed 1000",
		ProductInfo: m.ProductInfo{
			Speed:                 1000,
			ConnectionType:        m.FIBER,
			Tv:                    stringPtr("Premium"),
			UnthrottledCapacityMb: int32Ptr(100000),
		},
		Pricing: m.Pricing{
			MonthlyCostInCent:           4999,
			ContractDurationInMonths:    int32Ptr(24),
			MaxAgeInJears:               int32Ptr(65),
			InstallationServiceIncluded: true,
			AbsoluteDiscount:            &m.AbsoluteDiscount{ValueInCent: 500},
		},
	}

	prod2 = m.InternetProduct{
		Id:       "prod2",
		Provider: "ServusSpeed",
		Name:     "ServusSpeed DSL",
		ProductInfo: m.ProductInfo{
			Speed:                 50,
			ConnectionType:        m.DSL,
			Tv:                    stringPtr("Basic"),
			UnthrottledCapacityMb: nil,
		},
		Pricing: m.Pricing{
			MonthlyCostInCent:           2999,
			ContractDurationInMonths:    int32Ptr(12),
			MaxAgeInJears:               int32Ptr(75),
			InstallationServiceIncluded: false,
			AbsoluteDiscount:            nil,
		},
	}
)

func createUrlResponseMap(availableProducts []string) map[string]providertest.HTTPResponse {
	responseMap := make(map[string]providertest.HTTPResponse)

	productListString := `["` + strings.Join(availableProducts, `", "`) + `"]`

	responseMap["/api/external/available-products"] = providertest.HTTPResponse{
		StatusCode: 200,
		Headers:    applicationJson,
		Body:       `{"availableProducts": ` + productListString + `}`,
	}

	for _, product := range availableProducts {
		response, exists := servusSpeedTestResponses["/api/external/product-details/"+product]
		if !exists {
			log.Fatalf("No response defined for product %s", product)
		}
		responseMap["/api/external/product-details/"+product] = response
	}

	return responseMap
}

// Helper functions
func stringPtr(s string) *string { return &s }
func int32Ptr(i int32) *int32    { return &i }

func TestServusSpeed_ValidGermanAddress(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address:          testGermanAddress,
		URLResponseMap:   createUrlResponseMap([]string{"prod1", "prod2"}),
		ExpectedProducts: []m.InternetProduct{prod1, prod2},
		ExpectedError:    false,
		IsValidResponse:  true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

func TestServusSpeed_AddressWithoutHouseNumber(t *testing.T) {
	address := testGermanAddress
	address.HouseNumber = ""
	// Should behave the same as a valid address (empty house number allowed)
	testCase := providertest.ProviderTestCase{
		Address:          address,
		URLResponseMap:   createUrlResponseMap([]string{"prod1", "prod2"}),
		ExpectedProducts: []m.InternetProduct{}, // ServusSpeed does not return products for addresses without house number
		ExpectedError:    false,
		IsValidResponse:  true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

func TestServusSpeed_InvalidJSONResponse(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address: testGermanAddress,
		URLResponseMap: map[string]providertest.HTTPResponse{
			"/api/external/available-products": {
				StatusCode: 200,
				Headers:    applicationJson,
				Body:       `{"invalid": "json"`,
			},
		},
		ExpectedProducts: []m.InternetProduct{},
		ExpectedError:    true,
		IsValidResponse:  false,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

func TestServusSpeed_EmptyResponseArray(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address:          testGermanAddress,
		URLResponseMap:   createUrlResponseMap([]string{}),
		ExpectedProducts: []m.InternetProduct{},
		ExpectedError:    false,
		IsValidResponse:  true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

func TestServusSpeed_StringInsteadOfNumber(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address:          testGermanAddress,
		URLResponseMap:   createUrlResponseMap([]string{"strnum"}),
		ExpectedProducts: []m.InternetProduct{},
		ExpectedError:    true,
		IsValidResponse:  false,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

func TestServusSpeed_NumberInsteadOfString(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address:          testGermanAddress,
		URLResponseMap:   createUrlResponseMap([]string{"numstr"}),
		ExpectedProducts: []m.InternetProduct{},
		ExpectedError:    true,
		IsValidResponse:  false,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

func TestServusSpeed_AustriaCountryCode(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address:          m.Address{Street: "A", HouseNumber: "1", City: "Wien", PostalCode: "1010", CountryCode: "AT"},
		URLResponseMap:   nil, // No products available for AT
		ExpectedProducts: []m.InternetProduct{},
		ExpectedError:    false,
		IsValidResponse:  true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

func TestServusSpeed_SwitzerlandCountryCode(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address:          m.Address{Street: "B", HouseNumber: "2", City: "Zürich", PostalCode: "8000", CountryCode: "CH"},
		URLResponseMap:   nil, // No products available for CH
		ExpectedProducts: []m.InternetProduct{},
		ExpectedError:    false,
		IsValidResponse:  true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}
