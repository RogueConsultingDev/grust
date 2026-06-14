//go:build go1.27

package it

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFold_AppliesTheFolderFunctionOnAllValues(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}

	iter := New[int](values)

	res1, err := iter.Fold(0, func(cur int, v int) int { return cur + v })
	require.NoError(t, err)

	assert.Equal(t, 15, res1)

	iter2 := New[int](values)
	res2, err := iter2.Fold(
		"result: ",
		func(cur string, v int) string { return cur + strconv.Itoa(v) },
	)
	require.NoError(t, err)

	assert.Equal(t, "result: 12345", res2)
}

func TestFold_PropagatesError(t *testing.T) {
	iter := &Iterator[int]{
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

	output, err := iter.Fold(0, func(_ int, v int) int { return v })
	require.ErrorContains(t, err, "some error")
	// folding should have stopped at 1
	assert.Equal(t, 1, output)
}
