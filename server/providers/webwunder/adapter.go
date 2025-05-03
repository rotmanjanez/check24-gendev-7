package webwunder

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/rotmanjanez/check24-gendev-7/internal/utils"
	i "github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
	p "github.com/rotmanjanez/check24-gendev-7/pkg/provider"
)

const providerName = "WebWunder"

var connectionTypes = []string{"DSL", "CABLE", "FIBER", "MOBILE"}

type WebWunderAdapter struct {
	apiKey       string
	logger       *slog.Logger
	soapEndpoint string
	soapAction   string
	soapGs       string
	soapEnv      string
}

func init() {
	p.RegisterProvider(providerName, WebWunderFactory)
}

// WebWunderFactory creates a new instance of the WebWunderAdapter
func WebWunderFactory(options map[string]interface{}, cache i.Cache, logger *slog.Logger) (i.ProviderAdapter, error) {
	apiKey := utils.RequireEnv("WEBWUNDER_API_KEY")

	soapEndpoint, ok := options["soapEndpoint"].(string)
	if !ok || soapEndpoint == "" {
		return nil, fmt.Errorf("WebWunder provider requires soapEndpoint in options")
	}

	soapAction, ok := options["soapAction"].(string)
	if !ok || soapAction == "" {
		return nil, fmt.Errorf("WebWunder provider requires soapAction in options")
	}

	soapGs, ok := options["soapGs"].(string)
	if !ok || soapGs == "" {
		return nil, fmt.Errorf("WebWunder provider requires soapGs in options")
	}

	soapEnv, ok := options["soapEnv"].(string)
	if !ok || soapEnv == "" {
		return nil, fmt.Errorf("WebWunder provider requires soapEnv in options")
	}

	for k, v := range options {
		if k != "soapEndpoint" && k != "soapAction" && k != "soapGs" && k != "soapEnv" {
			logger.Warn("WebWunder provider ignoring unknown option", "key", k, "value", v)
		}
	}

	return NewWebWunderAdapter(apiKey, soapEndpoint, soapAction, soapGs, soapEnv, logger), nil
}

func NewWebWunderAdapter(apiKey string, soapEndpoint string, soapAction string, soapGs string, soapEnv string, logger *slog.Logger) *WebWunderAdapter {
	return &WebWunderAdapter{
		apiKey:       apiKey,
		soapEndpoint: soapEndpoint,
		soapAction:   soapAction,
		soapGs:       soapGs,
		soapEnv:      soapEnv,
		logger:       logger,
	}
}

func (w *WebWunderAdapter) Name() string {
	return providerName
}

type metadata struct {
	Installation bool
}

// PrepareRequest converts a general request into SOAP-specific HTTP requests
func (w *WebWunderAdapter) PrepareRequest(ctx context.Context, request i.Request) (i.ParsedResponse, error) {
	switch request.Address.CountryCode {
	case "DE", "AT", "CH":
		break
	default:
		return i.ParsedResponse{}, fmt.Errorf("unsupported country: %s", request.Address.CountryCode)
	}

	if request.Address.HouseNumber == "" {
		w.logger.Debug("No HouseNumber is not supported for WebWunder, skipping request preparation")
		return i.ParsedResponse{}, nil
	}

	var reqests []i.PreparedRequest
	for _, ct := range connectionTypes {
		// installation does not metter for this as price and everything else is the same
		// may change in the future
		// for _, installation := range []bool{false, true} {

		req, err := w.createHTTPRequestFromRequest(request, ct, true)
		if err != nil {
			return i.ParsedResponse{}, fmt.Errorf("error creating HTTP request: %w", err)
		}
		// Create a prepared request with the HTTP request
		reqests = append(reqests, i.PreparedRequest{
			Request: req,
			Metadata: metadata{
				Installation: true, // always true for this provider
			},
		})

		// }
	}

	return i.ParsedResponse{
		Requests: reqests,
	}, nil
}

// ParseResponse parses a SOAP response into offers
func (w *WebWunderAdapter) ParseResponse(ctx context.Context, response i.Response) (i.ParsedResponse, error) {
	httpResponse := response.HTTPResponse

	if httpResponse.StatusCode != http.StatusOK {
		return i.ParsedResponse{}, fmt.Errorf("HTTP request failed with status code: %d", httpResponse.StatusCode)
	}

	// Read the response body
	defer httpResponse.Body.Close()
	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return i.ParsedResponse{}, fmt.Errorf("error reading response body: %w", err)
	}

	w.logger.Debug("Received response", "body", string(body))

	// Parse the SOAP envelope
	var soapResponse SoapResponseEnvelope

	// Use a decoder to handle the SOAP response with appropriate type detection
	decoder := xml.NewDecoder(bytes.NewReader(body))
	if err := decoder.Decode(&soapResponse); err != nil {
		return i.ParsedResponse{}, fmt.Errorf("error decoding SOAP response: %w", err)
	}

	// Convert products to offers
	w.logger.Debug("Parsed response", "products", soapResponse.Body.Output.Products)

	products := make([]m.InternetProduct, 0)
	for _, product := range soapResponse.Body.Output.Products {
		metadata, ok := response.Request.Metadata.(metadata)
		if !ok {
			return i.ParsedResponse{}, fmt.Errorf("installation metadata not found")
		}

		internetProduct, err := soapProductToInternetProduct(product, metadata)
		if err != nil {
			w.logger.Error("Error converting SOAP product to InternetProduct", "error", err)
			continue
		}

		products = append(products, internetProduct)
	}

	return i.ParsedResponse{
		InternetProducts: products,
	}, nil
}

// Helper functions

func soapProductToInternetProduct(product Product, metadata metadata) (m.InternetProduct, error) {
	// Convert SOAP product to InternetProduct
	if product.ProductInfo == nil {
		return m.InternetProduct{}, fmt.Errorf("product info is nil")
	}

	info := m.ProductInfo{}
	pricing := m.Pricing{
		ContractDurationInMonths: &product.ProductInfo.ContractDurationInMonths,
		MonthlyCostInCent:        product.ProductInfo.MonthlyCostInCent,
		SubsequentCosts: &m.SubsequentCost{
			MonthlyCostInCent: product.ProductInfo.MonthlyCostInCentFrom25thMonth,
			StartMonth:        25,
		},
		InstallationServiceIncluded: metadata.Installation,
	}

	ct, err := m.NewConnectionTypeFromValue(product.ProductInfo.ConnectionType)
	if err != nil {
		return m.InternetProduct{}, fmt.Errorf("error creating connection type: %w", err)
	}
	info.ConnectionType = ct
	info.Speed = product.ProductInfo.Speed

	voucher := product.ProductInfo.Voucher
	if voucher != nil && (voucher.AbsoluteVoucher.DiscountInCent != 0 || voucher.AbsoluteVoucher.MinOrderValueInCent != 0) {
		pricing.AbsoluteDiscount = &m.AbsoluteDiscount{
			ValueInCent:         voucher.AbsoluteVoucher.DiscountInCent,
			MinOrderValueInCent: &voucher.AbsoluteVoucher.MinOrderValueInCent,
		}
	}

	if voucher != nil && voucher.PercentageVoucher.Percentage != 0 {
		pricing.PercentageDiscount = &m.PercentageDiscount{
			Percentage:        voucher.PercentageVoucher.Percentage,
			MaxDiscountInCent: &voucher.PercentageVoucher.MaxDiscountInCent,
		}
	}

	installationInt := 0
	if metadata.Installation {
		installationInt = 1
	}

	return m.InternetProduct{
		Id:          fmt.Sprintf("%s-%d.%d", providerName, product.Id, installationInt),
		Provider:    providerName,
		Name:        product.Name,
		ProductInfo: info,
		Pricing:     pricing,
	}, nil
}

// createInputFromRequest converts a generic request to a SOAP-specific input
func createInputFromRequest(request i.Request, connectionType string, installation bool) (Input, error) {
	// Create the input
	input := Input{
		Installation: installation,
		Connection:   connectionType,
		Address: Address{
			Street:      request.Address.Street,
			HouseNumber: request.Address.HouseNumber,
			City:        request.Address.City,
			PLZ:         request.Address.PostalCode,
			CountryCode: string(request.Address.CountryCode),
		},
	}

	return input, nil
}

func (w *WebWunderAdapter) createHTTPRequestFromRequest(request i.Request, connectionType string, installation bool) (*http.Request, error) {
	// Convert the generic request to our SOAP-specific input
	input, err := createInputFromRequest(request, connectionType, installation)
	if err != nil {
		return nil, fmt.Errorf("error creating SOAP input: %w", err)
	}

	// Create the SOAP envelope
	legacyRequest := LegacyGetInternetOffers{
		Input: input,
	}

	envelope := SoapRequestEnvelope{
		Soapenv: w.soapEnv,
		Gs:      w.soapGs,
		Body: SoapRequestBody{
			Content: legacyRequest,
		},
	}

	// Marshal the SOAP envelope to XML
	// todo: test marshal without indent
	xmlData, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling SOAP envelope: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, w.soapEndpoint, bytes.NewReader(xmlData))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Set required headers
	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Set("SOAPAction", w.soapAction)
	req.Header.Set("Accept", "text/xml")
	req.Header.Set("X-Api-Key", w.apiKey)

	return req, nil
}
