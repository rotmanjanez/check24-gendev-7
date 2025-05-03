package verbyndich

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/rotmanjanez/check24-gendev-7/internal/units"
	utils "github.com/rotmanjanez/check24-gendev-7/internal/utils"
	i "github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
	p "github.com/rotmanjanez/check24-gendev-7/pkg/provider"
)

const providerName = "VerbynDich"

type VerbynDichAdapter struct {
	url               string
	apiKey            string
	descriptionParser *DescriptionParser
	blockSize         uint
	logger            *slog.Logger
}

func init() {
	p.RegisterProvider(providerName, VerbynDichFactory)
}

func VerbynDichFactory(options map[string]interface{}, cache i.Cache, logger *slog.Logger) (i.ProviderAdapter, error) {
	url, ok := options["url"].(string)
	if !ok {
		return nil, fmt.Errorf("VerbynDich provider requires url in options")
	}
	blockSize, ok := options["blockSize"].(float64)
	if !ok {
		return nil, fmt.Errorf("VerbynDich provider requires blockSize in options")
	}

	for k, v := range options {
		if k != "url" && k != "blockSize" {
			logger.Warn("VerbynDich provider ignoring unknown option", "key", k, "value", v)
		}
	}
	return NewVerbynDichAdapter(url, utils.RequireEnv("VERBYNDICH_API_KEY"), uint(blockSize), logger), nil
}

func NewVerbynDichAdapter(url string, apiKey string, blockSize uint, logger *slog.Logger) *VerbynDichAdapter {
	return &VerbynDichAdapter{
		url:               url,
		apiKey:            apiKey,
		descriptionParser: NewDescriptionParser(logger),
		blockSize:         blockSize,
		logger:            logger,
	}
}

func (v *VerbynDichAdapter) Name() string {
	return providerName
}

func (v *VerbynDichAdapter) newAPIRequest(address m.Address, page uint) (i.PreparedRequest, error) {
	url := fmt.Sprintf("%s/check24/data?apiKey=%s&page=%d", v.url, v.apiKey, page)
	body := fmt.Sprintf(`%s;%s;%s;%s`, address.Street, address.HouseNumber, address.City, address.PostalCode)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		v.logger.Error("Error creating new request", "error", err)
		return i.PreparedRequest{}, err
	}

	req.Header.Set("Accept", "application/json")

	return i.PreparedRequest{
		Request:  req,
		Metadata: page,
	}, nil
}

func (v *VerbynDichAdapter) PrepareRequest(ctx context.Context, request i.Request) (i.ParsedResponse, error) {
	if request.Address.HouseNumber == "" {
		v.logger.Debug("No HouseNumber is not supported by VerbynDich provider")
		return i.ParsedResponse{}, nil
	}

	var requests []i.PreparedRequest
	for idx := uint(0); idx < v.blockSize; idx++ {
		req, err := v.newAPIRequest(request.Address, idx)
		if err != nil {
			return i.ParsedResponse{}, fmt.Errorf("error creating new request: %w", err)
		}
		requests = append(requests, req)
	}

	return i.ParsedResponse{Requests: requests}, nil
}

func (v *VerbynDichAdapter) ParseResponse(ctx context.Context, resp i.Response) (i.ParsedResponse, error) {
	status := resp.HTTPResponse.StatusCode
	if status != http.StatusOK {
		v.logger.Error("Error response from server", "statusCode", status)
		strinBody, err := io.ReadAll(resp.HTTPResponse.Body)
		if err != nil {
			v.logger.Error("Error reading response body", "error", err)
			return i.ParsedResponse{}, fmt.Errorf("error reading response body: %w", err)
		}
		v.logger.Debug("Response body", "body", string(strinBody))
		return i.ParsedResponse{}, fmt.Errorf("error response from server: %d", status)
	}

	body, err := io.ReadAll(resp.HTTPResponse.Body)
	if err != nil {
		v.logger.Error("Error reading response body", "error", err)
		return i.ParsedResponse{}, err
	}

	var parsedData Response

	err = json.Unmarshal(body, &parsedData)
	if err != nil {
		v.logger.Error("Error unmarshalling response", "error", err)
		return i.ParsedResponse{}, err
	}

	v.logger.Debug("provider response", "response", parsedData)

	parsedResponse := i.ParsedResponse{}

	if parsedData.Valid {
		data, err := v.descriptionParser.parse(parsedData.Description)
		if err != nil {
			v.logger.Error("Error parsing description", "error", err)
			return i.ParsedResponse{}, err
		}
		v.logger.Debug("Parsed description data", "data", data)

		product, err := v.productToInternetProduct(parsedData.Product, data)
		if err != nil {
			v.logger.Error("Error converting product to InternetProduct", "error", err)
			return i.ParsedResponse{}, err
		}
		v.logger.Debug("Parsed product", "product", product)

		parsedResponse.InternetProducts = []m.InternetProduct{product}

		if !parsedData.Last {
			prevPage := resp.Request.Metadata.(uint)

			req, err := v.newAPIRequest(resp.InitialRequestData.Address, prevPage+v.blockSize)
			if err != nil {
				v.logger.Error("Error creating new request", "error", err)
				return i.ParsedResponse{}, err
			}
			parsedResponse.Requests = []i.PreparedRequest{req}
		}
	}

	return parsedResponse, nil
}

func (v *VerbynDichAdapter) productToInternetProduct(product string, data DescriptionData) (m.InternetProduct, error) {
	if data.Speed == nil {
		return m.InternetProduct{}, fmt.Errorf("speed is nil")
	}
	var speed int32
	switch data.Speed.Unit {
	case "Mbit/s":
		speed = data.Speed.Value
	case "Gbit/s":
		speed = data.Speed.Value * int32(units.Gb/units.Mb)
	default:
		return m.InternetProduct{}, fmt.Errorf("unknown speed unit: %s", data.Speed.Unit)
	}

	if data.Price == nil {
		return m.InternetProduct{}, fmt.Errorf("price is nil")
	}
	priceInCent := *data.Price * 100

	var tv *string
	if data.IncludedTVSender != "" {
		tv = &data.IncludedTVSender
	}

	var unthrottledCapacity *int32
	if data.UnthrottledCapacity != nil {
		var scale int32
		switch strings.ToUpper(data.UnthrottledCapacity.Unit) {
		case "MB":
			scale = int32(units.Mb)
		case "GB":
			scale = int32(units.Gb / units.Mb)
		default:
			return m.InternetProduct{}, fmt.Errorf("unknown unthrottled capacity unit: %s", data.UnthrottledCapacity.Unit)
		}

		unthrottledCapacity = new(int32)
		*unthrottledCapacity = data.UnthrottledCapacity.Value * scale
	}

	pricing := m.Pricing{
		MaxAgeInJears:               data.MaxAge,
		MinAgeInYears:               data.MinAge,
		MonthlyCostInCent:           priceInCent,
		SubsequentCosts:             data.SubsequentCost,
		InstallationServiceIncluded: data.InstallationIncluded,
		AbsoluteDiscount:            data.AbsoluteDiscount,
		PercentageDiscount:          data.PercentageDiscount,
	}

	if data.MinOrderValue != nil {
		pricing.MinOrderValueInCent = new(int32)
		*pricing.MinOrderValueInCent = *data.MinOrderValue * units.Eur
	}

	if data.MinimalContractDuration != nil {
		minContractDurationInMonths := new(int32)
		switch strings.ToUpper(data.MinimalContractDuration.Unit) {
		case "MONAT", "MONATE":
			*minContractDurationInMonths = data.MinimalContractDuration.Value
		case "JAHRE", "JAHREN":
			*minContractDurationInMonths = data.MinimalContractDuration.Value * 12
		default:
			return m.InternetProduct{}, fmt.Errorf("unknown contract duration unit: %s", data.MinimalContractDuration.Unit)
		}
		pricing.MinContractDurationInMonths = minContractDurationInMonths
	}

	return m.InternetProduct{
		Id:       product,
		Name:     product,
		Provider: providerName,
		ProductInfo: m.ProductInfo{
			Speed:                 speed,
			ConnectionType:        data.ConnectionType,
			Tv:                    tv,
			UnthrottledCapacityMb: unthrottledCapacity,
		},
		Pricing: pricing,
	}, nil
}
