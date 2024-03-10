package delivery

import (
	"testing"

	"github.com/promotedai/schema/generated/go/proto/common"
	"github.com/promotedai/schema/generated/go/proto/delivery"
	"github.com/promotedai/schema/generated/go/proto/event"
	"github.com/stretchr/testify/assert"
)

func TestRequestMustBeSet(t *testing.T) {
	req := NewDeliveryRequest(nil, nil, false, 0, nil)
	errors := req.Validate()
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "Request builder must be set", errors[0])
}

func TestValidateRequestIdMustBeUnsetOnRequest(t *testing.T) {
	req := NewDeliveryRequest(
		&delivery.Request{
			RequestId: "z",
			UserInfo: &common.UserInfo{
				AnonUserId: "a",
			},
		},
		nil,
		false,
		0,
		nil,
	)
	errors := req.Validate()
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "Request.requestID should not be set", errors[0])
}

func TestValidateRetrievalInsertionOffsetMustBeNonNeg(t *testing.T) {
	req := NewDeliveryRequest(
		&delivery.Request{
			UserInfo: &common.UserInfo{
				AnonUserId: "a",
			},
			Insertion: []*delivery.Insertion{{ContentId: "z"}},
		},
		nil,
		false,
		-1,
		nil,
	)
	errors := req.Validate()
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "Insertion start must be greater or equal to 0", errors[0])
}

func TestValidateContentIdMustBeSet(t *testing.T) {
	req := NewDeliveryRequest(
		&delivery.Request{
			UserInfo: &common.UserInfo{
				AnonUserId: "a",
			},
			Insertion: []*delivery.Insertion{{ContentId: ""}},
		},
		nil,
		false,
		0,
		nil,
	)
	errors := req.Validate()
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "Insertion.contentID should be set", errors[0])
}

func TestValidateWithValidInsertion(t *testing.T) {
	req := NewDeliveryRequest(
		&delivery.Request{
			UserInfo: &common.UserInfo{
				AnonUserId: "a",
			},
			Insertion: []*delivery.Insertion{{ContentId: "z"}},
		},
		nil,
		false,
		0,
		nil,
	)
	errors := req.Validate()
	assert.Equal(t, 0, len(errors))
}

func TestValidateExperimentValid(t *testing.T) {
	req := NewDeliveryRequest(
		&delivery.Request{
			UserInfo: &common.UserInfo{
				AnonUserId: "a",
			},
			Insertion: []*delivery.Insertion{{ContentId: "z"}},
		},
		&event.CohortMembership{Arm: event.CohortArm_TREATMENT, CohortId: "my cohort"},
		false,
		0,
		nil,
	)
	errors := req.Validate()
	assert.Equal(t, 0, len(errors))
}

func TestValidateUserInfoOnRequest(t *testing.T) {
	req := NewDeliveryRequest(
		&delivery.Request{
			Insertion: []*delivery.Insertion{{ContentId: "z"}},
		},
		nil,
		false,
		0,
		nil,
	)
	errors := req.Validate()
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "Request.userInfo should be set", errors[0])
}

func TestValidateAnonUserIdOnRequest(t *testing.T) {
	req := NewDeliveryRequest(
		&delivery.Request{
			UserInfo: &common.UserInfo{
				AnonUserId: "",
			},
			Insertion: []*delivery.Insertion{{ContentId: "z"}},
		},
		&event.CohortMembership{Arm: event.CohortArm_TREATMENT, CohortId: "my cohort"},
		false,
		0,
		nil,
	)
	errors := req.Validate()
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "Request.userInfo.anonUserID should be set", errors[0])
}

func TestValidateCapturesMultipleErrors(t *testing.T) {
	req := NewDeliveryRequest(
		&delivery.Request{
			RequestId: "a",
			UserInfo: &common.UserInfo{
				AnonUserId: "",
			},
			Insertion: []*delivery.Insertion{{ContentId: "z"}},
		},
		&event.CohortMembership{Arm: event.CohortArm_TREATMENT, CohortId: "my cohort"},
		false,
		0,
		nil,
	)
	errors := req.Validate()
	assert.Equal(t, 2, len(errors))
	assert.Equal(t, "Request.requestID should not be set", errors[0])
	assert.Equal(t, "Request.userInfo.anonUserID should be set", errors[1])
}
