package delivery

import (
	"strconv"

	"github.com/promotedai/schema/generated/go/proto/delivery"
)

type TestApiFactory struct {
	metricsAPI  MetricsAPI
	deliveryAPI DeliveryAPI
	sdkDelivery DeliveryAPI
}

// CreateSDKDelivery creates an SDK delivery instance.
func (f *TestApiFactory) CreateSDKDelivery() DeliveryAPI {
	return f.sdkDelivery
}

// CreateDeliveryAPI creates an API delivery instance.
func (f *TestApiFactory) CreateDeliveryAPI(endpoint, apiKey string, timeoutMillis int64, maxRequestInsertions int, acceptGzip, warmup bool) DeliveryAPI {
	return f.deliveryAPI
}

// CreateApiMetrics creates an API metrics instance.
func (f *TestApiFactory) CreateApiMetrics(endpoint, apiKey string, timeoutMillis int64) MetricsAPI {
	return f.metricsAPI
}

type FakeSampler struct {
	samplesIn bool
}

func (s *FakeSampler) SampleRandom(threshold float32) bool {
	return s.samplesIn
}

func CreateTestRequestInsertions(num int) []*delivery.Insertion {
	var res []*delivery.Insertion
	for i := 0; i < num; i++ {
		res = append(res, &delivery.Insertion{ContentId: strconv.Itoa(i)})
	}
	return res
}

func CreateTestResponseInsertions(num, offset int) []*delivery.Insertion {
	var res []*delivery.Insertion
	for i := 0; i < num; i++ {
		res = append(res, &delivery.Insertion{
			ContentId:   strconv.Itoa(i),
			Position:    uint64Pointer(i + offset),
			InsertionId: "id" + strconv.Itoa(i),
		})
	}
	return res
}

func uint64Pointer(val int) *uint64 {
	newVal := uint64(val)
	return &newVal
}
