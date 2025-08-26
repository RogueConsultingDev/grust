package betteriter

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTake_YieldsTheFirstNElements(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}

	output, err := New(values).Take(3).Collect()

	require.NoError(t, err)
	assert.Equal(t, []*int{ptr(1), ptr(2), ptr(3)}, output)
}

func TestTake_StopsIfUnderlyingIteratorHasNoMoreElements(t *testing.T) {
	values := []int{1, 2}

	output, err := New(values).Take(3).Collect()

	require.NoError(t, err)
	assert.Equal(t, []*int{ptr(1), ptr(2)}, output)

	values = []int{}

	output, err = New(values).Take(3).Collect()

	require.NoError(t, err)
	assert.Equal(t, []*int{}, output)
}

func TestTake_StopsAfterN(t *testing.T) {
	iter := &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			if !yield(1, nil) {
				return
			}
			require.Fail(t, "Should not reach this point")
		},
	}

	_, err := iter.Take(1).Collect()
	assert.NoError(t, err)
}

func TestTake_PropagatesError(t *testing.T) {
	iter := &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			if !yield(1, nil) {
				return
			}
			if !yield(42, errors.New("some error")) {
				return
			}
			require.Fail(t, "Should not reach this point")
		},
	}

	output, err := iter.Take(5).Collect()
	assert.Empty(t, output)
	assert.ErrorContains(t, err, "some error")
}
