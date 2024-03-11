package delivery

import (
	"github.com/promotedai/schema/generated/go/proto/delivery"
)

// DeliveryRequestValidator performs validation on delivery requests during a deliver call when performChecks is true in the client.
type DeliveryRequestValidator interface {
	// Validate checks the state of the delivery request and collects/returns any validation errors as strings.
	Validate(request *DeliveryRequest) []string
}

// DefaultDeliveryRequestValidator implements the default delivery request validation logic.
type DefaultDeliveryRequestValidator struct{}

// Validate performs validation on the DeliveryRequest and returns any validation errors.
func (v *DefaultDeliveryRequestValidator) Validate(request *DeliveryRequest) []string {
	if request == nil {
		return []string{"DeliveryRequest is nil"}
	}

	var validationErrors []string

	reqBuilder := request.Request
	if reqBuilder == nil {
		validationErrors = append(validationErrors, "Request builder must be set")
		return validationErrors
	}

	// Check the IDs.
	validationErrors = append(validationErrors, v.validateIDs(request.Request)...)

	// Insertion start should be >= 0.
	if request.RetrievalInsertionOffset < 0 {
		validationErrors = append(validationErrors, "Insertion start must be greater or equal to 0")
	}

	return validationErrors
}

func (v *DefaultDeliveryRequestValidator) validateIDs(req *delivery.Request) []string {
	var validationErrors []string

	if req.RequestId != "" {
		validationErrors = append(validationErrors, "Request.requestID should not be set")
	}

	if req.UserInfo == nil {
		validationErrors = append(validationErrors, "Request.userInfo should be set")
	} else if req.GetUserInfo().GetAnonUserId() == "" {
		validationErrors = append(validationErrors, "Request.userInfo.anonUserID should be set")
	}

	for _, ins := range req.Insertion {
		if ins.ContentId == "" {
			validationErrors = append(validationErrors, "Insertion.contentID should be set")
		}
	}

	return validationErrors
}
