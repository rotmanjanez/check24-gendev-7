package exampleprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	i "github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
	p "github.com/rotmanjanez/check24-gendev-7/pkg/provider"
)

const providerName = "Example Provider"

func init() {
	p.RegisterProvider(providerName, ExampleProviderFactory)
}

type ExampleProvider struct {
	responses []m.InternetProduct
	delay     float64
	logger    *slog.Logger
}

func ExampleProviderFactory(options map[string]interface{}, cache i.Cache, logger *slog.Logger) (i.ProviderAdapter, error) {
	// Here you can initialize your provider with the options provided in the config
	delay, ok := options["delay"].(float64)
	if !ok {
		delay = 0 // Default delay if not provided
		logger.Warn("No delay provided, using default of 0 seconds")
	}

	responsesRaw, ok := options["responses"]
	if !ok {
		return nil, fmt.Errorf("no responses provided")
	}
	// parse the responses as json
	jsonData, err := json.Marshal(responsesRaw)
	if err != nil {
		return nil, fmt.Errorf("error marshalling responses: %w", err)
	}
	var responses []m.InternetProduct
	err = json.Unmarshal(jsonData, &responses)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling responses: %w", err)
	}
	if len(responses) == 0 {
		logger.Warn("No responses provided, using default empty response")
		responses = []m.InternetProduct{}
	}

	return NewExampleProvider(responses, delay, logger)
}

func NewExampleProvider(responses []m.InternetProduct, delay float64, logger *slog.Logger) (*ExampleProvider, error) {
	return &ExampleProvider{
		responses: responses,
		delay:     delay,
		logger:    logger,
	}, nil
}

func (*ExampleProvider) Name() string {
	return providerName
}

func (p *ExampleProvider) PrepareRequest(ctx context.Context, request i.Request) (i.ParsedResponse, error) {
	// Here you can modify the request as needed before sending it to the provider
	// For example, you might want to add headers or change the URL

	// sleep for the specified delay
	time.Sleep(time.Duration(p.delay) * time.Second)

	return i.ParsedResponse{
		InternetProducts: p.responses,
	}, nil
}

func (p *ExampleProvider) ParseResponse(ctx context.Context, response i.Response) (i.ParsedResponse, error) {
	// Here you can parse the response from the provider and extract the relevant data from a request initiated in a previous call to PrepareRequest or ParseResponse

	return i.ParsedResponse{}, fmt.Errorf("example provider does not make any http requests and therefore does not parse any responses")
}
