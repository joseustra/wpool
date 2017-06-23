package wpool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJob(t *testing.T) {
	wp := New(2)

	var first int
	var last int

	wp.Add(func() error {
		first = 10
		return nil
	})

	wp.Add(func() error {
		last = 100
		return nil
	})

	wp.Wait()

	assert.Equal(t, 10, first)
	assert.Equal(t, 100, last)
}
