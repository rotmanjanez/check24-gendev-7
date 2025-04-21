package interfaces

import (
	"net/http"

	"github.com/rotmanjanez/check24-gendev-7/pkg/models"
)

// Response represents the response from a provider to a request issued by the provider adapter.
// It contains the initial request data, the prepared request that was sent to the provider and the HTTP response from the provider.
type Response struct {
	InitialRequestData Request
	Request            PreparedRequest
	HTTPResponse       *http.Response
}

// ParsedResponse represents the parsed response from a provider.
// It contains a slice of offers and a slice of follow-up requests that may be necessary to gather all offers from the provider that match a request.
type ParsedResponse struct {
	InternetProducts []models.InternetProduct
	Requests         []PreparedRequest
}
