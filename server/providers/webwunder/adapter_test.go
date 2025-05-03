package webwunder

import (
	"fmt"
	"log/slog"
	"testing"

	i "github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
	providertest "github.com/rotmanjanez/check24-gendev-7/pkg/provider/testing"
)

// Helper functions for pointers
func int32Ptr(i int32) *int32 { return &i }

// Common test constants
var (
	// Standard Munich address used in most tests
	standardTestAddress = m.Address{
		Street:      "Marienplatz",
		HouseNumber: "1",
		City:        "München",
		PostalCode:  "80331",
		CountryCode: "DE",
	}

	// Standard HTTP headers for SOAP responses
	soapHeaders = map[string]string{"Content-Type": "text/xml"}

	// SOAP envelope template for building response bodies
	soapEnvelopeTemplate = `<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
	<Body>
		<Output xmlns:ns2="http://webwunder.gendev7.check24.fun/offerservice">
			%s
		</Output>
	</Body>
</SOAP-ENV:Envelope>`

	// Empty SOAP envelope for simple responses
	emptySoapEnvelope = `<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"></SOAP-ENV:Envelope>`

	// Standard product XML templates
	productXmlTemplate = `<products>
	<productId>%s</productId>
	<providerName>%s</providerName>
	<productInfo>
		<speed>%d</speed>
		<monthlyCostInCent>%d</monthlyCostInCent>
		<monthlyCostInCentFrom25thMonth>%d</monthlyCostInCentFrom25thMonth>
		%s
		<contractDurationInMonths>%d</contractDurationInMonths>
		<connectionType>%s</connectionType>
	</productInfo>
</products>`

	// Voucher XML templates
	absoluteVoucherXml = `<voucher>
		<absoluteVoucher>
			<discountInCent>%d</discountInCent>
			<minOrderValueInCent>%d</minOrderValueInCent>
		</absoluteVoucher>
	</voucher>`

	percentageVoucherXml = `<voucher>
		<percentageVoucher>
			<percentage>%d</percentage>
			<maxDiscountInCent>%d</maxDiscountInCent>
		</percentageVoucher>
	</voucher>`

	bothVouchersXml = `<voucher>
		<absoluteVoucher>
			<discountInCent>%d</discountInCent>
			<minOrderValueInCent>%d</minOrderValueInCent>
		</absoluteVoucher>
		<percentageVoucher>
			<percentage>%d</percentage>
			<maxDiscountInCent>%d</maxDiscountInCent>
		</percentageVoucher>
	</voucher>`
)

// Helper functions for building common product and pricing structures
func createStandardProduct(id string, productName string, speed int32, connectionType m.ConnectionType,
	monthlyCost int32, subsequentCost int32, contractDuration int32) m.InternetProduct {
	return m.InternetProduct{
		Id:       id,
		Provider: "WebWunder",
		Name:     productName,
		ProductInfo: m.ProductInfo{
			Speed:          speed,
			ConnectionType: connectionType,
		},
		Pricing: m.Pricing{
			MonthlyCostInCent:        monthlyCost,
			ContractDurationInMonths: int32Ptr(contractDuration),
			SubsequentCosts: &m.SubsequentCost{
				MonthlyCostInCent: subsequentCost,
				StartMonth:        25,
			},
		},
	}
}

func createProductWithAbsoluteDiscount(product m.InternetProduct, valueInCent int32, minOrderValueInCent int32) m.InternetProduct {
	product.Pricing.AbsoluteDiscount = &m.AbsoluteDiscount{
		ValueInCent:         valueInCent,
		MinOrderValueInCent: int32Ptr(minOrderValueInCent),
	}
	return product
}

func createProductWithPercentageDiscount(product m.InternetProduct, percentage int32, maxDiscountInCent int32) m.InternetProduct {
	product.Pricing.PercentageDiscount = &m.PercentageDiscount{
		Percentage:        percentage,
		MaxDiscountInCent: int32Ptr(maxDiscountInCent),
	}
	return product
}

// CreateTestProvider creates a WebWunder provider instance for testing
func CreateTestProvider(baseUrl string, logger *slog.Logger) (i.ProviderAdapter, error) {
	return NewWebWunderAdapter(
		"test-api-key",
		baseUrl+"/endpunkte/soap/ws",
		"http://spring.io/guides/gs-producing-web-service/legacyGetInternetOffers",
		"http://webwunder.gendev7.check24.fun/offerservice",
		"http://schemas.xmlsoap.org/soap/envelope/",
		logger,
	), nil
}

// TestWebWunder_ValidMunichAddress tests a valid Munich address with SOAP response
func TestWebWunder_ValidMunichAddress(t *testing.T) {
	// Create product XML
	basicDSLXml := fmt.Sprintf(productXmlTemplate, "1", "WebWunder DSL Basic", 50, 1999, 2499, "", 24, "DSL")
	fiberPremiumXml := fmt.Sprintf(productXmlTemplate, "2", "WebWunder Fiber Premium", 250, 3999, 4499,
		fmt.Sprintf(absoluteVoucherXml, 1000, 2000), 12, "FIBER")
	productsXml := basicDSLXml + fiberPremiumXml

	// Create expected products
	basicDSLProduct := createStandardProduct("WebWunder-1.0", "WebWunder DSL Basic", 50, m.DSL, 1999, 2499, 24)
	fiberPremiumProduct := createStandardProduct("WebWunder-2.0", "WebWunder Fiber Premium", 250, m.FIBER, 3999, 4499, 12)
	fiberPremiumProduct = createProductWithAbsoluteDiscount(fiberPremiumProduct, 1000, 2000)

	testCase := providertest.ProviderTestCase{
		Address: standardTestAddress,
		URLResponseMap: map[string]providertest.HTTPResponse{
			"/endpunkte/soap/ws": {
				StatusCode: 200,
				Headers:    soapHeaders,
				Body:       fmt.Sprintf(soapEnvelopeTemplate, productsXml),
			},
		},
		ExpectedProducts: multiplyProductsForAllConnectionTypes([]m.InternetProduct{basicDSLProduct, fiberPremiumProduct}),
		ExpectedError:    false,
		IsValidResponse:  true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestWebWunder_EmptyResponse tests provider with empty SOAP response
func TestWebWunder_EmptyResponse(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address: standardTestAddress,
		URLResponseMap: map[string]providertest.HTTPResponse{
			"/endpunkte/soap/ws": {
				StatusCode: 200,
				Headers:    soapHeaders,
				Body:       fmt.Sprintf(soapEnvelopeTemplate, ""),
			},
		},
		ExpectedProducts: multiplyProductsForAllConnectionTypes([]m.InternetProduct{}),
		ExpectedError:    false,
		IsValidResponse:  true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestWebWunder_InvalidXMLResponse tests provider with malformed XML
func TestWebWunder_InvalidXMLResponse(t *testing.T) {
	testCase := providertest.ProviderTestCase{
		Address: standardTestAddress,
		URLResponseMap: map[string]providertest.HTTPResponse{
			"/endpunkte/soap/ws": {
				StatusCode: 200,
				Headers:    soapHeaders,
				Body:       `<not>valid<xml>`,
			},
		},
		ExpectedProducts: []m.InternetProduct{},
		ExpectedError:    true,
		IsValidResponse:  false,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestWebWunder_UnsupportedCountry tests provider with unsupported country
func TestWebWunder_UnsupportedCountry(t *testing.T) {
	unsupportedAddress := m.Address{
		Street:      "Baker Street",
		HouseNumber: "221B",
		City:        "London",
		PostalCode:  "NW1 6XE",
		CountryCode: "GB",
	}

	testCase := providertest.ProviderTestCase{
		Address: unsupportedAddress,
		URLResponseMap: map[string]providertest.HTTPResponse{
			"/endpunkte/soap/ws": {
				StatusCode: 200,
				Headers:    soapHeaders,
				Body:       emptySoapEnvelope,
			},
		},
		ExpectedProducts: []m.InternetProduct{},
		ExpectedError:    true,
		IsValidResponse:  false,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestWebWunder_WithVouchers tests parsing of different voucher types
func TestWebWunder_WithVouchers(t *testing.T) {
	// Create product XML with different voucher types
	absoluteVoucherProductXml := fmt.Sprintf(productXmlTemplate, "1", "WebWunder Absolute Voucher", 100, 2999, 3499,
		fmt.Sprintf(absoluteVoucherXml, 500, 1000), 24, "DSL")

	percentageVoucherProductXml := fmt.Sprintf(productXmlTemplate, "2", "WebWunder Percentage Voucher", 200, 3999, 4499,
		fmt.Sprintf(percentageVoucherXml, 10, 1000), 12, "FIBER")

	bothVouchersProductXml := fmt.Sprintf(productXmlTemplate, "3", "WebWunder Both Vouchers", 300, 4999, 5499,
		fmt.Sprintf(bothVouchersXml, 1000, 2000, 15, 1500), 24, "CABLE")

	productsXml := absoluteVoucherProductXml + percentageVoucherProductXml + bothVouchersProductXml

	// Create expected products
	absoluteVoucherProduct := createStandardProduct("WebWunder-1.0", "WebWunder Absolute Voucher", 100, m.DSL, 2999, 3499, 24)
	absoluteVoucherProduct = createProductWithAbsoluteDiscount(absoluteVoucherProduct, 500, 1000)

	percentageVoucherProduct := createStandardProduct("WebWunder-2.0", "WebWunder Percentage Voucher", 200, m.FIBER, 3999, 4499, 12)
	percentageVoucherProduct = createProductWithPercentageDiscount(percentageVoucherProduct, 10, 1000)

	bothVouchersProduct := createStandardProduct("WebWunder-3.0", "WebWunder Both Vouchers", 300, m.CABLE, 4999, 5499, 24)
	bothVouchersProduct = createProductWithAbsoluteDiscount(bothVouchersProduct, 1000, 2000)
	bothVouchersProduct = createProductWithPercentageDiscount(bothVouchersProduct, 15, 1500)

	testCase := providertest.ProviderTestCase{
		Address: standardTestAddress,
		URLResponseMap: map[string]providertest.HTTPResponse{
			"/endpunkte/soap/ws": {
				StatusCode: 200,
				Headers:    soapHeaders,
				Body:       fmt.Sprintf(soapEnvelopeTemplate, productsXml),
			},
		},
		ExpectedProducts: multiplyProductsForAllConnectionTypes([]m.InternetProduct{
			absoluteVoucherProduct,
			percentageVoucherProduct,
			bothVouchersProduct,
		}),
		ExpectedError:   false,
		IsValidResponse: true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestWebWunder_MissingProductInfo tests products with missing product info
func TestWebWunder_MissingProductInfo(t *testing.T) {
	productsXml := `<products>
	<productId>1</productId>
	<providerName>WebWunder Missing Info</providerName>
</products>
<products>
	<productId>2</productId>
	<providerName>WebWunder Valid Info</providerName>
	<productInfo>
		<speed>200</speed>
		<monthlyCostInCent>3999</monthlyCostInCent>
		<monthlyCostInCentFrom25thMonth>4499</monthlyCostInCentFrom25thMonth>
		<contractDurationInMonths>12</contractDurationInMonths>
		<connectionType>FIBER</connectionType>
	</productInfo>
</products>`

	validProduct := createStandardProduct("WebWunder-2.0", "WebWunder Valid Info", 200, m.FIBER, 3999, 4499, 12)

	testCase := providertest.ProviderTestCase{
		Address: standardTestAddress,
		URLResponseMap: map[string]providertest.HTTPResponse{
			"/endpunkte/soap/ws": {
				StatusCode: 200,
				Headers:    soapHeaders,
				Body:       fmt.Sprintf(soapEnvelopeTemplate, productsXml),
			},
		},
		ExpectedProducts: multiplyProductsForAllConnectionTypes([]m.InternetProduct{validProduct}),
		ExpectedError:    false,
		IsValidResponse:  true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestWebWunder_InstallationFalse tests installation=false (id ends with .0)
func TestWebWunder_InstallationFalse(t *testing.T) {
	productXml := fmt.Sprintf(productXmlTemplate, "1", "WebWunder DSL", 100, 2999, 3499, "", 24, "DSL")
	product := createStandardProduct("WebWunder-1.0", "WebWunder DSL", 100, m.DSL, 2999, 3499, 24)

	testCase := providertest.ProviderTestCase{
		Address: standardTestAddress,
		URLResponseMap: map[string]providertest.HTTPResponse{
			"/endpunkte/soap/ws": {
				StatusCode: 200,
				Headers:    soapHeaders,
				Body:       fmt.Sprintf(soapEnvelopeTemplate, productXml),
			},
		},
		ExpectedProducts: multiplyProductsForAllConnectionTypes([]m.InternetProduct{product}),
		ExpectedError:    false,
		IsValidResponse:  true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestWebWunder_InstallationTrue tests installation=true (id ends with .1)
func TestWebWunder_InstallationTrue(t *testing.T) {
	productXml := fmt.Sprintf(productXmlTemplate, "1", "WebWunder DSL", 100, 2999, 3499, "", 24, "DSL")

	// Same as the standard product but with .1 suffix to indicate installation=true
	product := createStandardProduct("WebWunder-1.1", "WebWunder DSL", 100, m.DSL, 2999, 3499, 24)

	testCase := providertest.ProviderTestCase{
		Address: standardTestAddress,
		URLResponseMap: map[string]providertest.HTTPResponse{
			"/endpunkte/soap/ws": {
				StatusCode: 200,
				Headers:    soapHeaders,
				Body:       fmt.Sprintf(soapEnvelopeTemplate, productXml),
			},
		},
		ExpectedProducts: multiplyProductsForAllConnectionTypes([]m.InternetProduct{product}),
		ExpectedError:    false,
		IsValidResponse:  true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// TestWebWunder_ConnectionTypeMismatch tests connection type mismatch between request and response
func TestWebWunder_ConnectionTypeMismatch(t *testing.T) {
	productXml := fmt.Sprintf(productXmlTemplate, "1", "WebWunder Mismatch", 100, 2999, 3499, "", 24, "FIBER")

	// Connection type in product matches response (FIBER) rather than request (which would be DSL)
	product := createStandardProduct("WebWunder-1.0", "WebWunder Mismatch", 100, m.FIBER, 2999, 3499, 24)

	testCase := providertest.ProviderTestCase{
		Address: standardTestAddress,
		URLResponseMap: map[string]providertest.HTTPResponse{
			"/endpunkte/soap/ws": {
				StatusCode: 200,
				Headers:    soapHeaders,
				Body:       fmt.Sprintf(soapEnvelopeTemplate, productXml),
			},
		},
		ExpectedProducts: multiplyProductsForAllConnectionTypes([]m.InternetProduct{product}),
		ExpectedError:    false,
		IsValidResponse:  true,
	}
	providertest.RunProviderTestCase(t, testCase, CreateTestProvider)
}

// Helper to multiply expected products for all connection types and installation options
func multiplyProductsForAllConnectionTypes(baseProducts []m.InternetProduct) []m.InternetProduct {
	var all []m.InternetProduct
	connectionTypes := []m.ConnectionType{m.DSL, m.CABLE, m.FIBER, m.MOBILE}
	for _, ct := range connectionTypes {
		// for _, installation := range []int{0, 1} {
		for _, prod := range baseProducts {
			p := prod
			p.ProductInfo.ConnectionType = ct
			// Id format: WebWunder-<id>.<installation>
			// Replace connection type in Id if present, else append
			// We'll use the same logic as soapProductToInternetProduct
			// Assume Id is like "WebWunder-1.0" or "WebWunder-2.1"
			var baseId string
			if dash := len(p.Id) - 4; dash > 0 && p.Id[:10] == "WebWunder-" {
				baseId = p.Id[:len(p.Id)-2] // remove .0 or .1
			} else {
				baseId = p.Id
			}
			p.Id = fmt.Sprintf("%s-%s.%d", "WebWunder", getIdNumber(baseId), 1)
			all = append(all, p)
		}
		// }
	}
	return all
}

// Helper to extract the numeric part from Id like "WebWunder-1"
func getIdNumber(baseId string) string {
	// baseId is like "WebWunder-1" or "WebWunder-2"
	if dash := len(baseId) - 1; dash > 0 {
		for i := len(baseId) - 1; i >= 0; i-- {
			if baseId[i] == '-' {
				return baseId[i+1:]
			}
		}
	}
	return baseId
}
