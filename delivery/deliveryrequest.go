package delivery

import (
	"log"

	"github.com/promotedai/schema/generated/go/proto/delivery"
	"github.com/promotedai/schema/generated/go/proto/event"
	"google.golang.org/protobuf/proto"
)

// NoMaxRequestInsertions means we don't trim the number of request insertions when cloning delivery requests.
const NoMaxRequestInsertions = 0

// DeliveryRequest represents the input into delivery.
type DeliveryRequest struct {
	// Request is the underlying request for ranked content.
	Request *delivery.Request

	// OnlyLog indicates whether only logs should be sent to Metrics API.
	OnlyLog bool

	// RetrievalInsertionOffset is the start index in the request insertions in the list of ALL insertions.
	RetrievalInsertionOffset int

	// Experiment is the experiment that the user is in, may be nil, which means apply the treatment.
	Experiment *event.CohortMembership

	// validator is the request validator.
	validator DeliveryRequestValidator
}

// NewDeliveryRequest is a factory method for DeliveryRequest.
func NewDeliveryRequest(
	req *delivery.Request,
	experiment *event.CohortMembership,
	onlyLog bool,
	retrievalInsertionOffset int,
	validator DeliveryRequestValidator) *DeliveryRequest {

	if validator == nil {
		validator = &DefaultDeliveryRequestValidator{}
	}
	return &DeliveryRequest{
		Request:                  req,
		OnlyLog:                  onlyLog,
		RetrievalInsertionOffset: retrievalInsertionOffset,
		Experiment:               experiment,
		validator:                validator,
	}
}

// Clone creates a copy of the DeliveryRequest, optionally trimming to maximum request insertions.
func (d *DeliveryRequest) Clone(maxRequestInsertions int) *DeliveryRequest {
	copiedRequest := proto.Clone(d.Request).(*delivery.Request)

	if maxRequestInsertions != NoMaxRequestInsertions && len(copiedRequest.Insertion) > maxRequestInsertions {
		log.Printf("Too many request insertions, truncating at %d\n", maxRequestInsertions)
		copiedRequest.Insertion = copiedRequest.Insertion[:maxRequestInsertions]
	}

	return &DeliveryRequest{
		Request:                  copiedRequest,
		OnlyLog:                  d.OnlyLog,
		RetrievalInsertionOffset: d.RetrievalInsertionOffset,
		Experiment:               d.Experiment,
		validator:                d.validator,
	}
}

// Validate checks the state of the DeliveryRequest and returns any validation errors.
func (d *DeliveryRequest) Validate() []string {
	if d == nil {
		return []string{"DeliveryRequest is nil"}
	}
	if d.validator != nil {
		return d.validator.Validate(d)
	}
	return nil
}
