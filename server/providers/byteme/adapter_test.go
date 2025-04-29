package byteme

import (
	"log/slog"
	"testing"

	i "github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
	providertest "github.com/rotmanjanez/check24-gendev-7/pkg/provider/testing"
)

// CreateTestProvider creates a ByteMe provider instance for testing
func CreateTestProvider(baseURL string, logger *slog.Logger) (i.ProviderAdapter, error) {
	return NewByteMeAdapter(
		baseURL,
		"test-api-key",
		logger,
	), nil
}

// Helper functions for pointers
func stringPtr(s string) *string { return &s }
func int32Ptr(i int32) *int32    { return &i }

// TestByteMe_ValidMunichAddress tests a valid Munich address with CSV response
func TestByteMe_ValidMunichAddress(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address: m.Address{Street: "Marienplatz", HouseNumber: "1", City: "München", PostalCode: "80331", CountryCode: "DE"},
		URLResponseMap: map[string]providertest.HTTPResponse{
			"/": {StatusCode: 200, Headers: map[string]string{"Content-Type": "text/csv"}, Body: `productId,providerName,speed,monthlyCostInCent,afterTwoYearsMonthlyCost,durationInMonths,connectionType,installationService,tv,limitFrom,maxAge,voucherType,voucherValue
1,ByteMe,100,2999,1999,24,DSL,true,Premium,0,65,,0
2,ByteMe,50,1999,1499,12,FIBER,false,Basic,50,75,,0
`},
		},
		ExpectedProducts: []m.InternetProduct{
			{
				Id:          "1",
				Provider:    "ByteMe",
				Name:        "ByteMe",
				ProductInfo: m.ProductInfo{Speed: 100, ConnectionType: m.DSL, Tv: stringPtr("Premium"), UnthrottledCapacityMb: nil},
				Pricing:     m.Pricing{MonthlyCostInCent: 2999, ContractDurationInMonths: int32Ptr(24), MaxAgeInJears: int32Ptr(65), InstallationServiceIncluded: true, SubsequentCosts: &m.SubsequentCost{MonthlyCostInCent: 1999, StartMonth: 25}},
			},
			{
				Id:          "2",
				Provider:    "ByteMe",
				Name:        "ByteMe",
				ProductInfo: m.ProductInfo{Speed: 50, ConnectionType: m.FIBER, Tv: stringPtr("Basic"), UnthrottledCapacityMb: int32Ptr(50000)},
				Pricing:     m.Pricing{MonthlyCostInCent: 1999, ContractDurationInMonths: int32Ptr(12), MaxAgeInJears: int32Ptr(75), InstallationServiceIncluded: false, SubsequentCosts: &m.SubsequentCost{MonthlyCostInCent: 1499, StartMonth: 25}},
			},
		},
		ExpectedError:   false,
		IsValidResponse: true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestByteMe_EmptyResponseCSV tests provider with empty CSV header only
func TestByteMe_EmptyResponseCSV(t *testing.T) {
	head := "productId,providerName,speed,monthlyCostInCent,afterTwoYearsMonthlyCost,durationInMonths,connectionType,installationService,tv,limitFrom,maxAge,voucherType,voucherValue\n"
	testCase := providertest.ProviderTestCase{
		Address:          m.Address{Street: "Marienplatz", HouseNumber: "1", City: "München", PostalCode: "80331", CountryCode: "DE"},
		URLResponseMap:   map[string]providertest.HTTPResponse{"/": {StatusCode: 200, Headers: map[string]string{"Content-Type": "text/csv"}, Body: head}},
		ExpectedProducts: []m.InternetProduct{},
		ExpectedError:    false,
		IsValidResponse:  true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestByteMe_InvalidCSVResponse tests provider with malformed CSV
func TestByteMe_InvalidCSVResponse(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address:          m.Address{Street: "Marienplatz", HouseNumber: "1", City: "München", PostalCode: "80331", CountryCode: "DE"},
		URLResponseMap:   map[string]providertest.HTTPResponse{"/": {StatusCode: 200, Headers: map[string]string{"Content-Type": "text/csv"}, Body: `not,a,csv`}},
		ExpectedProducts: []m.InternetProduct{},
		ExpectedError:    true,
		IsValidResponse:  false,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestByteMe_UnknownVoucherType tests unrecognized voucher type handling
func TestByteMe_UnknownVoucherType(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address: m.Address{Street: "Marienplatz", HouseNumber: "1", City: "München", PostalCode: "80331", CountryCode: "DE"},
		URLResponseMap: map[string]providertest.HTTPResponse{"/": {StatusCode: 200, Headers: map[string]string{"Content-Type": "text/csv"}, Body: `productId,providerName,speed,monthlyCostInCent,afterTwoYearsMonthlyCostInCent,durationInMonths,connectionType,installationService,tv,limitFrom,maxAge,voucherType,voucherValue
1,ByteMe,100,2999,1999,24,DSL,true,,0,65,invalid,0
`}},
		ExpectedProducts: []m.InternetProduct{},
		ExpectedError:    true,
		IsValidResponse:  false,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestByteMe_AbsoluteAndPercentageVoucher tests parsing of absolute and percentage discounts
func TestByteMe_AbsoluteAndPercentageVoucher(t *testing.T) {
	head := "productId,providerName,speed,monthlyCostInCent,afterTwoYearsMonthlyCost,durationInMonths,connectionType,installationService,tv,limitFrom,maxAge,voucherType,voucherValue\n"
	csv := head +
		"1,ByteMe 1,100,2999,1999,24,DSL,true,,0,65,absolute,100\n" +
		"2,ByteMe 2,50,1999,1499,12,FIBER,false,,0,75,percentage,10\n"
	testCase := providertest.ProviderTestCase{
		Address:        m.Address{Street: "Marienplatz", HouseNumber: "1", City: "München", PostalCode: "80331", CountryCode: "DE"},
		URLResponseMap: map[string]providertest.HTTPResponse{"/": {StatusCode: 200, Headers: map[string]string{"Content-Type": "text/csv"}, Body: csv}},
		ExpectedProducts: []m.InternetProduct{
			{
				Id:          "1",
				Provider:    "ByteMe",
				Name:        "ByteMe 1",
				ProductInfo: m.ProductInfo{Speed: 100, ConnectionType: m.DSL, Tv: stringPtr(""), UnthrottledCapacityMb: nil},
				Pricing:     m.Pricing{MonthlyCostInCent: 2999, ContractDurationInMonths: int32Ptr(24), MaxAgeInJears: int32Ptr(65), InstallationServiceIncluded: true, SubsequentCosts: &m.SubsequentCost{MonthlyCostInCent: 1999, StartMonth: 25}, AbsoluteDiscount: &m.AbsoluteDiscount{ValueInCent: 100}},
			},
			{
				Id:          "2",
				Provider:    "ByteMe",
				Name:        "ByteMe 2",
				ProductInfo: m.ProductInfo{Speed: 50, ConnectionType: m.FIBER, Tv: stringPtr(""), UnthrottledCapacityMb: nil},
				Pricing:     m.Pricing{MonthlyCostInCent: 1999, ContractDurationInMonths: int32Ptr(12), MaxAgeInJears: int32Ptr(75), InstallationServiceIncluded: false, SubsequentCosts: &m.SubsequentCost{MonthlyCostInCent: 1499, StartMonth: 25}, PercentageDiscount: &m.PercentageDiscount{Percentage: 10}},
			},
		},
		ExpectedError:   false,
		IsValidResponse: true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}
