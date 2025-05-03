package verbyndich

import (
	"fmt"
	"log/slog"
	"testing"

	i "github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
	providertest "github.com/rotmanjanez/check24-gendev-7/pkg/provider/testing"
)

// CreateTestProvider creates a VerbynDich provider instance for testing
func CreateTestProvider(baseURL string, logger *slog.Logger) (i.ProviderAdapter, error) {
	return NewVerbynDichAdapter(
		baseURL,
		"test-api-key",
		1, // blockSize
		logger,
	), nil
}

// Helper functions for pointers
func stringPtr(s string) *string { return &s }
func int32Ptr(i int32) *int32    { return &i }

// Utility functions for creating test cases
func createTestAddress() m.Address {
	return m.Address{Street: "Teststrasse", HouseNumber: "1", City: "Teststadt", PostalCode: "12345", CountryCode: "DE"}
}

func createValidJSONResponse(product, description string, last bool) string {
	return fmt.Sprintf(`{"product":"%s","description":"%s","last":%t,"valid":true}`, product, description, last)
}

func createInvalidJSONResponse(product, description string, last bool) string {
	return fmt.Sprintf(`{"product":"%s","description":"%s","last":%t,"valid":false}`, product, description, last)
}

func createTestCase(address m.Address, urlPath, responseBody string, statusCode int, expectedProducts []m.InternetProduct, expectedError bool) providertest.ProviderTestCase {
	return providertest.ProviderTestCase{
		Address: address,
		URLResponseMap: map[string]providertest.HTTPResponse{
			urlPath: {StatusCode: statusCode, Headers: map[string]string{"Content-Type": "application/json"}, Body: responseBody},
		},
		ExpectedProducts: expectedProducts,
		ExpectedError:    expectedError,
		IsValidResponse:  !expectedError,
	}
}

func createBasicProduct(id, productName string, speed int32, connectionType m.ConnectionType, priceCents int32) m.InternetProduct {
	return m.InternetProduct{
		Id:          id,
		Name:        productName,
		Provider:    "VerbynDich",
		ProductInfo: m.ProductInfo{Speed: speed, ConnectionType: connectionType},
		Pricing:     m.Pricing{MonthlyCostInCent: priceCents},
	}
}

func runDescriptionTest(t *testing.T, testName, description string, expectedProducts []m.InternetProduct, expectedError bool) {
	address := createTestAddress()
	jsonBody := createValidJSONResponse("TestProduct", description, true)
	testCase := createTestCase(address, "/check24/data", jsonBody, 200, expectedProducts, expectedError)
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestVerbynDich_ValidResponse tests a valid response from the provider
func TestVerbynDich_ValidResponse(t *testing.T) {
	description := "Für nur 29€ im Monat erhalten Sie eine DSL-Verbindung mit einer Geschwindigkeit von 100 Mbit/s. Zusätzlich sind folgende Fernsehsender enthalten TestTV."
	expectedProduct := createBasicProduct("TestProduct", "TestProduct", 100, m.DSL, 2900)
	expectedProduct.ProductInfo.Tv = stringPtr("TestTV")
	runDescriptionTest(t, "ValidResponse", description, []m.InternetProduct{expectedProduct}, false)
}

// TestVerbynDich_InvalidJSON tests provider with invalid JSON response
func TestVerbynDich_InvalidJSON(t *testing.T) {
	address := createTestAddress()
	testCase := createTestCase(address, "/check24/data", `not a json`, 200, []m.InternetProduct{}, true)
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestVerbynDich_ErrorStatusCode tests provider with non-200 status code
func TestVerbynDich_ErrorStatusCode(t *testing.T) {
	address := createTestAddress()
	testCase := createTestCase(address, "/check24/data", `{"error":"server error"}`, 500, []m.InternetProduct{}, true)
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestVerbynDich_InvalidValidFlag tests response with valid=false
func TestVerbynDich_InvalidValidFlag(t *testing.T) {
	address := createTestAddress()
	description := "Für nur 29€ im Monat erhalten Sie eine DSL-Verbindung mit einer Geschwindigkeit von 100 Mbit/s."
	jsonBody := createInvalidJSONResponse("TestProduct", description, true)
	testCase := createTestCase(address, "/check24/data", jsonBody, 200, []m.InternetProduct{}, false)
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestVerbynDich_UnknownDescription tests provider with unknown description pattern
func TestVerbynDich_UnknownDescription(t *testing.T) {
	address := createTestAddress()
	description := "This is an unknown description pattern that should not match any regex."
	jsonBody := createValidJSONResponse("TestProduct", description, true)
	testCase := createTestCase(address, "/check24/data", jsonBody, 200, []m.InternetProduct{}, true)
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestVerbynDich_ConflictingInformation tests description with conflicting price information
func TestVerbynDich_ConflictingInformation(t *testing.T) {
	address := createTestAddress()
	// Two conflicting price statements
	description := "Für nur 29€ im Monat erhalten Sie eine DSL-Verbindung mit einer Geschwindigkeit von 100 Mbit/s. Für nur 39€ im Monat erhalten Sie eine FIBER-Verbindung mit einer Geschwindigkeit von 200 Mbit/s."
	jsonBody := createValidJSONResponse("TestProduct", description, true)
	testCase := createTestCase(address, "/check24/data", jsonBody, 200, []m.InternetProduct{}, true)
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestVerbynDich_InvalidConnectionType tests description with unsupported connection type
func TestVerbynDich_InvalidConnectionType(t *testing.T) {
	address := createTestAddress()
	description := "Für nur 29€ im Monat erhalten Sie eine QUANTUM-Verbindung mit einer Geschwindigkeit von 100 Mbit/s."
	jsonBody := createValidJSONResponse("TestProduct", description, true)
	testCase := createTestCase(address, "/check24/data", jsonBody, 200, []m.InternetProduct{}, true)
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestVerbynDich_ConflictingAgeRestrictions tests description with conflicting age restrictions
func TestVerbynDich_ConflictingAgeRestrictions(t *testing.T) {
	address := createTestAddress()
	description := "Für nur 29€ im Monat erhalten Sie eine DSL-Verbindung mit einer Geschwindigkeit von 100 Mbit/s. Dieses Angebot ist nur für Personen unter 65 Jahren verfügbar. Dieses Angebot ist nur für Personen über 18 Jahren verfügbar."
	jsonBody := createValidJSONResponse("TestProduct", description, true)

	expectedProducts := []m.InternetProduct{
		{
			Id:          "TestProduct",
			Name:        "TestProduct",
			Provider:    "VerbynDich",
			ProductInfo: m.ProductInfo{Speed: 100, ConnectionType: m.DSL},
			Pricing:     m.Pricing{MonthlyCostInCent: 2900, MaxAgeInJears: int32Ptr(65), MinAgeInYears: int32Ptr(18)},
		},
	}

	testCase := createTestCase(address, "/check24/data", jsonBody, 200, expectedProducts, false)
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestVerbynDich_ComplexValidDescription tests a complex but valid description
func TestVerbynDich_ComplexValidDescription(t *testing.T) {
	address := createTestAddress()
	description := "Für nur 49€ im Monat erhalten Sie eine FIBER-Verbindung mit einer Geschwindigkeit von 500 Mbit/s. Bitte beachten Sie, dass die Mindestvertragslaufzeit 24 Monate beträgt. Mit diesem Angebot erhalten Sie einen Rabatt von 10% auf Ihre monatliche Rechnung bis zum 12. Monat. Der maximale Rabatt beträgt 30€. Ab dem 25. Monat beträgt der monatliche Preis 59€. Ab 100GB pro Monat wird die Geschwindigkeit gedrosselt. Zusätzlich sind folgende Fernsehsender enthalten Premium+. Dieses Angebot ist nur für Personen unter 30 Jahren verfügbar. Mit diesem Angebot erhalten Sie einen einmaligen Rabatt von 50€ auf Ihre monatliche Rechnung. Der Mindestbestellwert beträgt 25€. Unsere Techniker kümmern sich um die Installation."
	jsonBody := createValidJSONResponse("ComplexProduct", description, true)

	expectedProducts := []m.InternetProduct{
		{
			Id:       "ComplexProduct",
			Name:     "ComplexProduct",
			Provider: "VerbynDich",
			ProductInfo: m.ProductInfo{
				Speed:                 500,
				ConnectionType:        m.FIBER,
				Tv:                    stringPtr("Premium+"),
				UnthrottledCapacityMb: int32Ptr(100000), // 100GB in MB
			},
			Pricing: m.Pricing{
				MonthlyCostInCent:           4900,
				MaxAgeInJears:               int32Ptr(30),
				MinOrderValueInCent:         int32Ptr(2500), // 25€ in cents
				InstallationServiceIncluded: true,
				MinContractDurationInMonths: int32Ptr(24),
				PercentageDiscount:          &m.PercentageDiscount{Percentage: 10, MaxDiscountInCent: int32Ptr(3000)},
				AbsoluteDiscount:            &m.AbsoluteDiscount{ValueInCent: 5000}, // 50€ in cents
				SubsequentCosts:             &m.SubsequentCost{MonthlyCostInCent: 5900, StartMonth: 25},
			},
		},
	}

	testCase := createTestCase(address, "/check24/data", jsonBody, 200, expectedProducts, false)
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestVerbynDich_PartiallyValidDescription tests description with some valid and some invalid parts
func TestVerbynDich_PartiallyValidDescription(t *testing.T) {
	address := createTestAddress()
	// Valid price/speed/type but unknown sentence. As we don't know wether this effects the product, we assume it is invalid.
	description := "Für nur 35€ im Monat erhalten Sie eine CABLE-Verbindung mit einer Geschwindigkeit von 250 Mbit/s. Hier ist ein unbekannter Satz der ignoriert werden sollte. Zusätzlich sind folgende Fernsehsender enthalten Sports."
	jsonBody := createValidJSONResponse("PartialProduct", description, true)

	testCase := createTestCase(address, "/check24/data", jsonBody, 200, []m.InternetProduct{}, true)
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}
