package interfaces

import (
	"net/http"

	"github.com/rotmanjanez/check24-gendev-7/pkg/models"
)

// Request represents the information provided by the user
// to the provider adapter.
type Request struct {
	Address models.Address
}

// PreparedRequest represents the requests that are initially sent to the provider based on the data in `Request`.
type PreparedRequest struct {
	// The request to be sent to the provider.
	// Note, that no request is gauaranteed to be sent.
	Request *http.Request

	// Optional Callback which provider adapter is responsible for handling the response
	// if not set, this field will be set to the provider returning the prepared request by the runtime
	Callback ProviderAdapter

	// additional information about the request to be used when parsing the response.
	Metadata interface{}
}
