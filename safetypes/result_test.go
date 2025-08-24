package st

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockError struct {
	e string
}

func (t *MockError) Error() string {
	return t.e
}

func TestFrom_ReturnsNewResultFromArgs(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		v := fake.RandomStringWithLength(8)

		f := func() (string, error) { //nolint:unparam // Needs this signature
			return v, nil
		}

		res := ResultOf(f())

		assert.True(t, res.IsOk())
	})

	t.Run("err", func(t *testing.T) {
		v := fake.RandomStringWithLength(8)
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))

		f := func() (string, error) {
			return v, err
		}

		res := ResultOf(f())

		assert.True(t, res.IsErr())
	})
}

func TestMapResult_ReturnsANewResultWithMappedResultValue(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := Ok[int](value)

		expected := Ok[string](strconv.Itoa(value))

		assert.Equal(t, expected, MapResult(o, strconv.Itoa))
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
		e := Err[any](err)

		f := func(any) int {
			assert.Fail(t, "mapper should not have been called")

			return 0
		}

		expected := &errT[int]{err}

		assert.Equal(t, expected, MapResult(e, f))
	})
}

func TestMapResultOr_ReturnsTheMappedResultValueOrDefault(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := Ok[int](value)

		def := fake.RandomStringWithLength(9)

		expected := strconv.Itoa(value)

		assert.Equal(t, expected, MapResultOr(o, def, strconv.Itoa))
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
		e := Err[any](err)

		def := fake.RandomStringWithLength(9)
		f := func(any) string {
			assert.Fail(t, "mapper should not have been called")

			return ""
		}

		assert.Equal(t, def, MapResultOr(e, def, f))
	})
}

func TestMapResultOrElse_ReturnsTheMappedResultValueOrCallsDefaultFactory(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := Ok[int](value)

		factory := func() string {
			assert.Fail(t, "factory should not have been called")

			return fake.RandomStringWithLength(9)
		}

		expected := strconv.Itoa(value)

		assert.Equal(t, expected, MapResultOrElse(o, factory, strconv.Itoa))
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
		e := Err[any](err)

		mapper := func(any) string {
			assert.Fail(t, "mapper should not have been called")

			return ""
		}
		def := fake.RandomStringWithLength(9)
		factory := func() string {
			return def
		}

		assert.Equal(t, def, MapResultOrElse(e, factory, mapper))
	})
}

func TestMapResultErr_ReturnsANewResultWithMappedResultValue(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := Ok[int](value)

		f := func(error) error {
			assert.Fail(t, "mapper should not have been called")

			return nil
		}

		expected := &ok[int]{value}

		assert.Equal(t, expected, MapResultErr(o, f))
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
		e := Err[any](err)

		newE := &MockError{e: fmt.Errorf("mapped error: %w", err).Error()}
		f := func(e error) error {
			assert.Equal(t, e, err)

			return newE
		}

		expected := &errT[any]{err: newE}

		assert.Equal(t, expected, MapResultErr(e, f))
	})
}
