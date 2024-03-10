package delivery

import (
	"github.com/promotedai/schema/generated/go/proto/delivery"
)

// DeliveryResponse is the output from delivery.
type DeliveryResponse struct {
	// The response from Delivery.
	Response *delivery.Response

	// ClientRequestID is for tracking purposes, auto-generated if not supplied on the request.
	ClientRequestID string

	// ExecutionServer indicates if delivery happened in the SDK or via Delivery API.
	ExecutionServer delivery.ExecutionServer
}
