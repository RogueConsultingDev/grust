//go:build go1.27

package st

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResult_Map(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := Ok[int](value)

		expected := Ok[string](strconv.Itoa(value))

		assert.Equal(t, expected, o.Map(strconv.Itoa))
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
		e := Err[any](err)

		f := func(any) int {
			assert.Fail(t, "mapper should not have been called")

			return 0
		}

		expected := Err[int](err)

		assert.Equal(t, expected, e.Map(f))
	})
}

func TestResult_MapOr(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := Ok[int](value)

		def := fake.RandomStringWithLength(9)

		expected := strconv.Itoa(value)

		assert.Equal(t, expected, o.MapOr(def, strconv.Itoa))
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
		e := Err[any](err)

		def := fake.RandomStringWithLength(9)
		f := func(any) string {
			assert.Fail(t, "mapper should not have been called")

			return ""
		}

		assert.Equal(t, def, e.MapOr(def, f))
	})
}

func TestResult_MapOrElse(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := Ok[int](value)

		factory := func() string {
			assert.Fail(t, "factory should not have been called")

			return fake.RandomStringWithLength(9)
		}

		expected := strconv.Itoa(value)

		assert.Equal(t, expected, o.MapOrElse(factory, strconv.Itoa))
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

		assert.Equal(t, def, e.MapOrElse(factory, mapper))
	})
}
