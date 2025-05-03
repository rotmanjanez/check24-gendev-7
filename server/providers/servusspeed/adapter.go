package servusspeed

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/rotmanjanez/check24-gendev-7/internal/units"
	"github.com/rotmanjanez/check24-gendev-7/internal/utils"
	i "github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
	p "github.com/rotmanjanez/check24-gendev-7/pkg/provider"
)

const providerName = "ServusSpeed"

func init() {
	p.RegisterProvider(providerName, ServusSpeedFactory)
}

func ServusSpeedFactory(options map[string]interface{}, cache i.Cache, logger *slog.Logger) (i.ProviderAdapter, error) {
	if len(options) == 0 {
		return nil, fmt.Errorf("ServusSpeed provider requires url in options")
	}
	url, ok := options["url"].(string)
	if !ok {
		return nil, fmt.Errorf("ServusSpeed provider requires url in options")
	}
	cacheDuration, ok := options["cacheDuration"].(float64)
	if !ok {
		cacheDuration = 5 // Default cache duration if not provided
		logger.Warn("No cache duration provided, using default of 5 seconds")
	}
	for k, v := range options {
		if k != "url" && k != "cacheDuration" {
			logger.Warn("ServusSpeed provider ignoring unknown option", "key", k, "value", v)
		}
	}
	return NewServusSpeedAdapter(
		utils.RequireEnv("SERVUS_SPEED_USERNAME"),
		utils.RequireEnv("SERVUS_SPEED_PASSWORD"),
		url,
		cache,
		time.Duration(cacheDuration)*time.Minute,
		logger,
	), nil
}

type ServusSpeedAdapter struct {
	username      string
	password      string
	url           string
	cacheDuration time.Duration
	cache         i.Cache
	logger        *slog.Logger
}

func NewServusSpeedAdapter(username string, password string, url string, cache i.Cache, cacheDuration time.Duration, logger *slog.Logger) *ServusSpeedAdapter {
	return &ServusSpeedAdapter{
		username:      username,
		password:      password,
		url:           url,
		cache:         cache,
		cacheDuration: cacheDuration,
		logger:        logger,
	}
}

func (*ServusSpeedAdapter) Name() string {
	return providerName
}

func (s *ServusSpeedAdapter) newAPIRequest(method, endpoint string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, s.url+endpoint, bytes.NewBuffer(body))
	if err != nil {
		s.logger.Error("Error creating new request", "error", err)
		return nil, err
	}
	// add basic auth
	req.SetBasicAuth(s.username, s.password)

	// set headers
	req.Header.Set("Content-Type", "application/json")

	s.logger.Debug("New request", "method", method, "url", req.URL.String(), "body", string(body))

	return req, nil
}

func (s *ServusSpeedAdapter) convertDEAddressToServusSpeed(address m.Address) Address {
	if address.CountryCode != "DE" {
		s.logger.Debug("Servus Speed only supports Germany as a country")
		return Address{}
	}

	return Address{
		Street:      address.Street,
		HouseNumber: address.HouseNumber,
		City:        address.City,
		PostalCode:  address.PostalCode,
		Country:     "DE",
	}
}

func (s *ServusSpeedAdapter) PrepareRequest(ctx context.Context, request i.Request) (i.ParsedResponse, error) {
	if request.Address.CountryCode != "DE" {
		s.logger.Debug("Servus Speed only supports Germany as a country")
		return i.ParsedResponse{}, nil
	}

	if request.Address.HouseNumber == "" {
		s.logger.Debug("Servus Speed requires a house number")
		return i.ParsedResponse{}, nil
	}

	requestAddress := s.convertDEAddressToServusSpeed(request.Address)

	data := AvailableProductsRequest{
		Address: requestAddress,
	}

	v, err := json.Marshal(data)
	if err != nil {
		s.logger.Error("Error marshalling request data", "error", err)
		return i.ParsedResponse{}, err
	}
	req, err := s.newAPIRequest("POST", "/api/external/available-products", v)
	if err != nil {
		return i.ParsedResponse{}, err
	}

	return i.ParsedResponse{
		Requests: []i.PreparedRequest{{Request: req}},
	}, nil
}

func (s *ServusSpeedAdapter) parseProductDetailsResponse(ctx context.Context, id string, body []byte) (i.ParsedResponse, error) {
	var response ProductDetailsResponse

	err := json.Unmarshal(body, &response)
	if err != nil {

		return i.ParsedResponse{}, fmt.Errorf("error unmarshalling response: %w", err)
	}

	product := response.ServusSpeedProduct

	ct, err := m.NewConnectionTypeFromValue(strings.ToUpper(product.ProductInfo.ConnectionType))
	if err != nil {
		return i.ParsedResponse{}, fmt.Errorf("error creating connection type: %w", err)
	}
	var tv *string
	if product.ProductInfo.TV != "" {
		tv = &product.ProductInfo.TV
	}

	var maxAgeInYears *int32
	if product.ProductInfo.MaxAge > 0 {
		maxAgeInYears = &product.ProductInfo.MaxAge
	}

	var contractDurationInMonths *int32
	if product.ProductInfo.ContractDurationInMonths > 0 {
		contractDurationInMonths = &product.ProductInfo.ContractDurationInMonths
	}

	var unthrottledCapacityMb *int32
	if product.ProductInfo.LimitFrom > 0 {
		unthrottledCapacityMb = &product.ProductInfo.LimitFrom
		*unthrottledCapacityMb = *unthrottledCapacityMb * int32(units.Gb)
	}

	pricing := m.Pricing{
		MonthlyCostInCent:           product.PricingDetails.MonthlyCostInCent,
		MaxAgeInJears:               maxAgeInYears,
		ContractDurationInMonths:    contractDurationInMonths,
		InstallationServiceIncluded: product.PricingDetails.InstallationService,
	}

	if product.Discount != 0 {
		pricing.AbsoluteDiscount = &m.AbsoluteDiscount{
			ValueInCent: product.Discount,
		}
	}

	internetProduct := m.InternetProduct{
		Id:       id,
		Provider: providerName,
		Name:     product.ProviderName,
		/* DateOffered */
		ProductInfo: m.ProductInfo{
			Speed:                 product.ProductInfo.Speed,
			ConnectionType:        ct,
			Tv:                    tv,
			UnthrottledCapacityMb: unthrottledCapacityMb,
		},
		Pricing: pricing,
	}

	err = s.cache.Set(ctx, internetProduct.Id, &internetProduct, 5*time.Minute)
	if err != nil {
		s.logger.Error("Error setting product in cache", "error", err)
	}

	return i.ParsedResponse{
		InternetProducts: []m.InternetProduct{internetProduct},
	}, nil
}

func (s *ServusSpeedAdapter) parseAvailableProductsResponse(ctx context.Context, resp i.Response, body []byte) (i.ParsedResponse, error) {
	var response AvailableProductsResponse

	err := json.Unmarshal(body, &response)
	if err != nil {
		return i.ParsedResponse{}, fmt.Errorf("error unmarshalling response: %w", err)
	}

	s.logger.Debug("Parsed response", "response", response)

	var followUpRequests []i.PreparedRequest
	var products []m.InternetProduct

	for _, product := range response.Products {
		internetProduct := new(m.InternetProduct)
		found, err := s.cache.Get(ctx, product, internetProduct)
		if err == nil && found {
			products = append(products, *internetProduct)
			continue
		} else if err != nil {
			s.logger.Error("error getting product from cache", "product", product, "error", err)
		} else {
			s.logger.Debug("Product not found in cache", "product", product)
		}

		requestAddress := s.convertDEAddressToServusSpeed(resp.InitialRequestData.Address)

		productDetailsRequest := ProductDetailsRequest{
			Address: requestAddress,
		}

		productDetailsRequestBody, err := json.Marshal(productDetailsRequest)
		if err != nil {
			s.logger.Error("Error marshalling product details request", "error", err)
			return i.ParsedResponse{}, err
		}
		followUpRequest, err := s.newAPIRequest("POST", "/api/external/product-details/"+product, productDetailsRequestBody)
		if err != nil {
			return i.ParsedResponse{}, err
		}

		followUpRequests = append(followUpRequests, i.PreparedRequest{
			Request:  followUpRequest,
			Metadata: product,
		})
	}

	return i.ParsedResponse{
		InternetProducts: products,
		Requests:         followUpRequests,
	}, nil
}

func (s *ServusSpeedAdapter) ParseResponse(ctx context.Context, resp i.Response) (i.ParsedResponse, error) {
	if resp.HTTPResponse.StatusCode != http.StatusOK {
		s.logger.Error("Error response from server", "statusCode", resp.HTTPResponse.StatusCode)
		return i.ParsedResponse{}, nil
	}

	// Read the response body
	body, err := io.ReadAll(resp.HTTPResponse.Body)
	if err != nil {
		s.logger.Error("Error reading response body", "error", err)
		return i.ParsedResponse{}, err
	}
	s.logger.Debug("Response body", "body", string(body))

	path := resp.HTTPResponse.Request.URL.Path

	if path == "/api/external/available-products" {
		return s.parseAvailableProductsResponse(ctx, resp, body)
	} else if strings.HasPrefix(path, "/api/external/product-details/") {
		id, ok := resp.Request.Metadata.(string)
		if !ok {
			return i.ParsedResponse{}, fmt.Errorf("error getting product ID from metadata")
		}
		return s.parseProductDetailsResponse(ctx, id, body)
	} else {
		s.logger.Error("Unknown endpoint", "path", resp.HTTPResponse.Request.URL.Path)
		return i.ParsedResponse{}, nil
	}
}
