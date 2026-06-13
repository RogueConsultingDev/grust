package it

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilter_FiltersElements(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	iter := New(values)

	predicate := func(i int) bool {
		return i%2 == 1
	}

	output, err := iter.Filter(predicate).Collect()

	require.NoError(t, err)

	assert.Equal(t, []int{1, 3, 5}, output)
}

func TestFilter_IsLazy(t *testing.T) {
	values := []int{1, 2, 3}
	iter := New(values)

	predicate := func(i int) bool {
		assert.LessOrEqualf(t, i, 2, "filter was called with unexpected value: %d", i)

		return true
	}

	for v := range iter.Filter(predicate).it {
		if v == 2 {
			break
		}
	}
}

func TestFilter_PropagatesError(t *testing.T) {
	iter := &Iterator[int]{
		it: func(yield func(int, error) bool) {
			if !yield(1, nil) {
				return
			}

			if !yield(0, errors.New("some error")) {
				return
			}

			require.Fail(t, "Should not reach this point")
		},
	}

	predicate := func(i int) bool {
		// We should only be called with value = 1
		assert.Equal(t, 1, i, "Filter was called with unexpected value: %d", i)

		return true
	}

	output, err := iter.Filter(predicate).Collect()
	assert.Empty(t, output)
	assert.ErrorContains(t, err, "some error")
}
