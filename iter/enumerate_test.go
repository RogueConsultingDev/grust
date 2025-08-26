package betteriter

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnumerate_ReturnsValuesAndTheirIndex(t *testing.T) {
	values := []string{"one", "two", "three"}
	iter := New(values)

	output, err := NewEnumeratorFrom(iter).Collect()
	require.NoError(t, err)

	expected := []Enumerator[*string]{
		{0, &values[0]},
		{1, &values[1]},
		{2, &values[2]},
	}
	assert.Equal(t, expected, output)
}

func TestEnumerate_StopsOnError(t *testing.T) {
	values := []int{1, 2, 3}
	iter := New2[int, int](values)

	mapper := func(i *int) (int, error) {
		// We will error on value 2, so mapper should never be called with value 3
		assert.LessOrEqualf(t, *i, 2, "Mapper was called with unexpected value: %d", i)

		if *i == 2 {
			return 0, errors.New("Invalid value")
		}

		return 0, nil
	}

	output, err := NewEnumeratorFrom(iter.Map(mapper)).Collect()
	assert.Empty(t, output)
	assert.ErrorContains(t, err, "Invalid value")
}
