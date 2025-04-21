package interfaces

import "context"

type ProviderAdapter interface {
	// Converts a Request into
	// 		1. a slice of http.Requests required to fulfill the request
	// 		2. a slice available internet products that may be already available without sending a request
	//
	// The requests are then sent to the provider by the user of the interface.
	// The responses of these requests are then passed to the ParseResponse method.
	// Even if the function returns an error, all internet products that are returned in the ParsedResponse are considered valid.
	// Preparing needs to be non-blocking (except for cache lookups) to ensure fast response times to the user.
	// Note: The requests are not sent in order.
	PrepareRequest(ctx context.Context, request Request) (ParsedResponse, error)

	// Parses the response of the provider to a request returned by PrepareRequest or ParseResponse.
	// ParseResponse is called once for each request in the slice returned by PrepareRequest or ParseResponse.
	// The functions is expected to return
	// 		1. a slice of valid offers
	// 		2. a slice of follow-up requests that may be necessary to gather all offers from the provider that match a request
	//
	// Even if the function returns an error, all internet products that are returned in the ParsedResponse are considered valid.
	// Note: If the functions returns an error, the Requests in the ParsedResponse will not be further processed.
	ParseResponse(ctx context.Context, response Response) (ParsedResponse, error)

	// Returns the name of the provider.
	Name() string
}
