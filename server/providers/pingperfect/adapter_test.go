package pingperfect

import (
	"log/slog"
	"testing"

	i "github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
	providertest "github.com/rotmanjanez/check24-gendev-7/pkg/provider/testing"
)

// CreateTestProvider creates a PingPerfect provider instance for testing
func CreateTestProvider(baseURL string, logger *slog.Logger) (i.ProviderAdapter, error) {
	// Use test credentials for testing
	return NewPingPerfectAdapter(
		baseURL,
		"test-client-id",
		"test-signature-secret",
		logger,
	), nil
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}

// TestPingPerfect_ValidMunichAddress tests a valid Munich address
func TestPingPerfect_ValidMunichAddress(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address: m.Address{
			Street:      "Marienplatz",
			HouseNumber: "1",
			City:        "München",
			PostalCode:  "80331",
			CountryCode: "DE",
		},
		URLResponseMap: map[string]providertest.HTTPResponse{
			"/": {
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body: `[
					{"providerName": "PingPerfect Fiber Pro", "productInfo": {"speed": 1000, "contractDurationInMonths": 24, "connectionType": "FIBER", "tv": "Premium Package", "limitFrom": 100, "maxAge": 65}, "pricingDetails": {"monthlyCostInCent": 4999, "installationService": "yes"}},
					{"providerName": "PingPerfect DSL Basic", "productInfo": {"speed": 50, "contractDurationInMonths": 12, "connectionType": "DSL", "tv": "Basic Package", "limitFrom": 0, "maxAge": 75}, "pricingDetails": {"monthlyCostInCent": 2999, "installationService": "no"}}
				]`,
			},
		},
		ExpectedProducts: []m.InternetProduct{
			{
				Id:       "PingPerfect Fiber Pro",
				Provider: "PingPerfect",
				Name:     "PingPerfect Fiber Pro",
				ProductInfo: m.ProductInfo{
					Speed:                 1000,
					ConnectionType:        m.FIBER,
					Tv:                    stringPtr("Premium Package"),
					UnthrottledCapacityMb: int32Ptr(100000),
				},
				Pricing: m.Pricing{
					MonthlyCostInCent:           4999,
					ContractDurationInMonths:    int32Ptr(24),
					MaxAgeInJears:               int32Ptr(65),
					InstallationServiceIncluded: true,
				},
			},
			{
				Id:       "PingPerfect DSL Basic",
				Provider: "PingPerfect",
				Name:     "PingPerfect DSL Basic",
				ProductInfo: m.ProductInfo{
					Speed:                 50,
					ConnectionType:        m.DSL,
					Tv:                    stringPtr("Basic Package"),
					UnthrottledCapacityMb: nil,
				},
				Pricing: m.Pricing{
					MonthlyCostInCent:           2999,
					ContractDurationInMonths:    int32Ptr(12),
					MaxAgeInJears:               int32Ptr(75),
					InstallationServiceIncluded: false,
				},
			},
		},
		ExpectedError:   false,
		IsValidResponse: true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestPingPerfect_AddressWithoutHouseNumber tests an address without a house number
func TestPingPerfect_AddressWithoutHouseNumber(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address:          m.Address{Street: "Marienplatz", HouseNumber: "", City: "München", PostalCode: "80331", CountryCode: "DE"},
		URLResponseMap:   map[string]providertest.HTTPResponse{},
		ExpectedProducts: []m.InternetProduct{},
		ExpectedError:    false,
		IsValidResponse:  true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestPingPerfect_InvalidJSONResponse tests how provider handles malformed JSON
func TestPingPerfect_InvalidJSONResponse(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address: m.Address{Street: "Marienplatz", HouseNumber: "1", City: "München", PostalCode: "80331", CountryCode: "DE"},
		URLResponseMap: map[string]providertest.HTTPResponse{
			"/": {
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       `{"invalid": "json"`,
			},
		},
		ExpectedProducts: []m.InternetProduct{},
		ExpectedError:    true,
		IsValidResponse:  false,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestPingPerfect_EmptyResponseArray tests provider with empty product array
func TestPingPerfect_EmptyResponseArray(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address: m.Address{Street: "Marienplatz", HouseNumber: "1", City: "München", PostalCode: "80331", CountryCode: "DE"},
		URLResponseMap: map[string]providertest.HTTPResponse{
			"/": {
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       `[]`,
			},
		},
		ExpectedProducts: []m.InternetProduct{},
		ExpectedError:    false,
		IsValidResponse:  true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestPingPerfect_StringInsteadOfNumber tests string in speed field
func TestPingPerfect_StringInsteadOfNumber(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address: m.Address{Street: "Marienplatz", HouseNumber: "1", City: "München", PostalCode: "80331", CountryCode: "DE"},
		URLResponseMap: map[string]providertest.HTTPResponse{
			"/": {
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       `[{"providerName":"StrNum","productInfo":{"speed":"fast","contractDurationInMonths":12,"connectionType":"DSL"},"pricingDetails":{"monthlyCostInCent":1999,"installationService":"no"}}]`,
			},
		},
		ExpectedProducts: []m.InternetProduct{},
		ExpectedError:    true,
		IsValidResponse:  false,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestPingPerfect_NumberInsteadOfString tests number in tv field
func TestPingPerfect_NumberInsteadOfString(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address: m.Address{Street: "Marienplatz", HouseNumber: "1", City: "München", PostalCode: "80331", CountryCode: "DE"},
		URLResponseMap: map[string]providertest.HTTPResponse{
			"/": {
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       `[{"providerName":"NumStr","productInfo":{"speed":10,"contractDurationInMonths":12,"connectionType":"DSL","tv":123},"pricingDetails":{"monthlyCostInCent":1999,"installationService":"no"}}]`,
			},
		},
		ExpectedProducts: []m.InternetProduct{},
		ExpectedError:    true,
		IsValidResponse:  false,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestPingPerfect_AustriaCountryCode tests an address in Austria
func TestPingPerfect_AustriaCountryCode(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address: m.Address{Street: "A", HouseNumber: "1", City: "Wien", PostalCode: "1010", CountryCode: "AT"},
		URLResponseMap: map[string]providertest.HTTPResponse{
			"/": {
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       `[]`,
			},
		},
		ExpectedProducts: []m.InternetProduct{},
		ExpectedError:    false,
		IsValidResponse:  true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestPingPerfect_SwitzerlandCountryCode tests an address in Switzerland
func TestPingPerfect_SwitzerlandCountryCode(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address: m.Address{Street: "B", HouseNumber: "2", City: "Zürich", PostalCode: "8000", CountryCode: "CH"},
		URLResponseMap: map[string]providertest.HTTPResponse{
			"/": {
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       `[]`,
			},
		},
		ExpectedProducts: []m.InternetProduct{},
		ExpectedError:    false,
		IsValidResponse:  true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}
