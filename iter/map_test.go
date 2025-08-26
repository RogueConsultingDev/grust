package betteriter

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMap_TransformsElements(t *testing.T) {
	values := []int{1, 2, 3}
	iter := New2[int, string](values)

	mapper := func(i *int) (string, error) {
		return strconv.Itoa(*i * *i), nil
	}

	output, err := iter.Map(mapper).Collect()

	require.NoError(t, err)

	assert.Equal(t, []*string{ptr("1"), ptr("4"), ptr("9")}, output)
}

func TestMap_IsLazy(t *testing.T) {
	values := []int{1, 2, 3}
	iter := New2[int, int](values)

	mapper := func(i *int) (int, error) {
		assert.LessOrEqualf(t, *i, 2, "Mapper was called with unexpected value: %d", i)

		return *i, nil
	}

	for v := range iter.Map(mapper).it {
		if *v == 2 {
			break
		}
	}
}

func TestMap_StopsOnError(t *testing.T) {
	values := []int{1, 2, 3}
	iter := New2[int, int](values)

	mapper := func(i *int) (int, error) {
		// We will error on value 2, so mapper should never be called with value 3
		assert.LessOrEqualf(t, *i, 2, "Mapper was called with unexpected value: %d", i)

		if *i == 2 {
			return 0, errors.New("Invalid value")
		}

		return *i, nil
	}

	output, err := iter.Map(mapper).Collect()
	assert.Empty(t, output)
	assert.ErrorContains(t, err, "Invalid value")
}

func TestMap_PropagatesError(t *testing.T) {
	iter := &Iterator[int, int]{
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

	mapper := func(i int) (int, error) {
		// We should only be called with value = 1
		assert.Equal(t, 1, i, "Mapper was called with unexpected value: %d", i)

		return i, nil
	}

	output, err := iter.Map(mapper).Collect()
	assert.Empty(t, output)
	assert.ErrorContains(t, err, "some error")
}
