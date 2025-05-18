package requestmanager

import (
	"context"
	"testing"
	"time"

	i "github.com/rotmanjanez/check24-gendev-7/pkg/interfaces"
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
	p "github.com/rotmanjanez/check24-gendev-7/pkg/provider"
)

// fakeAdapter implements ProviderAdapter for testing PrepareRequest only
type fakeAdapter struct {
	prepareResp i.ParsedResponse
	parseResp   i.ParsedResponse
	prepareErr  error
	parseErr    error
}

func (f *fakeAdapter) PrepareRequest(ctx context.Context, req i.Request) (i.ParsedResponse, error) {
	return f.prepareResp, f.prepareErr
}

func (f *fakeAdapter) ParseResponse(ctx context.Context, resp i.Response) (i.ParsedResponse, error) {
	return f.parseResp, f.parseErr
}

func (f *fakeAdapter) Name() string {
	return "fake"
}

// newProvider creates a ProviderConfig with the fake adapter
func newProvider(adapter i.ProviderAdapter) *p.ProviderConfig {
	// retries=0, timeout minimal, concurrency=1, backoff minimal
	return p.NewProviderConfig(adapter, 0, time.Millisecond, 1, time.Millisecond)
}

// collectChannels drains response and error channels with timeout
func collectChannels(responses <-chan m.InternetProduct, errs <-chan error) ([]m.InternetProduct, []error) {
	var got []m.InternetProduct
	var errsCollected []error
	for responses != nil || errs != nil {
		select {
		case r, ok := <-responses:
			if !ok {
				responses = nil
				break
			}
			got = append(got, r)
		case e, ok := <-errs:
			if !ok {
				errs = nil
				break
			}
			errsCollected = append(errsCollected, e)
		case <-time.After(50 * time.Millisecond):
			return got, errsCollected
		}
	}
	return got, errsCollected
}

// helpers for test values
func int32p(i int32) *int32 { return &i }
func pricing(cost, duration int32) m.Pricing {
	return m.Pricing{MonthlyCostInCent: cost, ContractDurationInMonths: int32p(duration)}
}
func info(speed int32, connType m.ConnectionType) m.ProductInfo {
	return m.ProductInfo{Speed: speed, ConnectionType: connType}
}

// Test valid product is returned with no errors
func TestSingleProviderValid(t *testing.T) {
	prod := m.InternetProduct{
		Id:          "1",
		Provider:    "p",
		Name:        "prod",
		DateOffered: time.Now(),
		ProductInfo: info(100, m.FIBER),
		Pricing:     pricing(1000, 12),
	}
	adapter := &fakeAdapter{prepareResp: i.ParsedResponse{InternetProducts: []m.InternetProduct{prod}}}
	coord := NewRequestCoordinator([]*p.ProviderConfig{newProvider(adapter)})
	res, errs := coord.Run(context.Background(), i.Request{}, 1, 1)
	out, errsOut := collectChannels(res, errs)
	if len(errsOut) != 0 {
		t.Errorf("expected no errors, got %v", errsOut)
	}
	if len(out) != 1 || out[0].Id != prod.Id {
		t.Errorf("expected product %v, got %v", prod, out)
	}
}

// Test missing required fields produces error and no response
func TestRequiredFieldError(t *testing.T) {
	// missing Id
	bad := m.InternetProduct{Provider: "p", Name: "prod", DateOffered: time.Now(), ProductInfo: info(1, m.DSL), Pricing: pricing(1, 1)}
	adapter := &fakeAdapter{prepareResp: i.ParsedResponse{InternetProducts: []m.InternetProduct{bad}}}
	coord := NewRequestCoordinator([]*p.ProviderConfig{newProvider(adapter)})
	res, errs := coord.Run(context.Background(), i.Request{}, 1, 1)
	out, errsOut := collectChannels(res, errs)
	if len(out) != 0 {
		t.Errorf("expected no responses, got %v", out)
	}
	if len(errsOut) != 1 {
		t.Errorf("expected one error, got %v", errsOut)
	}
}

// Test constraint violation produces error and no response
func TestConstraintError(t *testing.T) {
	bad := m.InternetProduct{Id: "1", Provider: "p", Name: "prod", DateOffered: time.Now(), ProductInfo: info(-1, m.DSL), Pricing: pricing(1, 1)}
	adapter := &fakeAdapter{prepareResp: i.ParsedResponse{InternetProducts: []m.InternetProduct{bad}}}
	coord := NewRequestCoordinator([]*p.ProviderConfig{newProvider(adapter)})
	res, errs := coord.Run(context.Background(), i.Request{}, 1, 1)
	out, errsOut := collectChannels(res, errs)
	if len(out) != 0 {
		t.Errorf("expected no responses, got %v", out)
	}
	if len(errsOut) != 1 {
		t.Errorf("expected one error, got %v", errsOut)
	}
}

// Test multiple providers emit all products concurrently
func TestMultipleProviders(t *testing.T) {
	prod1 := m.InternetProduct{Id: "1", Provider: "p1", Name: "a", DateOffered: time.Now(), ProductInfo: info(1, m.DSL), Pricing: pricing(1, 1)}
	prod2 := m.InternetProduct{Id: "2", Provider: "p2", Name: "b", DateOffered: time.Now(), ProductInfo: info(2, m.FIBER), Pricing: pricing(2, 2)}
	ad1 := &fakeAdapter{prepareResp: i.ParsedResponse{InternetProducts: []m.InternetProduct{prod1}}}
	ad2 := &fakeAdapter{prepareResp: i.ParsedResponse{InternetProducts: []m.InternetProduct{prod2}}}
	coord := NewRequestCoordinator([]*p.ProviderConfig{newProvider(ad1), newProvider(ad2)})
	res, errs := coord.Run(context.Background(), i.Request{}, 2, 2)
	out, errsOut := collectChannels(res, errs)
	if len(errsOut) != 0 {
		t.Errorf("expected no errors, got %v", errsOut)
	}
	if len(out) != 2 {
		t.Errorf("expected two products, got %v", out)
	}
}
