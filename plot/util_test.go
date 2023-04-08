package plot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRingBuffer(t *testing.T) {
	r := newRingBuffer[int](10)

	for i := 0; i < 10; i++ {
		r.Add(i)
		assert.Equal(t, i+1, r.Len())
	}

	for i := 0; i < 10; i++ {
		assert.Equal(t, i, r.Get(i))
	}

	for i := 0; i < 10; i++ {
		r.Add(i + 10)
		assert.Equal(t, 10, r.Len())
		assert.Equal(t, i+1, r.Get(0))
	}

	for i := 0; i < 10; i++ {
		assert.Equal(t, i+10, r.Get(i))
	}
}
