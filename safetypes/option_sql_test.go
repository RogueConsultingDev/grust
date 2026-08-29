package st

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOption_Scan(t *testing.T) {
	t.Run("scan int", func(t *testing.T) {
		o := Some(fake.Int())

		val := fake.Int()
		err := o.Scan(val)
		require.NoError(t, err)

		assert.Equal(t, Some(val), o)
	})

	t.Run("scan float", func(t *testing.T) {
		o := Some(fake.Float64(2, 0, 100000))

		val := fake.Float64(2, 0, 100000)
		err := o.Scan(val)
		require.NoError(t, err)

		assert.Equal(t, Some[float64](val), o)
	})

	t.Run("scan string", func(t *testing.T) {
		o := Some(fake.RandomStringWithLength(8))

		val := fake.RandomStringWithLength(8)
		err := o.Scan(val)
		require.NoError(t, err)

		assert.Equal(t, Some[string](val), o)
	})

	t.Run("scan bool", func(t *testing.T) {
		o := Some(fake.Bool())

		val := fake.Bool()
		err := o.Scan(val)
		require.NoError(t, err)

		assert.Equal(t, Some[bool](val), o)
	})

	t.Run("scan to pointer", func(t *testing.T) {
		o := Some(ptr(fake.Int()))

		val := fake.Int()
		err := o.Scan(val)
		require.NoError(t, err)

		assert.Equal(t, Some[*int](ptr(val)), o)
	})

	t.Run("type conversion", func(t *testing.T) {
		// We should be able to assign an int32 to an int64
		o := Some[int64](0)

		val := fake.Int32()
		err := o.Scan(val)
		require.NoError(t, err)

		assert.Equal(t, Some[int64](int64(val)), o)
	})

	t.Run("nil value scans to none", func(t *testing.T) {
		o := Some("")

		err := o.Scan(nil)
		require.NoError(t, err)

		assert.Equal(t, None[string](), o)
	})

	t.Run("incompatible types", func(t *testing.T) {
		o := Some(0)

		err := o.Scan(fake.RandomStringWithLength(8))
		require.ErrorContains(t, err, "invalid type for *st.Option[int]: string")
	})

	t.Run("can't assign to nil ptr", func(t *testing.T) {
		o := Some[*string](nil)

		err := o.Scan(fake.RandomStringWithLength(8))
		require.ErrorContains(t, err, "can't assign value to nil pointer of type *string")
	})
}

func TestOption_Value(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		val := fake.Int()
		o := Some(val)

		res, err := o.Value()
		require.NoError(t, err)

		assert.Equal(t, val, res)
	})

	t.Run("none", func(t *testing.T) {
		o := None[int]()

		res, err := o.Value()
		require.NoError(t, err)

		assert.Nil(t, res)
	})
}
