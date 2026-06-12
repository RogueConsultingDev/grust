//go:build go1.27

package it

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUniqueBy_FiltersOutRepeatedValues(t *testing.T) {
	values := []NonCmpT{
		{v: []int{0}},
		{v: []int{0, 1}},
		{v: []int{1, 0}},
		{v: []int{0}},
		{v: []int{1}},
		{v: []int{1, 0}},
	}
	iter := New(values)

	keyer := func(v NonCmpT) string {
		var elems []string
		for _, i := range v.v {
			elems = append(elems, strconv.Itoa(i))
		}
		return strings.Join(elems, ";")
	}

	output, err := iter.UniqueBy(keyer).Collect()
	require.NoError(t, err)
	expected := []NonCmpT{
		{v: []int{0}},
		{v: []int{0, 1}},
		{v: []int{1, 0}},
		{v: []int{1}},
	}
	assert.Equal(t, expected, output)
}

func TestUniqueBy_IsLazy(t *testing.T) {
	values := []int{1, 1, 2, 3}
	iter := New[int](values)

	keyer := func(i int) int {
		assert.LessOrEqualf(t, i, 2, "Mapper was called with unexpected value: %d", i)

		return i
	}

	for v := range iter.UniqueBy(keyer).it {
		if v == 2 {
			break
		}
	}
}

func TestUniqueBy_StopsOnError(t *testing.T) {
	values := []int{1, 1, 2, 3}
	iter := New[int](values)

	mapper := func(i int) (int, error) {
		// We will error on value 2, so mapper should never be called with value 3
		assert.LessOrEqualf(t, i, 2, "Mapper was called with unexpected value: %d", i)

		if i == 2 {
			return 0, errors.New("Invalid value")
		}

		return i, nil
	}

	output, err := iter.UniqueBy(func(i int) int { return i }).Map(mapper).Collect()
	assert.Empty(t, output)
	assert.ErrorContains(t, err, "Invalid value")
}

func TestUniqueBy_PropagatesError(t *testing.T) {
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

	output, err := iter.UniqueBy(func(i int) int { return i }).Collect()
	assert.Empty(t, output)
	assert.ErrorContains(t, err, "some error")
}
