package delivery

import (
	"testing"

	"github.com/promotedai/schema/generated/go/proto/event"
	"github.com/stretchr/testify/assert"
)

func TestTwoArmExperiment_CreateSuccess(t *testing.T) {
	exp, err := NewTwoArmExperiment("HOLD_OUT", 10, 50, 10, 50)
	assert.NoError(t, err)
	assert.Equal(t, "HOLD_OUT", exp.CohortID)
	assert.Equal(t, 50, exp.NumControlBuckets)
	assert.Equal(t, 50, exp.NumTreatmentBuckets)
	assert.Equal(t, 10, exp.NumActiveTreatmentBuckets)
	assert.Equal(t, 10, exp.NumActiveControlBuckets)
}

func TestTwoArmExperiment_CreateInvalidCohortID(t *testing.T) {
	_, err := NewTwoArmExperiment("", 10, 50, 10, 50)
	assert.NotNil(t, err)
	assert.Equal(t, "cohort ID must be non-empty", err.Error())

	_, err = NewTwoArmExperiment(" ", 10, 50, 10, 50)
	assert.NotNil(t, err)
	assert.Equal(t, "cohort ID must be non-empty", err.Error())
}

func TestTwoArmExperiment_CreateInvalidBucketActiveCounts(t *testing.T) {
	_, err := NewTwoArmExperiment("a", -1, 50, 10, 50)
	assert.NotNil(t, err)
	assert.Equal(t, "active control buckets must be between 0 and the total number of control buckets", err.Error())

	_, err = NewTwoArmExperiment("a", 51, 50, 10, 50)
	assert.NotNil(t, err)
	assert.Equal(t, "active control buckets must be between 0 and the total number of control buckets", err.Error())

	_, err = NewTwoArmExperiment("a", 10, 50, -1, 50)
	assert.NotNil(t, err)
	assert.Equal(t, "active treatment buckets must be between 0 and the total number of treatment buckets", err.Error())

	_, err = NewTwoArmExperiment("a", 10, 50, 51, 50)
	assert.NotNil(t, err)
	assert.Equal(t, "active treatment buckets must be between 0 and the total number of treatment buckets", err.Error())
}

func TestTwoArmExperiment_CreateInvalidBucketCounts(t *testing.T) {
	_, err := NewTwoArmExperiment("a", 0, -1, 10, 50)
	assert.NotNil(t, err)
	assert.Equal(t, "control buckets must be positive", err.Error())

	_, err = NewTwoArmExperiment("a", 10, 50, 0, -1)
	assert.NotNil(t, err)
	assert.Equal(t, "treatment buckets must be positive", err.Error())
}

func TestTwoArmExperiment_CreateTwoArmExperiment1PercentSuccess(t *testing.T) {
	exp, err := Create5050TwoArmExperimentConfig("HOLD_OUT", 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "HOLD_OUT", exp.CohortID)
	assert.Equal(t, 50, exp.NumControlBuckets)
	assert.Equal(t, 50, exp.NumTreatmentBuckets)
	assert.Equal(t, 1, exp.NumActiveTreatmentBuckets)
	assert.Equal(t, 1, exp.NumActiveControlBuckets)
}

func TestTwoArmExperiment_CreateTwoArmExperiment10And5PercentSuccess(t *testing.T) {
	exp, err := Create5050TwoArmExperimentConfig("HOLD_OUT", 10, 5)
	assert.NoError(t, err)
	assert.Equal(t, "HOLD_OUT", exp.CohortID)
	assert.Equal(t, 50, exp.NumControlBuckets)
	assert.Equal(t, 50, exp.NumTreatmentBuckets)
	assert.Equal(t, 5, exp.NumActiveTreatmentBuckets)
	assert.Equal(t, 10, exp.NumActiveControlBuckets)
}

func TestTwoArmExperiment_UserInControl(t *testing.T) {
	exp, err := Create5050TwoArmExperimentConfig("HOLD_OUT", 50, 50)
	assert.NoError(t, err)
	mem := exp.CheckMembership("user2")
	assert.NotNil(t, mem)
	assert.Equal(t, "HOLD_OUT", mem.CohortId)
	assert.Equal(t, event.CohortArm_CONTROL, mem.GetArm())
}

func TestTwoArmExperiment_UserNotActive(t *testing.T) {
	exp, err := Create5050TwoArmExperimentConfig("HOLD_OUT", 1, 1)
	assert.NoError(t, err)
	mem := exp.CheckMembership("user5")
	assert.Nil(t, mem)
}

func TestTwoArmExperiment_UserInTreatment(t *testing.T) {
	exp, err := Create5050TwoArmExperimentConfig("HOLD_OUT", 50, 50)
	assert.NoError(t, err)
	mem := exp.CheckMembership("user4")
	assert.NotNil(t, mem)
	assert.Equal(t, "HOLD_OUT", mem.CohortId)
	assert.Equal(t, event.CohortArm_TREATMENT, mem.GetArm())
}

func TestTwoArmExperiment_Create5050InvalidPercents(t *testing.T) {
	_, err := Create5050TwoArmExperimentConfig("HOLD_OUT", -1, 50)
	assert.NotNil(t, err)
	assert.Equal(t, "control percent must be in the range [0, 50]", err.Error())

	_, err = Create5050TwoArmExperimentConfig("HOLD_OUT", 51, 50)
	assert.NotNil(t, err)
	assert.Equal(t, "control percent must be in the range [0, 50]", err.Error())

	_, err = Create5050TwoArmExperimentConfig("HOLD_OUT", 50, -1)
	assert.NotNil(t, err)
	assert.Equal(t, "treatment percent must be in the range [0, 50]", err.Error())

	_, err = Create5050TwoArmExperimentConfig("HOLD_OUT", 50, 51)
	assert.NotNil(t, err)
	assert.Equal(t, "treatment percent must be in the range [0, 50]", err.Error())
}
