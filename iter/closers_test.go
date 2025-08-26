package betteriter

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollect_CollectsTheIterInASlice(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	output, err := New(values).Collect()
	require.NoError(t, err)

	assert.Len(t, output, len(values))

	for i, v := range output {
		assert.Equal(t, values[i], *v)

		// They should point to the same value
		assert.True(t, &values[i] == v)
	}
}

func TestCollect_PropagatesError(t *testing.T) {
	iter := &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			if !yield(1, nil) {
				return
			}
			if !yield(42, errors.New("some error")) {
				return
			}
			require.Fail(t, "Should not reach this point")
			if !yield(2, nil) {
				return
			}
		},
	}

	output, err := iter.Collect()
	require.ErrorContains(t, err, "some error")

	assert.Empty(t, output)
}

func TestReversed_CollectsTheIterInASliceInReverseOrder(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	output, err := New(values).Reversed()
	require.NoError(t, err)

	assert.Len(t, output, len(values))

	for i, v := range output {
		assert.Equal(t, values[4-i], *v)

		// They should point to the same value
		assert.True(t, &values[4-i] == v)
	}
}

func TestReversed_PropagatesError(t *testing.T) {
	iter := &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			if !yield(1, nil) {
				return
			}
			if !yield(42, errors.New("some error")) {
				return
			}
			require.Fail(t, "Should not reach this point")
			if !yield(2, nil) {
				return
			}
		},
	}

	output, err := iter.Reversed()
	require.ErrorContains(t, err, "some error")

	assert.Empty(t, output)
}

func TestAny_ReturnsTrueOnFirstElementThatMatchesThePredicate(t *testing.T) {
	predicate := func(*int) bool { return true }

	// Empty slice doesn't have any values
	output, err := New([]int{}).Any(predicate)
	require.NoError(t, err)
	assert.False(t, output)

	// A slice with only a zero value still has a value
	output, err = New([]int{0}).Any(predicate)
	require.NoError(t, err)
	assert.True(t, output)

	// A slice with only a nil pointer still has a value
	output, err = New([]*int{nil}).Any(func(i **int) bool { return true })
	require.NoError(t, err)
	assert.True(t, output)

	// A slice with many values...
	output, err = New([]int{1, 2, 3, 4, 5}).Any(predicate)
	require.NoError(t, err)
	assert.True(t, output)
}

func TestAny_ReturnsFalseIfNoElementsMatchThePredicate(t *testing.T) {
	output, err := New([]int{1, 3, 5}).Any(func(v *int) bool { return *v%2 == 0 })
	require.NoError(t, err)
	assert.False(t, output)
}

func TestAny_StopsAtTheFirstMatchingElement(t *testing.T) {
	iter := &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			if !yield(1, nil) {
				return
			}
			require.Fail(t, "Should not reach this point")
			if !yield(42, errors.New("some error")) {
				return
			}
		},
	}

	output, err := iter.Any(func(int) bool { return true })
	assert.NoError(t, err)
	assert.True(t, output)
}

func TestAny_PropagatesError(t *testing.T) {
	iter := &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			if !yield(42, errors.New("some error")) {
				return
			}
			require.Fail(t, "Should not reach this point")
		},
	}

	output, err := iter.Any(func(int) bool { return true })
	assert.ErrorContains(t, err, "some error")
	assert.False(t, output)
}

func TestAll_ReturnsTrueIfAllElementsMatchThePredicate(t *testing.T) {
	predicate := func(*int) bool { return true }

	// Empty slice doesn't have any values, so nothing _doesn't_ match
	output, err := New([]int{}).All(predicate)
	require.NoError(t, err)
	assert.True(t, output)

	// A slice with one value that matches
	output, err = New([]int{0}).All(predicate)
	require.NoError(t, err)
	assert.True(t, output)

	// A slice with many values...
	output, err = New([]int{1, 2, 3, 4, 5}).All(predicate)
	require.NoError(t, err)
	assert.True(t, output)
}

func TestAll_ReturnsFalseIfAnyElementDoesntMatchThePredicate(t *testing.T) {
	output, err := New([]int{1, 2, 3, 4, 5}).All(func(i *int) bool { return *i < 5 })
	assert.NoError(t, err)
	assert.False(t, output)
}

func TestAll_StopsAtTheFirstNonMatchingElement(t *testing.T) {
	iter := &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			if !yield(1, nil) {
				return
			}
			require.Fail(t, "Should not reach this point")
			if !yield(42, errors.New("some error")) {
				return
			}
		},
	}

	output, err := iter.All(func(int) bool { return false })
	assert.NoError(t, err)
	assert.False(t, output)
}

func TestAll_PropagatesError(t *testing.T) {
	iter := &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			if !yield(42, errors.New("some error")) {
				return
			}
			require.Fail(t, "Should not reach this point")
		},
	}

	output, err := iter.All(func(int) bool { return true })
	assert.ErrorContains(t, err, "some error")
	assert.False(t, output)
}

func TestFirst_ReturnsTheFirstElementOfTheIter(t *testing.T) {
	val := 42

	iter := &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			if !yield(val, nil) {
				return
			}
			require.Fail(t, "Should not reach this point")
		},
	}
	output, err := iter.First()
	assert.NoError(t, err)
	assert.Equal(t, val, output)
}

func TestFirst_ReturnsAnErrorIfIteratorIsEmpty(t *testing.T) {
	output, err := New([]int{}).First()
	assert.ErrorContains(t, err, "empty iterator")
	assert.Nil(t, output)
}

func TestFirst_PropagatesError(t *testing.T) {
	iter := &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			if !yield(42, errors.New("some error")) {
				return
			}
			require.Fail(t, "Should not reach this point")
		},
	}

	output, err := iter.First()
	assert.ErrorContains(t, err, "some error")
	assert.Zero(t, output)
}

func TestLast_ReturnsTheFirstElementOfTheIter(t *testing.T) {
	output, err := New([]int{1, 2, 3, 4, 5}).Last()
	assert.NoError(t, err)
	assert.Equal(t, 5, *output)
}

func TestLast_ReturnsAnErrorIfIteratorIsEmpty(t *testing.T) {
	output, err := New([]int{}).Last()
	assert.ErrorContains(t, err, "empty iterator")
	assert.Nil(t, output)
}

func TestLast_PropagatesError(t *testing.T) {
	iter := &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			if !yield(42, errors.New("some error")) {
				return
			}
			require.Fail(t, "Should not reach this point")
		},
	}

	output, err := iter.Last()
	assert.ErrorContains(t, err, "some error")
	assert.Zero(t, output)
}

func TestFind_ReturnsTheFirstElementThatMatches(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}

	output, ok, err := New(values).Find(func(v *int) bool { return *v%2 == 0 })
	require.NoError(t, err)

	require.True(t, ok)
	require.NotNil(t, output)
	assert.Equal(t, 2, *output)
}

func TestFind_ReturnsFalseIfNothingMatches(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}

	output, ok, err := New(values).Find(func(v *int) bool { return *v > 10 })
	require.NoError(t, err)

	assert.False(t, ok)
	assert.Zero(t, output)
}

func TestFind_PropagatesError(t *testing.T) {
	iter := &Iterator[int, int]{
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

	output, ok, err := iter.Find(func(int) bool { return false })
	assert.ErrorContains(t, err, "some error")
	assert.False(t, ok)
	assert.Zero(t, output)
}

func TestPosition_ReturnsThePositionOfTheFirstElementThatMatches(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}

	output, ok, err := New(values).Position(func(v *int) bool { return *v >= 4 })
	require.NoError(t, err)

	require.True(t, ok)
	require.NotNil(t, output)
	assert.Equal(t, 3, output)
}

func TestPosition_ReturnsFalseIfNothingMatches(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}

	output, ok, err := New(values).Position(func(v *int) bool { return *v > 10 })
	require.NoError(t, err)

	assert.False(t, ok)
	assert.Zero(t, output)
}

func TestPosition_PropagatesError(t *testing.T) {
	iter := &Iterator[int, int]{
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

	output, ok, err := iter.Position(func(int) bool { return false })
	assert.ErrorContains(t, err, "some error")
	assert.False(t, ok)
	assert.Zero(t, output)
}

func TestFold_AppliesTheFolderFunctionOnAllValues(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}

	iter := New2[int, int](values)

	res1, err := iter.Fold(0, func(cur int, v *int) int { return cur + *v })
	require.NoError(t, err)

	assert.Equal(t, 15, res1)

	iter2 := New2[int, string](values)
	res2, err := iter2.Fold(
		"result: ",
		func(cur string, v *int) string { return cur + strconv.Itoa(*v) },
	)
	require.NoError(t, err)

	assert.Equal(t, "result: 12345", res2)
}

func TestFold_PropagatesError(t *testing.T) {
	iter := &Iterator[int, int]{
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
	assert.ErrorContains(t, err, "some error")
	// folding should have stopped at 1
	assert.Equal(t, 1, output)
}

func TestCopied_CopiesTheIterInASlice(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	output, err := Copied(New(values))
	require.NoError(t, err)

	assert.Len(t, output, len(values))

	for i, v := range output {
		assert.Equal(t, values[i], v)

		// They should not point to the same value
		assert.False(t, &values[i] == &v)
	}
}

func TestCopied_PropagatesError(t *testing.T) {
	iter := &Iterator[*int, any]{
		it: func(yield func(*int, error) bool) {
			if !yield(ptr(1), nil) {
				return
			}
			if !yield(ptr(42), errors.New("some error")) {
				return
			}
			require.Fail(t, "Should not reach this point")
			if !yield(ptr(2), nil) {
				return
			}
		},
	}

	output, err := Copied(iter)
	require.ErrorContains(t, err, "some error")

	assert.Empty(t, output)
}
