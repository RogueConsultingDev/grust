package st

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrom_ReturnsNewOptionFromArgs(t *testing.T) {
	type S struct {
		value int
	}

	t.Run("some", func(t *testing.T) {
		res1 := OptionOf(fake.IntBetween(1, 100))
		assert.True(t, res1.IsSome())

		res2 := OptionOf(fake.Float(2, 1, 100))
		assert.True(t, res2.IsSome())

		res3 := OptionOf(fake.RandomStringWithLength(8))
		assert.True(t, res3.IsSome())

		res4 := OptionOf(true)
		assert.True(t, res4.IsSome())

		res5 := OptionOf(S{value: 1})
		assert.True(t, res5.IsSome())

		zeroInt := 0
		res6 := OptionOf(&zeroInt)
		assert.True(t, res6.IsSome())

		zeroStr := ""
		res7 := OptionOf(&zeroStr)
		assert.True(t, res7.IsSome())

		zeroS := S{} //nolint:exhaustruct  // We want the zero value
		res8 := OptionOf(&zeroS)
		assert.True(t, res8.IsSome())
	})

	t.Run("none", func(t *testing.T) {
		res1 := OptionOf(0)
		assert.True(t, res1.IsNone())

		res2 := OptionOf(0.0)
		assert.True(t, res2.IsNone())

		res3 := OptionOf("")
		assert.True(t, res3.IsNone())

		res4 := OptionOf(false)
		assert.True(t, res4.IsNone())

		res5 := OptionOf(S{}) //nolint:exhaustruct  // We want the zero value
		assert.True(t, res5.IsNone())

		res6 := OptionOf((*string)(nil))
		assert.True(t, res6.IsNone())
	})
}

func TestMapOption_ReturnsANewOptionWithMappedOptionValue(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		val := fake.Int()
		s := Some(val)

		expected := Some(strconv.Itoa(val))

		assert.Equal(t, expected, MapOption(s, strconv.Itoa))
	})

	t.Run("none", func(t *testing.T) {
		n := None[int]()

		f := func(int) int {
			assert.Fail(t, "mapper should not have been called")

			return 0
		}

		assert.Equal(t, None[int](), MapOption(n, f))
	})
}

func TestMapOptionOr_ReturnsTheMappedOptionValueOrDefault(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		val := fake.Int()
		s := Some(val)

		def := fake.RandomStringWithLength(9)

		expected := strconv.Itoa(val)

		assert.Equal(t, expected, MapOptionOr(s, def, strconv.Itoa))
	})

	t.Run("none", func(t *testing.T) {
		n := None[int]()

		def := fake.RandomStringWithLength(9)
		f := func(int) string {
			assert.Fail(t, "mapper should not have been called")

			return ""
		}

		assert.Equal(t, def, MapOptionOr(n, def, f))
	})
}

func TestMapOptionOrElse_ReturnsTheMappedOptionValueOrCallsDefaultFactory(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		val := fake.Int()
		s := Some(val)

		factory := func() string {
			assert.Fail(t, "factory should not have been called")

			return fake.RandomStringWithLength(9)
		}

		expected := strconv.Itoa(val)

		assert.Equal(t, expected, MapOptionOrElse(s, factory, strconv.Itoa))
	})

	t.Run("none", func(t *testing.T) {
		n := None[int]()

		mapper := func(int) string {
			assert.Fail(t, "mapper should not have been called")

			return ""
		}
		def := fake.RandomStringWithLength(9)
		factory := func() string {
			return def
		}

		assert.Equal(t, def, MapOptionOrElse(n, factory, mapper))
	})
}

func TestAnd_ReturnsOtherOrNone(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		val := fake.Int()
		s := Some(val)

		other := Some(fake.RandomStringWithLength(8))

		assert.Equal(t, other, And(s, other))
	})

	t.Run("none", func(t *testing.T) {
		n := None[int]()

		other := Some(fake.RandomStringWithLength(8))

		assert.Equal(t, None[string](), And(n, other))
	})
}

func TestAndThen_ReturnsMappedOptionOrNone(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		val := fake.Int()
		s := Some(val)

		other := Some(fake.RandomStringWithLength(8))
		f := func(o int) Option[string] {
			assert.Equal(t, val, o)

			return other
		}

		assert.Equal(t, other, AndThen(s, f))
	})

	t.Run("none", func(t *testing.T) {
		n := None[int]()

		f := func(int) Option[string] {
			assert.Fail(t, "mapper should not have been called")

			return &none[string]{}
		}

		assert.Equal(t, None[string](), AndThen(n, f))
	})
}
