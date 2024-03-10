package delivery

import (
	"testing"

	"github.com/promotedai/schema/generated/go/proto/delivery"
	"github.com/stretchr/testify/assert"
)

func TestSdkDelivery_InvalidPagingOffsetAndRetrievalInsertionOffset(t *testing.T) {
	req := &delivery.Request{
		Paging:    NewPaging(5, 10),
		Insertion: CreateTestRequestInsertions(10),
	}
	dreq := NewDeliveryRequest(req, nil, false, 100, nil)
	resp, err := NewSDKDelivery().RunDelivery(dreq)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "offset should be >= insertion start")
}

func TestSdkDelivery_ValidPagingOffsetAndRetrievalInsertionOffset(t *testing.T) {
	req := &delivery.Request{
		Paging:    NewPaging(5, 10),
		Insertion: CreateTestRequestInsertions(10),
	}
	dreq := NewDeliveryRequest(req, nil, false, 5, nil)
	resp, err := NewSDKDelivery().RunDelivery(dreq)
	assert.NotNil(t, resp)
	assert.NoError(t, err)
}

func TestSdkDelivery_NoPagingReturnsAll(t *testing.T) {
	insertions := CreateTestRequestInsertions(10)
	req := &delivery.Request{
		Insertion: insertions,
	}
	dreq := NewDeliveryRequest(req, nil, false, 0, nil)

	resp, err := NewSDKDelivery().RunDelivery(dreq)
	assert.NoError(t, err)
	assert.True(t, len(req.RequestId) > 0)
	assert.True(t, len(resp.RequestId) > 0)
	assertAllResponseInsertions(t, resp)
}

func TestSdkDelivery_RetrievalInsertionOffsetSetToOffset(t *testing.T) {
	retrievalInsertionOffset := 5
	insertions := CreateTestRequestInsertions(3)
	req := &delivery.Request{
		Insertion: insertions,
		Paging:    NewPaging(2, 5),
	}
	dreq := NewDeliveryRequest(req, nil, false, retrievalInsertionOffset, nil)

	resp, err := NewSDKDelivery().RunDelivery(dreq)
	assert.NoError(t, err)
	assert.True(t, len(req.RequestId) > 0)
	assert.True(t, len(resp.RequestId) > 0)
	assert.Equal(t, 2, len(resp.Insertion))
}

func TestSdkDelivery_PagingZeroOffset(t *testing.T) {
	insertions := CreateTestRequestInsertions(10)
	req := &delivery.Request{
		Insertion: insertions,
		Paging:    NewPaging(5, 0),
	}
	dreq := NewDeliveryRequest(req, nil, false, 0, nil)

	resp, err := NewSDKDelivery().RunDelivery(dreq)
	assert.NoError(t, err)
	assert.True(t, len(req.RequestId) > 0)
	assert.True(t, len(resp.RequestId) > 0)
	assert.Equal(t, 5, len(resp.Insertion))
	for i := 0; i < 5; i++ {
		insertion := resp.Insertion[i]
		assert.Equal(t, uint64(i), *insertion.Position)
	}
}

func TestSdkDelivery_PagingNonZeroOffset(t *testing.T) {
	insertions := CreateTestRequestInsertions(10)
	req := &delivery.Request{
		Insertion: insertions,
		Paging:    NewPaging(5, 5),
	}
	dreq := NewDeliveryRequest(req, nil, false, 0, nil)

	resp, err := NewSDKDelivery().RunDelivery(dreq)
	assert.NoError(t, err)
	assert.True(t, len(req.RequestId) > 0)
	assert.True(t, len(resp.RequestId) > 0)
	assert.Equal(t, 5, len(resp.Insertion))
	for i := 5; i < 10; i++ {
		insertion := resp.Insertion[i-5]
		assert.Equal(t, uint64(i), *insertion.Position)
	}
}

func TestSdkDelivery_PagingSizeMoreThanInsertions(t *testing.T) {
	insertions := CreateTestRequestInsertions(10)
	req := &delivery.Request{
		Insertion: insertions,
		Paging:    NewPaging(11, 0),
	}
	dreq := NewDeliveryRequest(req, nil, false, 0, nil)

	resp, err := NewSDKDelivery().RunDelivery(dreq)
	assert.NoError(t, err)
	assert.True(t, len(req.RequestId) > 0)
	assert.True(t, len(resp.RequestId) > 0)
	assertAllResponseInsertions(t, resp)
}

func assertAllResponseInsertions(t *testing.T, resp *delivery.Response) {
	assert.Equal(t, 10, len(resp.Insertion))
	for i := 0; i < 10; i++ {
		insertion := resp.Insertion[i]
		assert.Equal(t, uint64(i), *insertion.Position)
		assert.True(t, len(insertion.InsertionId) > 0)
	}
}
