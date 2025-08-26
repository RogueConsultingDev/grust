package it

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilter_FiltersElements(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	iter := New(values)

	filter := func(i int) bool {
		return i%2 == 1
	}

	output, err := Filter(iter, filter).Collect()

	require.NoError(t, err)

	assert.Equal(t, []int{1, 3, 5}, output)
}

func TestFilter_IsLazy(t *testing.T) {
	values := []int{1, 2, 3}
	iter := New(values)

	filter := func(i int) bool {
		assert.LessOrEqualf(t, i, 2, "filter was called with unexpected value: %d", i)

		return true
	}

	for v := range Filter(iter, filter).it {
		if v == 2 {
			break
		}
	}
}
