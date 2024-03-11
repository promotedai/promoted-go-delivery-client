package delivery

import (
	"github.com/promotedai/schema/generated/go/proto/event"
	"github.com/stretchr/testify/mock"
)

// MockMetrics is a mock implementation of the MetricsAPI interface.
type MockMetrics struct {
	mock.Mock
}

// RunMetricsLogging mocks the RunMetricsLogging method of the MetricsAPI interface.
func (m *MockMetrics) RunMetricsLogging(logRequest *event.LogRequest) error {
	args := m.Called(logRequest)
	return args.Error(0)
}
