package delivery

import "github.com/promotedai/schema/generated/go/proto/event"

// ApplyTreatmentChecker provides a method to determine whether treatment should be applied.
type ApplyTreatmentChecker interface {
	ShouldApplyTreatment(cohortMembership *event.CohortMembership) bool
}
