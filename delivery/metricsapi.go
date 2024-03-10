package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/promotedai/schema/generated/go/proto/event"
)

type MetricsAPI interface {
	RunMetricsLogging(logRequest *event.LogRequest) error
}

// PromotedMetricsAPI is an API client for Promoted.ai's Metrics API.
type PromotedMetricsAPI struct {
	// Endpoint is the metrics API endpoint.
	Endpoint string

	// APIKey needed to access the metrics endpoint.
	APIKey string

	// HTTPClient used for making RPCs.
	HTTPClient *http.Client

	// TimeoutDuration is used for the http client as well as the overall metrics processing.
	TimeoutDuration time.Duration
}

// NewPromotedMetricsAPI instantiates a new Metrics API client.
func NewPromotedMetricsAPI(endpoint, apiKey string, timeoutMillis int64) *PromotedMetricsAPI {
	timeout := time.Duration(timeoutMillis) * time.Millisecond
	httpClient := &http.Client{
		Timeout: timeout,
	}

	return &PromotedMetricsAPI{
		Endpoint:        endpoint,
		APIKey:          apiKey,
		HTTPClient:      httpClient,
		TimeoutDuration: time.Duration(timeoutMillis) * time.Millisecond,
	}
}

// RunMetricsLogging performs metrics logging.
func (m *PromotedMetricsAPI) RunMetricsLogging(logRequest *event.LogRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.TimeoutDuration)
	defer cancel()

	requestBody, err := json.Marshal(logRequest)
	if err != nil {
		return fmt.Errorf("error marshaling log request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.Endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", m.APIKey)

	resp, err := m.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failure calling Metrics API; statusCode=%d", resp.StatusCode)
	}

	return nil
}
