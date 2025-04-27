/*
 * CHECK24 GenDev 7 API
 *
 * API for the 7th CHECK24 GenDev challenge providing product offerings from five different internet providers
 *
 * API version: dev
 */

package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/rotmanjanez/check24-gendev-7/config"
	i "github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
)

// InternetProductsAPIService is a service that implements the logic for the InternetProductsAPIServicer
// This service should implement the business logic for every endpoint for the InternetProductsAPI API.
// Include any external packages or services that will be required by this service.
type InternetProductsAPIService struct {
	config    *config.Config
	providers []i.ProviderAdapter
}

// NewInternetProductsAPIService creates a default api service
func NewInternetProductsAPIService(cfg *config.Config, providers []i.ProviderAdapter) *InternetProductsAPIService {
	return &InternetProductsAPIService{
		config:    cfg,
		providers: providers,
	}
}

func handleAdaperResponse(products *[]m.InternetProduct, registeredRequestsCtx *[]i.Response, initialRequest i.Request, response i.ParsedResponse, provider i.ProviderAdapter) {
	// still continue with the returned requests and products
	*products = append(*products, response.InternetProducts...)

	for _, prepReq := range response.Requests {
		// set the callback provider to the one that prepared the request if not already set
		if prepReq.Callback == nil {
			prepReq.Callback = provider
		}

		// add the request to the list of requests to be sent
		*registeredRequestsCtx = append(*registeredRequestsCtx, i.Response{
			InitialRequestData: initialRequest,
			Request:            prepReq,
			HTTPResponse:       nil,
		})
	}
}

// GetInternetProducts -
func (s *InternetProductsAPIService) GetInternetProducts(ctx context.Context, address m.Address) (ImplResponse, error) {

	var products []m.InternetProduct

	// Response also contains the initial Request data an the prepared requests returned from the provider adapter as it may be needed to parse the response
	// So we con use it directly to store all required information, even if the naming is a bit confusing
	var registeredRequestsCtx []i.Response

	for _, provider := range s.providers {
		reqData := i.Request{
			Address: address,
		}

		resp, err := provider.PrepareRequest(reqData)

		if err != nil {
			slog.Error("Error preparing request", "provider", provider.Name(), "error", err)
		}

		handleAdaperResponse(&products, &registeredRequestsCtx, reqData, resp, provider)
	}

	for len(registeredRequestsCtx) > 0 {
		requestCtx := registeredRequestsCtx[0]
		registeredRequestsCtx = registeredRequestsCtx[1:]

		httpReq := requestCtx.Request.Request

		client := &http.Client{}

		// do the request
		httpResp, err := client.Do(httpReq)

		if err != nil {
			slog.Error("Error doing request", "request", httpReq, "error", err)
			continue
		}
		defer httpResp.Body.Close()

		requestCtx.HTTPResponse = httpResp

		callbackProvider := requestCtx.Request.Callback

		parsedResp, err := callbackProvider.ParseResponse(requestCtx)

		if err != nil {
			slog.Error("Error parsing response", "provider", callbackProvider.Name(), "error", err)
		}

		handleAdaperResponse(&products, &registeredRequestsCtx, requestCtx.InitialRequestData, parsedResp, callbackProvider)
	}

	return Response(200, products), nil
}
