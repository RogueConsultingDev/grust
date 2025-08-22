package option

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSome_ReturnsNewSomeOption(t *testing.T) {
	v := fake.Lorem().Word()

	res := Some(v)

	s, ok := res.(some[string])
	require.True(t, ok, "result should be a some[string]: %#v", res)

	expected := some[string]{
		val: v,
	}
	assert.Equal(t, expected, s)
}

func TestSome_IsNone(t *testing.T) {
	s := Some(fake.Int())

	assert.False(t, s.IsNone())
}

func TestSome_IsNoneOr(t *testing.T) {
	value := fake.Int()
	s := Some(value)

	for _, res := range []bool{true, false} {
		name := fmt.Sprintf("predicate returns %v", res)
		t.Run(name, func(t *testing.T) {
			called := false

			f := func(v int) bool {
				called = true

				assert.Equal(t, value, v)

				return res
			}

			assert.Equal(t, res, s.IsNoneOr(f))
			assert.True(t, called, "predicate should have been called")
		})
	}
}

func TestSome_IsSome(t *testing.T) {
	s := Some(fake.Int())

	assert.True(t, s.IsSome())
}

func TestSome_IsSomeAnd(t *testing.T) {
	value := fake.Int()
	s := Some(value)

	for _, res := range []bool{true, false} {
		name := fmt.Sprintf("predicate returns %v", res)
		t.Run(name, func(t *testing.T) {
			called := false

			f := func(v int) bool {
				called = true

				assert.Equal(t, value, v)

				return res
			}

			assert.Equal(t, res, s.IsSomeAnd(f))
			assert.True(t, called, "predicate should have been called")
		})
	}
}

func TestSome_Expect(t *testing.T) {
	value := fake.Int()
	s := Some(value)

	assert.Equal(t, value, s.Expect(fake.RandomStringWithLength(8)))
}

func TestSome_Unwrap(t *testing.T) {
	value := fake.Int()
	s := Some(value)

	assert.Equal(t, value, s.Unwrap())
}

func TestSome_UnwrapOr(t *testing.T) {
	value := fake.Int()
	s := Some(value)

	def := fake.Int()

	assert.Equal(t, value, s.UnwrapOr(def))
}

func TestSome_UnwrapOrElse(t *testing.T) {
	value := fake.Int()
	s := Some(value)

	def := fake.Int()
	f := func() int {
		assert.Fail(t, "should not be called")

		return def
	}

	assert.Equal(t, value, s.UnwrapOrElse(f))
}

func TestSome_UnwrapOrDefault(t *testing.T) {
	valStr := fake.RandomStringWithLength(8)
	valInt := fake.IntBetween(1, 100)
	valFloat := fake.Float64(2, 1, 100)

	assert.Equal(t, valStr, Some(valStr).UnwrapOrDefault())
	assert.Equal(t, valInt, Some(valInt).UnwrapOrDefault())
	assert.Equal( //nolint:testifylint  // Value should be _exactly_ equal
		t,
		valFloat,
		Some(valFloat).UnwrapOrDefault(),
	)
}

func TestSome_Inspect(t *testing.T) {
	value := fake.Int()
	s := Some(value)

	called := false
	f := func(v int) {
		called = true

		assert.Equal(t, value, v)
	}

	assert.Equal(t, s, s.Inspect(f))
	assert.True(t, called, "predicate should have been called")
}

func TestSome_Filter(t *testing.T) {
	value := fake.Int()
	s := Some(value)

	for _, res := range []bool{true, false} {
		name := fmt.Sprintf("predicate returns %v", res)
		t.Run(name, func(t *testing.T) {
			called := false
			f := func(v int) bool {
				called = true

				assert.Equal(t, value, v)

				return res
			}

			if res {
				assert.Equal(t, s, s.Filter(f))
			} else {
				assert.Equal(t, none[int]{}, s.Filter(f))
			}

			assert.True(t, called, "predicate should have been called")
		})
	}
}

func TestSome_Or(t *testing.T) {
	value := fake.Int()
	s := Some(value)

	other := Some(fake.Int())

	assert.Equal(t, s, s.Or(other))
}

func TestSome_OrElse(t *testing.T) {
	value := fake.Int()
	s := Some(value)

	f := func() Option[int] {
		assert.Fail(t, "should not be called")

		return Some(fake.Int())
	}

	assert.Equal(t, s, s.OrElse(f))
}

func TestSome_Xor(t *testing.T) {
	value := fake.Int()
	s := Some(value)

	t.Run("other is some", func(t *testing.T) {
		other := Some(fake.Int())
		assert.Equal(t, none[int]{}, s.Xor(other))
	})

	t.Run("other is none", func(t *testing.T) {
		other := none[int]{}
		assert.Equal(t, s, s.Xor(other))
	})
}

func TestSome_String(t *testing.T) {
	value := fake.Int()
	s := Some(value)

	expected := fmt.Sprintf("Some(%v)", value)

	assert.Equal(t, expected, s.String())
}
