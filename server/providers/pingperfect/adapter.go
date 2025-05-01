package pingperfect

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rotmanjanez/check24-gendev-7/internal/units"
	utils "github.com/rotmanjanez/check24-gendev-7/internal/utils"
	i "github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
	p "github.com/rotmanjanez/check24-gendev-7/pkg/provider"
)

const providerName = "PingPerfect"

func init() {
	p.RegisterProvider(providerName, PingPerfectFactory)
}

type PingPerfectAdapter struct {
	signatureSecret string
	clientId        string
	hasher          hash.Hash
	url             string
	logger          *slog.Logger
}

func NewPingPerfectAdapter(url string, clientId string, signatureSecret string, logger *slog.Logger) *PingPerfectAdapter {
	return &PingPerfectAdapter{
		signatureSecret: signatureSecret,
		clientId:        clientId,
		hasher:          hmac.New(sha256.New, []byte(signatureSecret)),
		url:             url,
		logger:          logger,
	}
}

func PingPerfectFactory(options map[string]interface{}, cache i.Cache, logger *slog.Logger) (i.ProviderAdapter, error) {
	if len(options) == 0 {
		return nil, fmt.Errorf("PingPerfect provider requires url in options")
	}

	url, ok := options["url"].(string)
	if !ok {
		return nil, fmt.Errorf("PingPerfect provider requires url in options")
	}

	for k, v := range options {
		if k != "url" {
			logger.Warn("ByteMe provider ignoring unknown option", "key", k, "value", v)
		}
	}

	signatureSecret := utils.RequireEnv("PING_PERFECT_SIGNATURE_SECRET")
	clientId := utils.RequireEnv("PING_PERFECT_CLIENT_ID")

	return NewPingPerfectAdapter(url, clientId, signatureSecret, logger), nil
}

func (p *PingPerfectAdapter) Name() string {
	return providerName
}

type Signature struct {
	Signature string
	Timestamp int64
}

func (p *PingPerfectAdapter) getSignature(body []byte) (Signature, error) {
	defer p.hasher.Reset()

	timeStamp := time.Now().Unix()

	str := fmt.Sprintf("%d:%s", timeStamp, string(body))

	_, err := p.hasher.Write([]byte(str))
	if err != nil {
		p.logger.Error("Error writing to hasher", "error", err)
		return Signature{}, err
	}

	signature := p.hasher.Sum(nil)

	return Signature{
		Signature: hex.EncodeToString(signature),
		Timestamp: timeStamp,
	}, nil
}

// Converts an address and wantsFiber flag into a request in the provider's format
func (p *PingPerfectAdapter) getRequestBody(address m.Address, wantsFiber bool) ([]byte, error) {
	if address.HouseNumber == "" {
		// There are addresses without a house number. PingPerfect does not support these.
		return nil, nil
	}

	houseNumber, err := strconv.ParseInt(address.HouseNumber, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("error converting house number to int: %w", err)
	}

	data := Request{
		Street:      address.Street,
		PostalCode:  address.PostalCode,
		HouseNumber: int32(houseNumber),
		City:        address.City,
		WantsFiber:  wantsFiber,
	}

	body, err := json.Marshal(data)
	if err != nil {
		p.logger.Error("Error marshalling request", "error", err, "request", data)
		return nil, err
	}
	return body, nil
}

// Adds required headers to the request
func (p *PingPerfectAdapter) addHeaders(req *http.Request, signature Signature) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Client-ID", p.clientId)
	req.Header.Set("X-Signature", signature.Signature)
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", signature.Timestamp))
}

func (p *PingPerfectAdapter) prepareRequest(address i.Request, wantsFiber bool) (*i.PreparedRequest, error) {
	body, err := p.getRequestBody(address.Address, wantsFiber)

	if err != nil {
		p.logger.Error("Error getting request body", "error", err)
		return nil, err
	}
	if body == nil {
		p.logger.Debug("No house number provided, skipping request")
		return nil, nil
	}

	signature, err := p.getSignature(body)
	if err != nil {
		p.logger.Error("Error getting signature", "error", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", p.url, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		p.logger.Error("Error creating request", "error", err)
		return nil, err
	}

	// Add headers to the request
	p.addHeaders(req, signature)
	return &i.PreparedRequest{
		Request: req,
	}, nil

}

func (p *PingPerfectAdapter) PrepareRequest(ctx context.Context, request i.Request) (i.ParsedResponse, error) {
	// Prepare the request
	preparedRequest, err := p.prepareRequest(request, false)
	if err != nil {
		p.logger.Error("Error preparing request", "error", err)
		return i.ParsedResponse{}, err
	}
	if preparedRequest == nil {
		p.logger.Debug("No request prepared, skipping")
		return i.ParsedResponse{}, nil
	}

	return i.ParsedResponse{
		Requests: []i.PreparedRequest{*preparedRequest},
	}, nil
}

func (p *PingPerfectAdapter) ParseResponse(ctx context.Context, response i.Response) (i.ParsedResponse, error) {
	httpResponse := response.HTTPResponse
	if httpResponse.StatusCode != http.StatusOK {
		p.logger.Error("Error response", "status", httpResponse.StatusCode)
		return i.ParsedResponse{}, fmt.Errorf("error response: %s", httpResponse.Status)
	}

	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		p.logger.Error("Error reading response body", "error", err)
		return i.ParsedResponse{}, err
	}

	// Unmarshal the response
	var offers []InternetProduct

	err = json.Unmarshal(responseBody, &offers)
	if err != nil {
		p.logger.Error("Error unmarshalling response", "error", err, "response", string(responseBody))
		return i.ParsedResponse{}, err
	}
	p.logger.Debug("Parsed response", "products", offers)

	var internetProducts []m.InternetProduct

	var errs []error
	for _, offer := range offers {
		internetProduct, err := p.offerToInternetProduct(offer)
		if err != nil {
			p.logger.Debug("Error converting offer to internet product", "offer", offer, "error", err)
			errs = append(errs, err)
			continue
		}

		internetProducts = append(internetProducts, internetProduct)
	}

	err = nil
	if len(errs) > 0 {
		err = fmt.Errorf("errors occurred while parsing response: %v", errs)
	}

	return i.ParsedResponse{
		InternetProducts: internetProducts,
	}, err
}

func (p *PingPerfectAdapter) offerToInternetProduct(offer InternetProduct) (m.InternetProduct, error) {
	ct, err := m.NewConnectionTypeFromValue(strings.ToUpper(offer.ProductInfo.ConnectionType))
	if err != nil {
		return m.InternetProduct{}, fmt.Errorf("error creating connection type: %w", err)
	}

	var installationServiceIncluded bool
	switch offer.PricingDetails.InstallationService {
	case "yes":
		installationServiceIncluded = true
	case "no":
		installationServiceIncluded = false
	default:
		return m.InternetProduct{}, fmt.Errorf("unknown installation service: %s", offer.PricingDetails.InstallationService)
	}

	var unthrottledCapacityMb *int32
	if offer.ProductInfo.LimitFrom > 0 {
		unthrottledCapacityMb = new(int32)
		*unthrottledCapacityMb = offer.ProductInfo.LimitFrom * int32(units.Gb/units.Mb)
	}

	return m.InternetProduct{
		Id:       offer.ProviderName,
		Provider: providerName,
		Name:     offer.ProviderName,
		ProductInfo: m.ProductInfo{
			Speed:                 offer.ProductInfo.Speed,
			Tv:                    &offer.ProductInfo.Tv,
			ConnectionType:        ct,
			UnthrottledCapacityMb: unthrottledCapacityMb,
		},
		Pricing: m.Pricing{
			MonthlyCostInCent:           offer.PricingDetails.MonthlyCostInCent,
			ContractDurationInMonths:    &offer.ProductInfo.ContractDurationInMonths,
			MaxAgeInJears:               &offer.ProductInfo.MaxAge,
			InstallationServiceIncluded: installationServiceIncluded,
		},
	}, nil
}
