package delivery

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxThreshold(t *testing.T) {
	s := NewDefaultSampler()
	assert.True(t, s.SampleRandom(1))
}

func TestMinThreshold(t *testing.T) {
	s := NewDefaultSampler()
	assert.False(t, s.SampleRandom(0))
}

func TestRandomness(t *testing.T) {
	rand := rand.New(rand.NewSource(0))
	sampler := &DefaultSampler{rand: rand}
	assert.False(t, sampler.SampleRandom(0.5))
	assert.True(t, sampler.SampleRandom(0.5))
}
