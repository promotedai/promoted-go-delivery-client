package delivery

import (
	"testing"

	"github.com/promotedai/schema/generated/go/proto/common"
	"github.com/promotedai/schema/generated/go/proto/delivery"
	"github.com/promotedai/schema/generated/go/proto/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSendShadowTrafficForOnlyLogSampledIn(t *testing.T) {
	mockSdkDelivery := new(MockDelivery)
	mockApiDelivery := new(MockDelivery)
	mockMetrics := new(MockMetrics)

	apiFactory := &TestApiFactory{
		sdkDelivery: mockSdkDelivery,
		deliveryAPI: mockApiDelivery,
		metricsAPI:  mockMetrics,
	}

	client := createDefaultClient(apiFactory, true, 0.5)

	req := &delivery.Request{
		Insertion: CreateTestRequestInsertions(10),
		ClientInfo: &common.ClientInfo{
			TrafficType: common.ClientInfo_PRODUCTION,
		},
	}
	dreq := NewDeliveryRequest(req, nil, true, 0, nil)

	mockSdkDelivery.On("RunDelivery", mock.Anything).Return(&delivery.Response{
		Insertion: req.Insertion,
	}, nil)
	mockApiDelivery.On("RunDelivery", mock.Anything).Return(&delivery.Response{
		Insertion: req.Insertion,
	}, nil)
	mockMetrics.On("RunMetricsLogging", mock.Anything).Return(nil)

	client.Deliver(dreq)

	mockSdkDelivery.AssertCalled(t, "RunDelivery", dreq)
	mockApiDelivery.AssertCalled(t, "RunDelivery", mock.AnythingOfType("*delivery.DeliveryRequest"))
	verifyShadowTrafficRequest(t, mockApiDelivery, common.ClientInfo_SHADOW)
}

func TestSendShadowTrafficForUserInControl(t *testing.T) {
	mockSdkDelivery := new(MockDelivery)
	mockApiDelivery := new(MockDelivery)
	mockMetrics := new(MockMetrics)

	apiFactory := &TestApiFactory{
		sdkDelivery: mockSdkDelivery,
		deliveryAPI: mockApiDelivery,
		metricsAPI:  mockMetrics,
	}

	client := createDefaultClient(apiFactory, true, 0.5)
	cm := &event.CohortMembership{
		Arm:      event.CohortArm_CONTROL,
		CohortId: "testing",
	}

	req := &delivery.Request{
		Insertion: CreateTestRequestInsertions(10),
		ClientInfo: &common.ClientInfo{
			TrafficType: common.ClientInfo_PRODUCTION,
		},
	}
	dreq := NewDeliveryRequest(req, cm, false, 0, nil)

	mockSdkDelivery.On("RunDelivery", mock.Anything).Return(&delivery.Response{
		Insertion: req.Insertion,
	}, nil)
	mockApiDelivery.On("RunDelivery", mock.Anything).Return(&delivery.Response{
		Insertion: req.Insertion,
	}, nil)
	mockMetrics.On("RunMetricsLogging", mock.Anything).Return(nil)

	client.Deliver(dreq)

	mockSdkDelivery.AssertCalled(t, "RunDelivery", dreq)
	mockApiDelivery.AssertCalled(t, "RunDelivery", mock.AnythingOfType("*delivery.DeliveryRequest"))
	verifyShadowTrafficRequest(t, mockApiDelivery, common.ClientInfo_SHADOW)
}

func TestDontSendShadowTrafficForOnlyLogSampledOut(t *testing.T) {
	mockSdkDelivery := new(MockDelivery)
	mockApiDelivery := new(MockDelivery)
	mockMetrics := new(MockMetrics)

	apiFactory := &TestApiFactory{
		sdkDelivery: mockSdkDelivery,
		deliveryAPI: mockApiDelivery,
		metricsAPI:  mockMetrics,
	}

	client := createDefaultClient(apiFactory, false, 0.5)

	req := &delivery.Request{
		Insertion: CreateTestRequestInsertions(10),
		ClientInfo: &common.ClientInfo{
			TrafficType: common.ClientInfo_PRODUCTION,
		},
	}
	dreq := NewDeliveryRequest(req, nil, true, 0, nil)

	mockSdkDelivery.On("RunDelivery", mock.Anything).Return(&delivery.Response{
		Insertion: req.Insertion,
	}, nil)
	mockMetrics.On("RunMetricsLogging", mock.Anything).Return(nil)

	client.Deliver(dreq)

	mockSdkDelivery.AssertCalled(t, "RunDelivery", dreq)
	mockApiDelivery.AssertNotCalled(t, "RunDelivery")
}

func TestDontSendShadowTrafficForUserInTreatment(t *testing.T) {
	mockSdkDelivery := new(MockDelivery)
	mockApiDelivery := new(MockDelivery)
	mockMetrics := new(MockMetrics)

	apiFactory := &TestApiFactory{
		sdkDelivery: mockSdkDelivery,
		deliveryAPI: mockApiDelivery,
		metricsAPI:  mockMetrics,
	}

	client := createDefaultClient(apiFactory, false, 0.5)
	cm := &event.CohortMembership{
		Arm:      event.CohortArm_TREATMENT,
		CohortId: "testing",
	}

	req := &delivery.Request{
		Insertion: CreateTestRequestInsertions(10),
		ClientInfo: &common.ClientInfo{
			TrafficType: common.ClientInfo_PRODUCTION,
		},
	}
	dreq := NewDeliveryRequest(req, cm, false, 0, nil)

	mockApiDelivery.On("RunDelivery", mock.Anything).Return(&delivery.Response{
		Insertion: req.Insertion,
	}, nil)
	mockMetrics.On("RunMetricsLogging", mock.Anything).Return(nil)

	client.Deliver(dreq)

	mockApiDelivery.AssertCalled(t, "RunDelivery", dreq)
}

func TestDontSendShadowTrafficForOnlyLogWhenTurnedOff(t *testing.T) {
	mockSdkDelivery := new(MockDelivery)
	mockApiDelivery := new(MockDelivery)
	mockMetrics := new(MockMetrics)

	apiFactory := &TestApiFactory{
		sdkDelivery: mockSdkDelivery,
		deliveryAPI: mockApiDelivery,
		metricsAPI:  mockMetrics,
	}

	client := createDefaultClient(apiFactory, true, 0)

	req := &delivery.Request{
		Insertion: CreateTestRequestInsertions(10),
		ClientInfo: &common.ClientInfo{
			TrafficType: common.ClientInfo_PRODUCTION,
		},
	}
	dreq := NewDeliveryRequest(req, nil, true, 0, nil)

	mockSdkDelivery.On("RunDelivery", mock.Anything).Return(&delivery.Response{
		Insertion: req.Insertion,
	}, nil)
	mockMetrics.On("RunMetricsLogging", mock.Anything).Return(nil)

	client.Deliver(dreq)

	mockSdkDelivery.AssertCalled(t, "RunDelivery", dreq)
	mockApiDelivery.AssertNotCalled(t, "RunDelivery")
}

func verifyShadowTrafficRequest(t *testing.T, mockApiDelivery *MockDelivery, expectedTrafficType common.ClientInfo_TrafficType) {
	captured := mockApiDelivery.Calls[0].Arguments[0].(*DeliveryRequest)
	assert.Equal(t, expectedTrafficType, captured.Request.ClientInfo.TrafficType)
}

func createDefaultClient(apiFactory APIFactory, samplesIn bool, shadowTrafficRate float32) *PromotedDeliveryClient {
	builder := NewPromotedDeliveryClientBuilder()
	client, _ := builder.
		WithBlockingShadowTraffic(true).
		WithSampler(&FakeSampler{samplesIn: samplesIn}).
		WithAPIFactory(apiFactory).
		WithShadowTrafficDeliveryRate(shadowTrafficRate).
		Build()
	return client
}
