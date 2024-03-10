package delivery

import (
	"errors"
	"math"

	"github.com/google/uuid"
	"github.com/promotedai/schema/generated/go/proto/delivery"
)

// SDKDelivery implements SDK-side delivery.
type SDKDelivery struct{}

// NewSDKDelivery is a factory method for SDKDelivery.
func NewSDKDelivery() *SDKDelivery {
	return &SDKDelivery{}
}

// RunDelivery performs delivery.
func (sdk *SDKDelivery) RunDelivery(deliveryRequest *DeliveryRequest) (*delivery.Response, error) {
	req := deliveryRequest.Request

	// Set a request id.
	req.RequestId = uuid.New().String()

	var paging *delivery.Paging
	if req.Paging == nil {
		paging = NewPaging(int32(len(req.Insertion)), 0)
	} else {
		paging = req.Paging
	}

	// Adjust offset and size.
	offset := uint64(math.Max(0, float64(paging.GetOffset())))
	index := offset - uint64(deliveryRequest.RetrievalInsertionOffset)
	if offset < uint64(deliveryRequest.RetrievalInsertionOffset) {
		return nil, errors.New("offset should be >= insertion start (specifically, the global position)")
	}

	size := int(paging.Size)
	if size <= 0 {
		size = len(req.Insertion)
	}

	finalInsertionSize := int(math.Min(float64(size), float64(uint64(len(req.Insertion))-index)))
	resp := &delivery.Response{RequestId: req.RequestId, Insertion: make([]*delivery.Insertion, 0, finalInsertionSize)}
	for i := 0; i < finalInsertionSize; i++ {
		reqIns := req.Insertion[index]
		resp.Insertion = append(resp.Insertion, newResponseInsertion(reqIns, offset))
		index++
		offset++
	}
	return resp, nil
}

// newResponseInsertion prepares the response insertion.
func newResponseInsertion(reqIns *delivery.Insertion, offset uint64) *delivery.Insertion {
	insID := reqIns.InsertionId
	if len(insID) == 0 {
		insID = uuid.NewString()
	}
	respIns := &delivery.Insertion{
		ContentId:   reqIns.ContentId,
		InsertionId: insID,
		Position:    &offset,
	}
	return respIns
}
