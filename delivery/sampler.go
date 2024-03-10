package delivery

import (
	"math/rand"
	"time"
)

type Sampler interface {
	SampleRandom(threshold float32) bool
}

// DefaultSampler is a basic implementation of a random sampler.
type DefaultSampler struct {
	rand *rand.Rand
}

// NewDefaultSampler creates a new instance of SamplerImpl.
func NewDefaultSampler() *DefaultSampler {
	return &DefaultSampler{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// SampleRandom samples random values based on the given threshold.
func (s *DefaultSampler) SampleRandom(threshold float32) bool {
	if threshold >= 1 {
		return true
	}
	if threshold <= 0 {
		return false
	}
	return s.rand.Float32() < threshold
}
