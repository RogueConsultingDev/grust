package it

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMap_TransformsElements(t *testing.T) {
	values := []int{1, 2, 3}
	iter := New(values)

	mapper := func(i int) (string, error) {
		return strconv.Itoa(i * i), nil
	}

	output, err := Map(iter, mapper).Collect()

	require.NoError(t, err)

	assert.Equal(t, []string{"1", "4", "9"}, output)
}

func TestMap_IsLazy(t *testing.T) {
	values := []int{1, 2, 3}
	iter := New(values)

	mapper := func(i int) (int, error) {
		assert.LessOrEqualf(t, i, 2, "Mapper was called with unexpected value: %d", i)

		return i, nil
	}

	for v := range Map(iter, mapper).it {
		if v == 2 {
			break
		}
	}
}

func TestMap_StopsOnError(t *testing.T) {
	values := []int{1, 2, 3}
	iter := New(values)

	mapper := func(i int) (int, error) {
		// We will error on value 2, so mapper should never be called with value 3
		assert.LessOrEqualf(t, i, 2, "Mapper was called with unexpected value: %d", i)

		if i == 2 {
			return 0, errors.New("Invalid value")
		}

		return i, nil
	}

	output, err := Map(iter, mapper).Collect()
	assert.Empty(t, output)
	assert.ErrorContains(t, err, "Invalid value")
}
