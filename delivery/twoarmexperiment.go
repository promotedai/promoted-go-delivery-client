package delivery

import (
	"errors"
	"hash/fnv"
	"math"
	"strings"

	"github.com/promotedai/schema/generated/go/proto/event"
)

// TwoArmExperiment represents a two arm Experiment configuration.
type TwoArmExperiment struct {
	// Name of cohort.
	CohortID string

	// Hash of cohort ID.
	CohortIDHash int

	// Number of the numControlBuckets that are active.
	NumActiveControlBuckets int

	// Number of control buckets.
	NumControlBuckets int

	// Number of the numTreatmentBuckets that are active.
	NumActiveTreatmentBuckets int

	// Number of treatment buckets.
	NumTreatmentBuckets int

	// Total number of buckets.
	NumTotalBuckets int
}

// Create5050TwoArmExperimentConfig is a factory method for a 50/50 experiment.
func Create5050TwoArmExperimentConfig(cohortID string, controlPercent, treatmentPercent int) (*TwoArmExperiment, error) {
	if controlPercent < 0 || controlPercent > 50 {
		return nil, errors.New("control percent must be in the range [0, 50]")
	}
	if treatmentPercent < 0 || treatmentPercent > 50 {
		return nil, errors.New("treatment percent must be in the range [0, 50]")
	}
	return &TwoArmExperiment{
		CohortID:                  cohortID,
		CohortIDHash:              hash(cohortID),
		NumActiveControlBuckets:   controlPercent,
		NumControlBuckets:         50,
		NumActiveTreatmentBuckets: treatmentPercent,
		NumTreatmentBuckets:       50,
		NumTotalBuckets:           100,
	}, nil
}

// NewTwoArmExperiment creates a two-arm experiment config with the given parameters.
func NewTwoArmExperiment(cohortID string, numActiveControlBuckets, numControlBuckets, numActiveTreatmentBuckets, numTreatmentBuckets int) (*TwoArmExperiment, error) {
	if strings.TrimSpace(cohortID) == "" {
		return nil, errors.New("cohort ID must be non-empty")
	}
	if numControlBuckets < 0 {
		return nil, errors.New("control buckets must be positive")
	}
	if numTreatmentBuckets < 0 {
		return nil, errors.New("treatment buckets must be positive")
	}
	if numActiveControlBuckets < 0 || numActiveControlBuckets > numControlBuckets {
		return nil, errors.New("active control buckets must be between 0 and the total number of control buckets")
	}
	if numActiveTreatmentBuckets < 0 || numActiveTreatmentBuckets > numTreatmentBuckets {
		return nil, errors.New("active treatment buckets must be between 0 and the total number of treatment buckets")
	}
	return &TwoArmExperiment{
		CohortID:                  cohortID,
		CohortIDHash:              hash(cohortID),
		NumActiveControlBuckets:   numActiveControlBuckets,
		NumControlBuckets:         numControlBuckets,
		NumActiveTreatmentBuckets: numActiveTreatmentBuckets,
		NumTreatmentBuckets:       numTreatmentBuckets,
		NumTotalBuckets:           numControlBuckets + numTreatmentBuckets,
	}, nil
}

// CheckMembership evaluates the experiment membership for a given user.
func (e *TwoArmExperiment) CheckMembership(userID string) *event.CohortMembership {
	hash := e.combineHash(hash(userID), e.CohortIDHash)
	bucket := int(math.Abs(float64(hash))) % e.NumTotalBuckets
	if bucket < e.NumActiveControlBuckets {
		return &event.CohortMembership{
			CohortId: e.CohortID,
			Arm:      event.CohortArm_CONTROL,
		}
	}
	if e.NumControlBuckets <= bucket && bucket < e.NumControlBuckets+e.NumActiveTreatmentBuckets {
		return &event.CohortMembership{
			CohortId: e.CohortID,
			Arm:      event.CohortArm_TREATMENT,
		}
	}
	return nil
}

// combineHash returns a simple combined hash of two other hashes.
func (e *TwoArmExperiment) combineHash(hash1, hash2 int) int {
	hash := 17
	hash = hash*31 + hash1
	hash = hash*31 + hash2
	return hash
}

// hash returns the hash value of a string.
func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}
