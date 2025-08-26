package it

import (
	"testing"

	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fake = faker.New()

func TestNewRepeat_ReturnsAnInfiniteIterator(t *testing.T) {
	val := 4 // Random value chosen by a fair dice roll

	repeat := NewRepeat(val)

	c := 0

	for v := range repeat.it {
		require.Equal(t, val, v)

		c++

		// Let's stop right there
		if c == 10_000 {
			break
		}
	}
}

func TestNewRepeatN_ReturnsAnInteratorOfNValues(t *testing.T) {
	val := 4 // Random value chosen by a fair dice roll
	n := fake.IntBetween(1, 10_000)

	repeat := NewRepeatN(val, n)

	c := 0

	for v := range repeat.it {
		require.Equal(t, val, v)

		c++

		// Safety stop
		if c == n+1 {
			require.Fail(t, "Iterator didn't stop", "It should have stopped at n=%d", n)
		}
	}

	assert.Equal(t, n, c)
}

func TestNewRepeatN_ReturnsNothingIfNIsZero(t *testing.T) {
	val := 4 // Random value chosen by a fair dice roll

	repeat := NewRepeatN(val, 0)

	for range repeat.it {
		require.Fail(t, "The iterator returned something")
	}
}

func TestZip_ReturnsValuesFromTwoSlices(t *testing.T) {
	iterA := []int{1, 2, 3, 4, 5}
	iterB := []string{"one", "two", "three", "four", "five"}

	idx := 0
	for v := range Zip(iterA, iterB).it {
		assert.Equal(t, iterA[idx], v.A)
		assert.Equal(t, iterB[idx], v.B)

		idx++
	}
}

func TestZip_StopsAtTheShortestOfTwoSlices(t *testing.T) {
	iterA := []int{1, 2, 3, 4, 5}
	iterB := []string{"one", "two", "three", "four", "five"}

	// Slice A is shorter
	idx := 0
	for range Zip(iterA[:3], iterB).it {
		idx++
	}

	assert.Equal(t, 3, idx)

	// Slice B is shorter
	idx = 0
	for range Zip(iterA, iterB[:3]).it {
		idx++
	}

	assert.Equal(t, 3, idx)
}

func TestZipEq_ReturnsValuesFromTwoSlices(t *testing.T) {
	iterA := []int{1, 2, 3, 4, 5}
	iterB := []string{"one", "two", "three", "four", "five"}

	idx := 0
	for v := range ZipEq(iterA, iterB).it {
		assert.Equal(t, iterA[idx], v.A)
		assert.Equal(t, iterB[idx], v.B)

		idx++
	}
}

func TestZipEq_ReturnsAnErrorIfSlicesAreNotSameLength(t *testing.T) {
	iterA := []int{1, 2, 3, 4, 5}
	iterB := []string{"one", "two", "three", "four", "five"}

	// Slice A is shorter
	output, err := ZipEq(iterA[:3], iterB).Collect()
	require.ErrorContains(t, err, "slices are not the same length")
	assert.Empty(t, output)

	// Slice B is shorter
	output, err = ZipEq(iterA, iterB[:3]).Collect()
	require.ErrorContains(t, err, "slices are not the same length")
	assert.Empty(t, output)
}
