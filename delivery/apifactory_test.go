package delivery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultAPIFactory_IsAPIFactory(t *testing.T) {
	// Make sure this builds.
	var apiFactory APIFactory
	apiFactory = &DefaultAPIFactory{}
	assert.NotNil(t, apiFactory)
}
