package delivery

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/promotedai/schema/generated/go/proto/delivery"
)

const deliveryEndpointSuffix = "/deliver"
const healthEndpointSuffix = "/healthz"

// DeliveryAPI runs the main delivery workflow.
type DeliveryAPI interface {
	RunDelivery(deliveryRequest *DeliveryRequest) (*delivery.Response, error)
}

// PromotedDeliveryAPI is the API client for Promoted.ai's Delivery API.
type PromotedDeliveryAPI struct {
	// deliveryHTTPEndpoint is the Delivery API endpoint.
	deliveryHTTPEndpoint string

	// healthHTTPEndpoint is the API endpoint for healthchecks, also used for warmup.
	healthHTTPEndpoint string

	// apiKey required for access to Delivery API.
	apiKey string

	// httpClient for the remote call.
	httpClient *http.Client

	// timeoutDuration is the timeout set on the HTTP client as well as on the entire delivery process including SDK time.
	timeoutDuration time.Duration

	// maxRequestInsertions is the maximum number of request insertions passed to the delivery API.
	maxRequestInsertions int

	// acceptGzip indicates whether or not to try gzip processing on the request handling.
	acceptGzip bool
}

// NewPromotedDeliveryAPI instantiates a new Delivery API client.
func NewPromotedDeliveryAPI(
	endpoint,
	apiKey string,
	timeoutMillis int64,
	maxRequestInsertions int,
	acceptGzip,
	warmup bool) *PromotedDeliveryAPI {
	timeout := time.Duration(timeoutMillis) * time.Millisecond
	httpClient := &http.Client{
		Timeout: timeout,
	}

	uri, err := url.Parse(endpoint)
	if err != nil {
		log.Panic("invalid delivery endpoint")
	}
	scheme := uri.Scheme
	authority := uri.Host

	api := &PromotedDeliveryAPI{
		deliveryHTTPEndpoint: scheme + "://" + authority + deliveryEndpointSuffix,
		healthHTTPEndpoint:   scheme + "://" + authority + healthEndpointSuffix,
		apiKey:               apiKey,
		httpClient:           httpClient,
		timeoutDuration:      time.Duration(timeoutMillis) * time.Millisecond,
		maxRequestInsertions: maxRequestInsertions,
		acceptGzip:           acceptGzip,
	}

	if warmup {
		api.runWarmup()
	}

	return api
}

// RunDelivery performs delivery.
func (d *PromotedDeliveryAPI) RunDelivery(deliveryRequest *DeliveryRequest) (*delivery.Response, error) {
	var resp *delivery.Response

	ctx, cancel := context.WithTimeout(context.Background(), d.timeoutDuration)
	defer cancel()

	requestBody, err := json.Marshal(deliveryRequest.Clone(d.maxRequestInsertions).Request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling delivery request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.deliveryHTTPEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", d.apiKey)
	if d.acceptGzip {
		req.Header.Set("Accept-Encoding", "gzip")
	}

	respHTTP, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer respHTTP.Body.Close()

	if respHTTP.StatusCode < 200 || respHTTP.StatusCode >= 300 {
		return nil, fmt.Errorf("failure calling Delivery API; statusCode=%d", respHTTP.StatusCode)
	}

	if d.acceptGzip && respHTTP.Header.Get("Content-Encoding") == "gzip" {
		resp, err = d.processCompressedResponse(respHTTP.Body)
	} else {
		resp, err = d.processUncompressedResponse(respHTTP.Body)
	}
	if err != nil {
		return nil, err
	}

	if resp.RequestId == "" {
		return nil, fmt.Errorf("delivery response should contain a requestId")
	}

	return resp, nil
}

func (d *PromotedDeliveryAPI) processUncompressedResponse(body io.Reader) (*delivery.Response, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var resp delivery.Response
	err = json.Unmarshal(buf.Bytes(), &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON response: %v", err)
	}
	return &resp, nil
}

func (d *PromotedDeliveryAPI) processCompressedResponse(body io.Reader) (*delivery.Response, error) {
	gzipReader, err := gzip.NewReader(body)
	if err != nil {
		return nil, fmt.Errorf("error creating gzip reader: %v", err)
	}
	defer gzipReader.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gzipReader)
	if err != nil {
		return nil, fmt.Errorf("error reading compressed response body: %v", err)
	}

	var resp delivery.Response
	err = json.Unmarshal(buf.Bytes(), &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON response: %v", err)
	}
	return &resp, nil
}

// runWarmup performs a warmup by making GET requests to the healthzEndpoint.
func (d *PromotedDeliveryAPI) runWarmup() {
	for i := 0; i < 20; i++ {
		req, err := http.NewRequest("GET", d.healthHTTPEndpoint, nil)
		if err != nil {
			log.Print("error during warmup")
			continue
		}
		req.Header.Set("x-api-key", d.apiKey)

		_, err = d.httpClient.Do(req)
		if err != nil {
			log.Print("error during warmup")
			continue
		}
	}
}
