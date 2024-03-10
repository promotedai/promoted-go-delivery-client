package delivery

import (
	"errors"
)

type PromotedDeliveryClientBuilder struct {
	deliveryEndpoint          string
	deliveryAPIKey            string
	deliveryTimeoutMillis     int64
	metricsEndpoint           string
	metricsAPIKey             string
	metricsTimeoutMillis      int64
	warmup                    bool
	maxRequestInsertions      int
	applyTreatmentChecker     ApplyTreatmentChecker
	sampler                   Sampler
	apiFactory                APIFactory
	shadowTrafficDeliveryRate float32
	performChecks             bool
	blockingShadowTraffic     bool
	acceptsGzip               bool
}

// NewPromotedDeliveryClientBuilder implements a builder interface for PromotedDeliveryClient.
func NewPromotedDeliveryClientBuilder() *PromotedDeliveryClientBuilder {
	return &PromotedDeliveryClientBuilder{
		deliveryTimeoutMillis: defaultDeliveryTimeoutMillis,
		metricsTimeoutMillis:  defaultMetricsTimeoutMillis,
		maxRequestInsertions:  defaultMaxRequestInsertions,
	}
}

func (b *PromotedDeliveryClientBuilder) WithDeliveryEndpoint(deliveryEndpoint string) *PromotedDeliveryClientBuilder {
	b.deliveryEndpoint = deliveryEndpoint
	return b
}

func (b *PromotedDeliveryClientBuilder) WithDeliveryAPIKey(deliveryAPIKey string) *PromotedDeliveryClientBuilder {
	b.deliveryAPIKey = deliveryAPIKey
	return b
}

func (b *PromotedDeliveryClientBuilder) WithMetricsEndpoint(metricsEndpoint string) *PromotedDeliveryClientBuilder {
	b.metricsEndpoint = metricsEndpoint
	return b
}

func (b *PromotedDeliveryClientBuilder) WithMetricsAPIKey(metricsAPIKey string) *PromotedDeliveryClientBuilder {
	b.metricsAPIKey = metricsAPIKey
	return b
}

func (b *PromotedDeliveryClientBuilder) WithDeliveryTimeoutMillis(deliveryTimeoutMillis int64) *PromotedDeliveryClientBuilder {
	b.deliveryTimeoutMillis = deliveryTimeoutMillis
	return b
}

func (b *PromotedDeliveryClientBuilder) WithMetricsTimeoutMillis(metricsTimeoutMillis int64) *PromotedDeliveryClientBuilder {
	b.metricsTimeoutMillis = metricsTimeoutMillis
	return b
}

func (b *PromotedDeliveryClientBuilder) WithMaxRequestInsertions(maxRequestInsertions int) *PromotedDeliveryClientBuilder {
	b.maxRequestInsertions = maxRequestInsertions
	return b
}

func (b *PromotedDeliveryClientBuilder) WithApplyTreatmentChecker(applyTreatmentChecker ApplyTreatmentChecker) *PromotedDeliveryClientBuilder {
	b.applyTreatmentChecker = applyTreatmentChecker
	return b
}

func (b *PromotedDeliveryClientBuilder) WithSampler(sampler Sampler) *PromotedDeliveryClientBuilder {
	b.sampler = sampler
	return b
}

func (b *PromotedDeliveryClientBuilder) WithAPIFactory(apiFactory APIFactory) *PromotedDeliveryClientBuilder {
	b.apiFactory = apiFactory
	return b
}

func (b *PromotedDeliveryClientBuilder) WithShadowTrafficDeliveryRate(shadowTrafficDeliveryRate float32) *PromotedDeliveryClientBuilder {
	b.shadowTrafficDeliveryRate = shadowTrafficDeliveryRate
	return b
}

func (b *PromotedDeliveryClientBuilder) WithPerformChecks(performChecks bool) *PromotedDeliveryClientBuilder {
	b.performChecks = performChecks
	return b
}

func (b *PromotedDeliveryClientBuilder) WithBlockingShadowTraffic(blockingShadowTraffic bool) *PromotedDeliveryClientBuilder {
	b.blockingShadowTraffic = blockingShadowTraffic
	return b
}

func (b *PromotedDeliveryClientBuilder) WithAcceptsGzip(acceptsGzip bool) *PromotedDeliveryClientBuilder {
	b.acceptsGzip = acceptsGzip
	return b
}

func (b *PromotedDeliveryClientBuilder) Build() (*PromotedDeliveryClient, error) {
	if b.deliveryTimeoutMillis <= 0 {
		b.deliveryTimeoutMillis = defaultDeliveryTimeoutMillis
	}

	if b.metricsTimeoutMillis <= 0 {
		b.metricsTimeoutMillis = defaultMetricsTimeoutMillis
	}

	if b.maxRequestInsertions <= 0 {
		b.maxRequestInsertions = defaultMaxRequestInsertions
	}

	if b.shadowTrafficDeliveryRate < 0 || b.shadowTrafficDeliveryRate > 1 {
		return nil, errors.New("shadowTrafficDeliveryRate must be between 0 and 1")
	}

	deliveryAPI := b.apiFactory.CreateDeliveryAPI(
		b.deliveryEndpoint,
		b.deliveryAPIKey,
		b.deliveryTimeoutMillis,
		b.maxRequestInsertions,
		b.acceptsGzip,
		b.warmup,
	)

	metricsAPI := b.apiFactory.CreateApiMetrics(
		b.metricsEndpoint,
		b.metricsAPIKey,
		b.metricsTimeoutMillis,
	)

	if b.sampler == nil {
		b.sampler = NewDefaultSampler()
	}

	return &PromotedDeliveryClient{
		deliveryAPI:               deliveryAPI,
		metricsAPI:                metricsAPI,
		sdkDelivery:               b.apiFactory.CreateSDKDelivery(),
		deliveryEndpoint:          b.deliveryEndpoint,
		deliveryAPIKey:            b.deliveryAPIKey,
		deliveryTimeoutMillis:     b.deliveryTimeoutMillis,
		metricsEndpoint:           b.metricsEndpoint,
		metricsAPIKey:             b.metricsAPIKey,
		metricsTimeoutMillis:      b.metricsTimeoutMillis,
		maxRequestInsertions:      b.maxRequestInsertions,
		applyTreatmentChecker:     b.applyTreatmentChecker,
		shadowTrafficDeliveryRate: b.shadowTrafficDeliveryRate,
		sampler:                   b.sampler,
		performChecks:             b.performChecks,
		blockingShadowTraffic:     b.blockingShadowTraffic,
	}, nil
}
