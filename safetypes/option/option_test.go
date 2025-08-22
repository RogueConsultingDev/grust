package option

import (
	"strconv"
	"testing"
	"time"

	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
)

var fake = faker.NewWithSeedInt64(time.Now().UnixNano())

func TestFrom_ReturnsNewOptionFromArgs(t *testing.T) {
	type S struct {
		value int
	}

	t.Run("some", func(t *testing.T) {
		res1 := Of(fake.IntBetween(1, 100))
		assert.True(t, res1.IsSome())

		res2 := Of(fake.Float(2, 1, 100))
		assert.True(t, res2.IsSome())

		res3 := Of(fake.RandomStringWithLength(8))
		assert.True(t, res3.IsSome())

		res4 := Of(true)
		assert.True(t, res4.IsSome())

		res5 := Of(S{value: 1})
		assert.True(t, res5.IsSome())

		zeroInt := 0
		res6 := Of(&zeroInt)
		assert.True(t, res6.IsSome())

		zeroStr := ""
		res7 := Of(&zeroStr)
		assert.True(t, res7.IsSome())

		zeroS := S{} //nolint:exhaustruct  // We want the zero value
		res8 := Of(&zeroS)
		assert.True(t, res8.IsSome())
	})

	t.Run("none", func(t *testing.T) {
		res1 := Of(0)
		assert.True(t, res1.IsNone())

		res2 := Of(0.0)
		assert.True(t, res2.IsNone())

		res3 := Of("")
		assert.True(t, res3.IsNone())

		res4 := Of(false)
		assert.True(t, res4.IsNone())

		res5 := Of(S{}) //nolint:exhaustruct  // We want the zero value
		assert.True(t, res5.IsNone())

		res6 := Of((*string)(nil))
		assert.True(t, res6.IsNone())
	})
}

func TestMap_ReturnsANewOptionWithMappedValue(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		val := fake.Int()
		s := Some(val)

		expected := Some(strconv.Itoa(val))

		assert.Equal(t, expected, Map(s, strconv.Itoa))
	})

	t.Run("none", func(t *testing.T) {
		n := None()

		f := func(any) int {
			assert.Fail(t, "mapper should not have been called")

			return 0
		}

		assert.Equal(t, none[int]{}, Map(n, f))
	})
}

func TestMapOr_ReturnsTheMappedValueOrDefault(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		val := fake.Int()
		s := Some(val)

		def := fake.RandomStringWithLength(9)

		expected := strconv.Itoa(val)

		assert.Equal(t, expected, MapOr(s, def, strconv.Itoa))
	})

	t.Run("none", func(t *testing.T) {
		n := None()

		def := fake.RandomStringWithLength(9)
		f := func(any) string {
			assert.Fail(t, "mapper should not have been called")

			return ""
		}

		assert.Equal(t, def, MapOr(n, def, f))
	})
}

func TestMapOrElse_ReturnsTheMappedValueOrCallsDefaultFactory(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		val := fake.Int()
		s := Some(val)

		factory := func() string {
			assert.Fail(t, "factory should not have been called")

			return fake.RandomStringWithLength(9)
		}

		expected := strconv.Itoa(val)

		assert.Equal(t, expected, MapOrElse(s, factory, strconv.Itoa))
	})

	t.Run("none", func(t *testing.T) {
		n := None()

		mapper := func(any) string {
			assert.Fail(t, "mapper should not have been called")

			return ""
		}
		def := fake.RandomStringWithLength(9)
		factory := func() string {
			return def
		}

		assert.Equal(t, def, MapOrElse(n, factory, mapper))
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
		n := None()

		other := Some(fake.RandomStringWithLength(8))

		assert.Equal(t, none[string]{}, And(n, other))
	})
}

func TestAndThen_ReturnsMappedOrNone(t *testing.T) {
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
		n := None()

		f := func(_ any) Option[string] {
			assert.Fail(t, "mapper should not have been called")

			return none[string]{}
		}

		assert.Equal(t, none[string]{}, AndThen(n, f))
	})
}
