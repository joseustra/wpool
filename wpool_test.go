package wpool_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ustrajunior/wpool"
)

func TestJob(t *testing.T) {
	wp := wpool.New(2, 100)

	var first int
	var last int

	fn1 := func() error {
		first = 10
		return nil
	}

	fn2 := func() error {
		last = 100
		return nil
	}

	wp.Add(wpool.Job{Task: fn1, Retry: 0})
	wp.Add(wpool.Job{Task: fn2, Retry: 0})

	wp.Wait()

	assert.Equal(t, 10, first)
	assert.Equal(t, 100, last)
}

func TestJobSomeVar(t *testing.T) {
	wp := wpool.New(2, 100)

	var counter int
	var mu sync.Mutex

	fn1 := func() error {
		mu.Lock()
		counter++
		mu.Unlock()
		return nil
	}

	fn2 := func() error {
		mu.Lock()
		counter++
		mu.Unlock()
		return nil
	}

	wp.Add(wpool.Job{Task: fn1, Retry: 0})
	wp.Add(wpool.Job{Task: fn2, Retry: 0})

	wp.Wait()

	assert.Equal(t, 2, counter)
}

func TestRetry(t *testing.T) {
	wp := wpool.New(2, 100)

	s := &struct {
		Value int
	}{
		Value: 0,
	}

	fn := func() error {
		s.Value++
		return errors.New("fail")
	}

	wp.Add(wpool.Job{Task: fn, Retry: 5})

	wp.Wait()

	assert.Equal(t, 6, s.Value)
}
