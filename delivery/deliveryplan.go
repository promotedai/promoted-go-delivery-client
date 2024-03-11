package delivery

// DeliveryPlan represents a plan object that indicates how to execute delivery.
type DeliveryPlan struct {
	ClientRequestID string
	UseAPIResponse  bool
}

// NewDeliveryPlan is a factory method for DeliveryPlan.
func NewDeliveryPlan(clientRequestID string, useAPIResponse bool) *DeliveryPlan {
	return &DeliveryPlan{
		ClientRequestID: clientRequestID,
		UseAPIResponse:  useAPIResponse,
	}
}
