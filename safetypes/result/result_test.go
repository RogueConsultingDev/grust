package result

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
)

var fake = faker.NewWithSeedInt64(time.Now().UnixNano())

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

		res := Of(f())

		assert.True(t, res.IsOk())
	})

	t.Run("err", func(t *testing.T) {
		v := fake.RandomStringWithLength(8)
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))

		f := func() (string, error) {
			return v, err
		}

		res := Of(f())

		assert.True(t, res.IsErr())
	})
}

func TestMap_ReturnsANewResultWithMappedValue(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := Ok[int, error](value)

		expected := Ok[string, error](strconv.Itoa(value))

		assert.Equal(t, expected, Map(o, strconv.Itoa))
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
		e := Err[any, error](err)

		f := func(any) int {
			assert.Fail(t, "mapper should not have been called")

			return 0
		}

		expected := errT[int, error]{err}

		assert.Equal(t, expected, Map(e, f))
	})
}

func TestMapOr_ReturnsTheMappedValueOrDefault(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := Ok[int, error](value)

		def := fake.RandomStringWithLength(9)

		expected := strconv.Itoa(value)

		assert.Equal(t, expected, MapOr(o, def, strconv.Itoa))
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
		e := Err[any, error](err)

		def := fake.RandomStringWithLength(9)
		f := func(any) string {
			assert.Fail(t, "mapper should not have been called")

			return ""
		}

		assert.Equal(t, def, MapOr(e, def, f))
	})
}

func TestMapOrElse_ReturnsTheMappedValueOrCallsDefaultFactory(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := Ok[int, error](value)

		factory := func() string {
			assert.Fail(t, "factory should not have been called")

			return fake.RandomStringWithLength(9)
		}

		expected := strconv.Itoa(value)

		assert.Equal(t, expected, MapOrElse(o, factory, strconv.Itoa))
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
		e := Err[any, error](err)

		mapper := func(any) string {
			assert.Fail(t, "mapper should not have been called")

			return ""
		}
		def := fake.RandomStringWithLength(9)
		factory := func() string {
			return def
		}

		assert.Equal(t, def, MapOrElse(e, factory, mapper))
	})
}

func TestMapErr_ReturnsANewResultWithMappedValue(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := Ok[int, error](value)

		f := func(error) *MockError {
			assert.Fail(t, "mapper should not have been called")

			return nil
		}

		expected := ok[int, *MockError]{value}

		assert.Equal(t, expected, MapErr(o, f))
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
		e := Err[any, error](err)

		newE := &MockError{e: fmt.Errorf("mapped error: %w", err).Error()}
		f := func(e error) *MockError {
			assert.Equal(t, e, err)

			return newE
		}

		expected := errT[any, *MockError]{err: newE}

		assert.Equal(t, expected, MapErr(e, f))
	})
}
