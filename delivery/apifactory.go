package delivery

// APIFactory is a factory interface for creating API clients.
type APIFactory interface {
	CreateSDKDelivery() DeliveryAPI
	CreateDeliveryAPI(endpoint, apiKey string, timeoutMillis int64, maxRequestInsertions int, acceptGzip, warmup bool) DeliveryAPI
	CreateApiMetrics(endpoint, apiKey string, timeoutMillis int64) MetricsAPI
}

// DefaultAPIFactory is the default implementation of ApiFactory.
type DefaultAPIFactory struct{}

// CreateSDKDelivery creates an SDK delivery instance.
func (f *DefaultAPIFactory) CreateSDKDelivery() DeliveryAPI {
	return NewSDKDelivery()
}

// CreateDeliveryAPI creates an API delivery instance.
func (f *DefaultAPIFactory) CreateDeliveryAPI(
	endpoint,
	apiKey string,
	timeoutMillis int64,
	maxRequestInsertions int,
	acceptGzip,
	warmup bool) DeliveryAPI {
	return NewPromotedDeliveryAPI(endpoint, apiKey, timeoutMillis, maxRequestInsertions, acceptGzip, warmup)
}

// CreateMetricsAPI creates an API metrics instance.
func (f *DefaultAPIFactory) CreateMetricsAPI(endpoint, apiKey string, timeoutMillis int64) MetricsAPI {
	return NewPromotedMetricsAPI(endpoint, apiKey, timeoutMillis)
}
