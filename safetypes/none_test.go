package st

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNone_ReturnsNewNoneOption(t *testing.T) {
	res := None[int]()

	n, ok := res.(none[int])
	require.True(t, ok, "result should be a none: %#v", res)

	assert.Equal(t, none[int]{}, n)
}

func TestNone_IsNone(t *testing.T) {
	n := None[int]()

	assert.True(t, n.IsNone())
}

func TestNone_IsNoneOr(t *testing.T) {
	n := None[int]()

	for _, res := range []bool{true, false} {
		name := fmt.Sprintf("predicate returns %v", res)
		t.Run(name, func(t *testing.T) {
			f := func(int) bool {
				// The predicate should _not_ be called, there's just no point in doing so.
				assert.Fail(t, "predicate should not have been called")

				return res
			}

			assert.True(t, n.IsNoneOr(f))
		})
	}
}

func TestNone_IsSome(t *testing.T) {
	n := None[int]()

	assert.False(t, n.IsSome())
}

func TestNone_IsSomeAnd(t *testing.T) {
	n := None[int]()

	for _, res := range []bool{true, false} {
		name := fmt.Sprintf("predicate returns %v", res)
		t.Run(name, func(t *testing.T) {
			f := func(int) bool {
				// The predicate should _not_ be called, there's just no point in doing so.
				assert.Fail(t, "predicate should not have been called")

				return res
			}

			assert.False(t, n.IsSomeAnd(f))
		})
	}
}

func TestNone_Expect(t *testing.T) {
	s := None[int]()

	msg := fake.RandomStringWithLength(8)
	assert.PanicsWithError(t, msg, func() {
		s.Expect(msg)
	})
}

func TestNone_Unwrap(t *testing.T) {
	s := None[int]()

	assert.PanicsWithError(t, "called `Option.Unwrap()` on a `None` value", func() {
		s.Unwrap()
	})
}

func TestNone_UnwrapOr(t *testing.T) {
	s := None[string]()
	def := fake.RandomStringWithLength(8)

	assert.Equal(t, def, s.UnwrapOr(def))
}

func TestNone_UnwrapOrElse(t *testing.T) {
	s := None[string]()

	def := fake.RandomStringWithLength(8)
	f := func() string {
		return def
	}

	assert.Equal(t, def, s.UnwrapOrElse(f))
}

func TestNone_UnwrapOrDefault(t *testing.T) {
	noneStr := none[string]{}
	noneInt := none[int]{}
	noneFloat := none[float64]{}

	assert.Zero( //nolint:testifylint  // Using assert.Zero for consistency
		t,
		noneStr.UnwrapOrDefault(),
	)
	assert.Zero(t, noneInt.UnwrapOrDefault())
	assert.Zero(t, noneFloat.UnwrapOrDefault())
}

func TestNone_Inspect(t *testing.T) {
	n := None[int]()

	f := func(int) {
		assert.Fail(t, "predicate should not have been called")
	}

	assert.Equal(t, n, n.Inspect(f))
}

func TestNone_Filter(t *testing.T) {
	n := None[int]()

	for _, res := range []bool{true, false} {
		name := fmt.Sprintf("predicate returns %v", res)
		t.Run(name, func(t *testing.T) {
			f := func(int) bool {
				// The predicate should _not_ be called, there's just no point in doing so.
				assert.Fail(t, "predicate should not have been called")

				return res
			}

			assert.Equal(t, n, n.Filter(f))
		})
	}
}

func TestNone_Or(t *testing.T) {
	n := none[int]{}

	t.Run("other is some", func(t *testing.T) {
		other := Some(fake.Int())

		assert.Equal(t, other, n.Or(other))
	})

	t.Run("other is none", func(t *testing.T) {
		other := none[int]{}

		assert.Equal(t, other, n.Or(other))
	})
}

func TestNone_OrElse(t *testing.T) {
	n := none[int]{}

	t.Run("other is some", func(t *testing.T) {
		other := Some(fake.Int())

		f := func() Option[int] {
			return other
		}

		assert.Equal(t, other, n.OrElse(f))
	})

	t.Run("other is none", func(t *testing.T) {
		other := none[int]{}

		f := func() Option[int] {
			return other
		}

		assert.Equal(t, other, n.OrElse(f))
	})
}

func TestNone_Xor(t *testing.T) {
	n := none[int]{}

	t.Run("other is some", func(t *testing.T) {
		other := Some(fake.Int())
		assert.Equal(t, other, n.Xor(other))
	})

	t.Run("other is none", func(t *testing.T) {
		other := none[int]{}
		assert.Equal(t, n, n.Xor(other))
	})
}

func TestNone_String(t *testing.T) {
	n := none[int]{}

	assert.Equal(t, "None", n.String())
}
