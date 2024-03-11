package delivery

import (
	"github.com/promotedai/schema/generated/go/proto/delivery"
	"github.com/stretchr/testify/mock"
)

// MockDelivery is a mock implementation of the Delivery interface.
type MockDelivery struct {
	mock.Mock
}

// RunDelivery is a mocked method for the Delivery interface.
func (m *MockDelivery) RunDelivery(deliveryRequest *DeliveryRequest) (*delivery.Response, error) {
	args := m.Called(deliveryRequest)
	return args.Get(0).(*delivery.Response), args.Error(1)
}
