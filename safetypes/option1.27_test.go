//go:build go1.27

package st

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOption_Map(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		value := fake.Int()
		o := Some(value)

		res := o.Map(strconv.Itoa)
		expected := strconv.Itoa(value)

		assert.Equal(t, Some(expected), res)
	})

	t.Run("none", func(t *testing.T) {
		o := None[int]()

		f := func(int) string {
			assert.Fail(t, "mapper should not have been called")

			return ""
		}

		res := o.Map(f)

		assert.Equal(t, None[string](), res)
	})
}

func TestOption_MapOr(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		value := fake.Int()
		o := Some(value)
		defaultVal := fake.RandomStringWithLength(8)

		res := o.MapOr(defaultVal, strconv.Itoa)
		expected := strconv.Itoa(value)

		assert.Equal(t, expected, res)
	})

	t.Run("none", func(t *testing.T) {
		o := None[int]()
		defaultVal := fake.RandomStringWithLength(8)

		f := func(int) string {
			assert.Fail(t, "mapper should not have been called")

			return ""
		}

		res := o.MapOr(defaultVal, f)

		assert.Equal(t, defaultVal, res)
	})
}

func TestOption_MapOrElse(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		value := fake.Int()
		o := Some(value)
		defaultVal := fake.RandomStringWithLength(8)

		res := o.MapOrElse(func() string { return defaultVal }, strconv.Itoa)
		expected := strconv.Itoa(value)

		assert.Equal(t, expected, res)
	})

	t.Run("none", func(t *testing.T) {
		o := None[int]()
		defaultVal := fake.RandomStringWithLength(8)

		f := func(int) string {
			assert.Fail(t, "mapper should not have been called")

			return ""
		}

		res := o.MapOrElse(func() string { return defaultVal }, f)

		assert.Equal(t, defaultVal, res)
	})
}

func TestOption_And(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		o := Some(1)

		value := fake.RandomStringWithLength(8)
		other := Some(value)

		res := o.And(other)

		assert.Equal(t, other, res)
	})

	t.Run("none", func(t *testing.T) {
		o := None[int]()
		value := fake.RandomStringWithLength(8)
		other := Some(value)

		res := o.And(other)

		assert.Equal(t, None[string](), res)
	})
}

func TestOption_AndThen(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		o := Some(1)

		value := fake.RandomStringWithLength(8)
		other := Some(value)

		res := o.AndThen(func(int) *Option[string] { return other })

		assert.Equal(t, other, res)
	})

	t.Run("none", func(t *testing.T) {
		o := None[int]()

		f := func(int) *Option[string] {
			assert.Fail(t, "mapper should not have been called")

			return Some("")
		}

		res := o.AndThen(f)

		assert.Equal(t, None[string](), res)
	})
}
