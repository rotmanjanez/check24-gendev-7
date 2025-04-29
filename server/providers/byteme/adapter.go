package byteme

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/google/go-querystring/query"

	"github.com/rotmanjanez/check24-gendev-7/internal/units"
	"github.com/rotmanjanez/check24-gendev-7/internal/utils"
	i "github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
	p "github.com/rotmanjanez/check24-gendev-7/pkg/provider"
)

const providerName = "ByteMe"

func init() {
	gocsv.FailIfUnmatchedStructTags = true
	gocsv.FailIfDoubleHeaderNames = true

	p.RegisterProvider(providerName, ByteMeFactory)
}

type ByteMeAdapter struct {
	apiKey string
	url    string
	logger *slog.Logger
}

func ByteMeFactory(options map[string]interface{}, cache i.Cache, logger *slog.Logger) (i.ProviderAdapter, error) {
	if len(options) == 0 {
		return nil, fmt.Errorf("ByteMe provider requires url in options")
	}

	url, ok := options["url"].(string)
	if !ok {
		return nil, fmt.Errorf("ByteMe provider requires url in options")
	}

	for k, v := range options {
		if k != "url" {
			logger.Warn("ByteMe provider ignoring unknown option", "key", k, "value", v)
		}
	}

	return NewByteMeAdapter(url, utils.RequireEnv("BYTEME_API_KEY"), logger), nil
}

func NewByteMeAdapter(url string, apiKey string, logger *slog.Logger) *ByteMeAdapter {
	return &ByteMeAdapter{
		apiKey: apiKey,
		url:    url,
		logger: logger,
	}
}

func (*ByteMeAdapter) Name() string {
	return providerName
}

func (b *ByteMeAdapter) PrepareRequest(ctx context.Context, request i.Request) (i.ParsedResponse, error) {
	if request.Address.HouseNumber == "" {
		b.logger.Debug("No HouseNumber is not supported by ByteMe provider")
		return i.ParsedResponse{}, nil
	}

	data := Request{
		Street:      request.Address.Street,
		PostalCode:  request.Address.PostalCode,
		HouseNumber: request.Address.HouseNumber,
		City:        request.Address.City,
	}
	v, err := query.Values(data)
	if err != nil {
		b.logger.Error("Error marshalling query parameters", "error", err)
		return i.ParsedResponse{}, err
	}
	queryParams := v.Encode()

	req, err := http.NewRequest("GET", b.url+"?"+queryParams, nil)
	if err != nil {
		b.logger.Error("Error creating new request", "error", err)
		return i.ParsedResponse{}, err
	}

	req.Header.Set("X-Api-Key", b.apiKey)

	b.logger.Debug("Request", "method", req.Method, "url", req.URL, "queryparams", queryParams)
	return i.ParsedResponse{
		Requests: []i.PreparedRequest{{Request: req}}}, nil
}

func (b *ByteMeAdapter) responseRowToInternetProduct(r ResponseRow) (m.InternetProduct, error) {
	ct, err := m.NewConnectionTypeFromValue(strings.ToUpper(r.ConnectionType))
	if err != nil {
		return m.InternetProduct{}, fmt.Errorf("error creating connection type: %w", err)
	}

	b.logger.Debug("Response", "row", r)

	var unthrottledCapacityMb *int32
	if r.LimitFrom != 0 {
		unthrottledCapacityMb = new(int32)
		*unthrottledCapacityMb = r.LimitFrom * int32(units.Gb/units.Mb)
	}

	var installationService bool
	switch r.InstallationService {
	case "true":
		installationService = true
	case "false":
		installationService = false
	default:
		return m.InternetProduct{}, fmt.Errorf("unknown installation service: %s", r.InstallationService)
	}

	pricing := m.Pricing{
		MonthlyCostInCent:           r.MonthlyCostInCent,
		ContractDurationInMonths:    &r.DurationInMonths,
		MaxAgeInJears:               &r.MaxAge,
		InstallationServiceIncluded: installationService,
		SubsequentCosts: &m.SubsequentCost{
			MonthlyCostInCent: r.AfterTwoYearsMonthlyCostInCent,
			StartMonth:        25,
		},
	}

	switch r.VoucherType {
	case "absolute":
		if r.VoucherValue != 0 {
			pricing.AbsoluteDiscount = &m.AbsoluteDiscount{
				ValueInCent: r.VoucherValue,
			}
		}
	case "percentage":
		if r.VoucherValue != 0 {
			pricing.PercentageDiscount = &m.PercentageDiscount{
				Percentage: r.VoucherValue,
			}
		}
	case "":
		// No voucher
	default:
		return m.InternetProduct{}, fmt.Errorf("unknown voucher type: %s", r.VoucherType)
	}
	name := r.ProviderName
	description := ""

	parts := strings.Split(r.ProviderName, ",")
	if len(parts) == 2 {
		name = parts[0]
		description = parts[1]
	}

	return m.InternetProduct{
		Id:          r.Id,
		Name:        name,
		Description: description,
		Provider:    providerName,
		ProductInfo: m.ProductInfo{
			Speed:                 r.Speed,
			Tv:                    &r.TV,
			ConnectionType:        ct,
			UnthrottledCapacityMb: unthrottledCapacityMb,
		},
		Pricing: pricing,
	}, nil
}

func (b *ByteMeAdapter) ParseResponse(ctx context.Context, resp i.Response) (i.ParsedResponse, error) {
	if resp.HTTPResponse.StatusCode != http.StatusOK {
		b.logger.Error("Error response from server", "statusCode", resp.HTTPResponse.StatusCode)
		return i.ParsedResponse{}, fmt.Errorf("error response from server: %s", resp.HTTPResponse.Status)
	}

	// Read the response body
	var response Response
	err := gocsv.Unmarshal(resp.HTTPResponse.Body, &response.Offers)
	if err != nil {
		b.logger.Error("Error unmarshalling response", "error", err)
		return i.ParsedResponse{}, err
	}
	b.logger.Debug("Parsed response", "offers", response.Offers)

	responseSet := make(map[ResponseRow]bool)
	var InternetProducts []m.InternetProduct

	var errs []error
	for _, offer := range response.Offers {
		if responseSet[offer] {
			continue
		}

		product, err := b.responseRowToInternetProduct(offer)
		if err != nil {
			b.logger.Debug("Error converting response row to internet product", "error", err)
			errs = append(errs, err)
			continue
		}

		InternetProducts = append(InternetProducts, product)
		responseSet[offer] = true
	}

	err = nil
	if len(errs) > 0 {
		err = fmt.Errorf("errors occurred while parsing response: %v", errs)
	}

	return i.ParsedResponse{
		InternetProducts: InternetProducts,
	}, err
}
