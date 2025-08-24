package st

import (
	"errors"
	"fmt"
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

			return newNone[string]()
		}

		assert.Equal(t, None[string](), AndThen(n, f))
	})
}

func TestOption_IsNone(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()
		o := newSome[int](val)

		assert.False(t, o.IsNone())
	})

	t.Run("None", func(t *testing.T) {
		o := newNone[int]()

		assert.True(t, o.IsNone())
	})
}

func TestOption_IsNoneOr(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()
		o := newSome[int](val)

		for _, res := range []bool{true, false} {
			name := fmt.Sprintf("predicate returns %v", res)
			t.Run(name, func(t *testing.T) {
				called := false

				f := func(v int) bool {
					called = true

					assert.Equal(t, val, v)

					return res
				}

				assert.Equal(t, res, o.IsNoneOr(f))
				assert.True(t, called, "predicate should have been called")
			})
		}
	})

	t.Run("None", func(t *testing.T) {
		o := newNone[int]()

		for _, res := range []bool{true, false} {
			name := fmt.Sprintf("predicate returns %v", res)
			t.Run(name, func(t *testing.T) {
				f := func(int) bool {
					assert.Fail(t, "predicate should not be called")

					return res
				}

				assert.True(t, o.IsNoneOr(f))
			})
		}
	})
}

func TestOption_IsSome(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()
		o := newSome[int](val)

		assert.True(t, o.IsSome())
	})

	t.Run("None", func(t *testing.T) {
		o := newNone[int]()

		assert.False(t, o.IsSome())
	})
}

func TestOption_IsSomeAnd(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()
		o := newSome[int](val)

		for _, res := range []bool{true, false} {
			name := fmt.Sprintf("predicate returns %v", res)
			t.Run(name, func(t *testing.T) {
				called := false

				f := func(v int) bool {
					called = true

					assert.Equal(t, val, v)

					return res
				}

				assert.Equal(t, res, o.IsSomeAnd(f))
				assert.True(t, called, "predicate should have been called")
			})
		}
	})

	t.Run("None", func(t *testing.T) {
		o := newNone[int]()

		for _, res := range []bool{true, false} {
			name := fmt.Sprintf("predicate returns %v", res)
			t.Run(name, func(t *testing.T) {
				f := func(int) bool {
					assert.Fail(t, "predicate should not be called")

					return res
				}

				assert.False(t, o.IsSomeAnd(f))
			})
		}
	})
}

func TestOption_Expect(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()
		o := newSome[int](val)

		assert.Equal(t, val, o.Expect(fake.RandomStringWithLength(8)))
	})

	t.Run("None", func(t *testing.T) {
		o := newNone[int]()

		msg := fake.RandomStringWithLength(8)
		assert.PanicsWithError(t, msg, func() {
			o.Expect(msg)
		})
	})
}

func TestOption_Unwrap(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()
		o := newSome[int](val)

		assert.Equal(t, val, o.Unwrap())
	})

	t.Run("None", func(t *testing.T) {
		o := newNone[int]()

		assert.PanicsWithError(t, "called `Option.Unwrap()` on a `None` value", func() {
			o.Unwrap()
		})
	})
}

func TestOption_UnwrapOr(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()
		o := newSome[int](val)

		assert.Equal(t, val, o.UnwrapOr(fake.Int()))
	})

	t.Run("None", func(t *testing.T) {
		o := newNone[int]()

		fallback := fake.Int()

		assert.Equal(t, fallback, o.UnwrapOr(fallback))
	})
}

func TestOption_UnwrapOrElse(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()
		o := newSome[int](val)

		predicate := func() int {
			assert.Fail(t, "predicate should not be called")

			return 0
		}

		assert.Equal(t, val, o.UnwrapOrElse(predicate))
	})

	t.Run("None", func(t *testing.T) {
		o := newNone[int]()

		fallback := fake.Int()

		predicate := func() int {
			return fallback
		}

		assert.Equal(t, fallback, o.UnwrapOrElse(predicate))
	})
}

func TestOption_UnwrapOrDefault(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()
		o := newSome[int](val)

		assert.Equal(t, val, o.UnwrapOrDefault())
	})

	t.Run("None", func(t *testing.T) {
		oi := newNone[int]()

		assert.Equal(t, 0, oi.UnwrapOrDefault())

		os := newNone[string]()

		assert.Empty(t, os.UnwrapOrDefault())
	})
}

func TestAsOkOr(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		val := fake.Int()
		o := newSome[int](val)

		expected := Ok[int](val)

		assert.Equal(t, expected, o.AsOkOr(errors.New(fake.RandomStringWithLength(8))))
	})

	t.Run("none", func(t *testing.T) {
		o := newNone[int]()
		err := errors.New(fake.RandomStringWithLength(8))

		expected := Err[int](err)

		assert.Equal(t, expected, o.AsOkOr(err))
	})
}

func TestAsOkOrElse(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		val := fake.Int()
		o := newSome[int](val)

		f := func() error {
			assert.Fail(t, "should not be called")

			return errors.New(fake.RandomStringWithLength(8))
		}

		expected := Ok[int](val)

		assert.Equal(t, expected, o.AsOkOrElse(f))
	})

	t.Run("none", func(t *testing.T) {
		o := newNone[int]()
		err := errors.New(fake.RandomStringWithLength(8))

		f := func() error {
			return err
		}

		expected := Err[int](err)

		assert.Equal(t, expected, o.AsOkOrElse(f))
	})
}

func TestOption_Inspect(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()
		o := newSome[int](val)

		called := false
		predicate := func(v int) {
			called = true

			assert.Equal(t, val, v)
		}

		assert.Same(t, o, o.Inspect(predicate))
		assert.True(t, called)
	})

	t.Run("None", func(t *testing.T) {
		o := newNone[int]()

		predicate := func(int) {
			assert.Fail(t, "predicate should not be called")
		}

		assert.Same(t, o, o.Inspect(predicate))
	})
}

func TestOption_Filter(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()
		o := newSome[int](val)

		for _, res := range []bool{true, false} {
			name := fmt.Sprintf("predicate returns %v", res)
			t.Run(name, func(t *testing.T) {
				called := false
				f := func(v int) bool {
					called = true

					assert.Equal(t, val, v)

					return res
				}

				if res {
					assert.Equal(t, o, o.Filter(f))
				} else {
					assert.Equal(t, newNone[int](), o.Filter(f))
				}

				assert.True(t, called, "predicate should have been called")
			})
		}
	})

	t.Run("None", func(t *testing.T) {
		o := newNone[int]()

		for _, res := range []bool{true, false} {
			name := fmt.Sprintf("predicate returns %v", res)
			t.Run(name, func(t *testing.T) {
				f := func(int) bool {
					assert.Fail(t, "predicate should not be called")

					return res
				}

				assert.Equal(t, newNone[int](), o.Filter(f))
			})
		}
	})
}

func TestOption_Or(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()
		o := newSome[int](val)

		other := newSome[int](fake.Int())

		assert.Equal(t, o, o.Or(other))
	})

	t.Run("None", func(t *testing.T) {
		o := newNone[int]()

		other := newSome[int](fake.Int())

		assert.Equal(t, other, o.Or(other))
	})
}

func TestOption_OrElse(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()
		o := newSome[int](val)

		otherFactory := func() Option[int] {
			assert.Fail(t, "factory should not be called")

			return newSome[int](fake.Int())
		}

		assert.Equal(t, o, o.OrElse(otherFactory))
	})

	t.Run("None", func(t *testing.T) {
		o := newNone[int]()

		other := newSome[int](fake.Int())
		otherFactory := func() Option[int] {
			return other
		}

		assert.Equal(t, other, o.OrElse(otherFactory))
	})
}

func TestOption_Xor(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()
		o := newSome[int](val)

		t.Run("other is some", func(t *testing.T) {
			other := newSome(fake.Int())
			assert.Equal(t, newNone[int](), o.Xor(other))
		})

		t.Run("other is none", func(t *testing.T) {
			other := newNone[int]()
			assert.Same(t, o, o.Xor(other))
		})
	})

	t.Run("None", func(t *testing.T) {
		o := newNone[int]()

		t.Run("other is some", func(t *testing.T) {
			other := newSome(fake.Int())
			assert.Same(t, other, o.Xor(other))
		})

		t.Run("other is none", func(t *testing.T) {
			other := newNone[int]()
			assert.Equal(t, newNone[int](), o.Xor(other))
		})
	})
}

func TestOption_Insert(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()
		o := newSome(val)

		newVal := fake.Int()
		res := o.Insert(newVal)

		assert.Equal(t, newVal, *res)

		expected := &option[int]{
			ok:  true,
			val: newVal,
		}
		assert.Equal(t, expected, o)

		// Assert that we can directly change the option's value with the pointer
		newVal = fake.Int()
		*res = newVal
		assert.Equal(t, newVal, o.Unwrap())
	})

	t.Run("None", func(t *testing.T) {
		o := newNone[int]()

		newVal := fake.Int()
		res := o.Insert(newVal)

		assert.Equal(t, newVal, *res)

		expected := &option[int]{
			ok:  true,
			val: newVal,
		}
		assert.Equal(t, expected, o)

		// Assert that we can directly change the option's value with the pointer
		newVal = fake.Int()
		*res = newVal
		assert.Equal(t, newVal, o.Unwrap())
	})
}

func TestOption_GetOrInsert(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()
		o := newSome(val)

		newVal := fake.Int()
		res := o.GetOrInsert(newVal)

		assert.Equal(t, val, *res)

		expected := &option[int]{
			ok:  true,
			val: val,
		}
		assert.Equal(t, expected, o)

		// Assert that we can directly change the option's value with the pointer
		newVal = fake.Int()
		*res = newVal
		assert.Equal(t, newVal, o.Unwrap())
	})

	t.Run("None", func(t *testing.T) {
		o := newNone[int]()

		newVal := fake.Int()
		res := o.GetOrInsert(newVal)

		assert.Equal(t, newVal, *res)

		expected := &option[int]{
			ok:  true,
			val: newVal,
		}
		assert.Equal(t, expected, o)

		// Assert that we can directly change the option's value with the pointer
		newVal = fake.Int()
		*res = newVal
		assert.Equal(t, newVal, o.Unwrap())
	})
}

func TestOption_GetOrInsertDefault(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()
		o := newSome(val)

		res := o.GetOrInsertDefault()

		assert.Equal(t, val, *res)

		expected := &option[int]{
			ok:  true,
			val: val,
		}
		assert.Equal(t, expected, o)

		// Assert that we can directly change the option's value with the pointer
		newVal := fake.Int()
		*res = newVal
		assert.Equal(t, newVal, o.Unwrap())
	})

	t.Run("None", func(t *testing.T) {
		o := newNone[int]()

		res := o.GetOrInsertDefault()

		assert.Equal(t, 0, *res)

		expected := &option[int]{
			ok:  true,
			val: 0,
		}
		assert.Equal(t, expected, o)

		// Assert that we can directly change the option's value with the pointer
		newVal := fake.Int()
		*res = newVal
		assert.Equal(t, newVal, o.Unwrap())
	})
}

func TestOption_GetOrInsertWith(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()
		o := newSome(val)

		newVal := fake.Int()
		factory := func() int {
			assert.Fail(t, "factory should not be called")

			return newVal
		}
		res := o.GetOrInsertWith(factory)

		assert.Equal(t, val, *res)

		expected := &option[int]{
			ok:  true,
			val: val,
		}
		assert.Equal(t, expected, o)

		// Assert that we can directly change the option's value with the pointer
		newVal = fake.Int()
		*res = newVal
		assert.Equal(t, newVal, o.Unwrap())
	})

	t.Run("None", func(t *testing.T) {
		o := newNone[int]()

		newVal := fake.Int()
		factory := func() int {
			return newVal
		}
		res := o.GetOrInsertWith(factory)

		assert.Equal(t, newVal, *res)

		expected := &option[int]{
			ok:  true,
			val: newVal,
		}
		assert.Equal(t, expected, o)

		// Assert that we can directly change the option's value with the pointer
		newVal = fake.Int()
		*res = newVal
		assert.Equal(t, newVal, o.Unwrap())
	})
}

func TestOption_Take(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()
		o := newSome(val)

		res := o.Take()

		expected := newSome(val)
		assert.Equal(t, expected, res)

		// Original opt should now be a None
		assert.True(t, o.IsNone())
	})

	t.Run("None", func(t *testing.T) {
		o := newNone[int]()

		res := o.Take()

		assert.True(t, res.IsNone())

		// Original opt should still be a None
		assert.True(t, o.IsNone())
	})
}

func TestOption_TakeIf(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		val := fake.Int()

		for _, res := range []bool{true, false} {
			name := fmt.Sprintf("predicate returns %v", res)
			o := newSome(val)

			t.Run(name, func(t *testing.T) {
				called := false
				f := func(v int) bool {
					called = true

					assert.Equal(t, val, v)

					return res
				}

				taken := o.TakeIf(f)

				assert.True(t, called, "predicate should have been called")

				if res {
					expected := newSome(val)
					assert.Equal(t, expected, taken)

					// Original opt should now be a None
					assert.True(t, o.IsNone())
				} else {
					expected := newNone[int]()
					assert.Equal(t, expected, taken)

					// Original opt should have stayed intact
					assert.True(t, o.IsSome())
					assert.Equal(t, val, o.Unwrap())
				}
			})
		}
	})

	t.Run("None", func(t *testing.T) {
		for _, res := range []bool{true, false} {
			name := fmt.Sprintf("predicate returns %v", res)
			o := newNone[int]()

			t.Run(name, func(t *testing.T) {
				f := func(int) bool {
					assert.Fail(t, "predicate should not be called")

					return res
				}

				taken := o.TakeIf(f)

				assert.True(t, taken.IsNone())

				// Original opt should still be a None
				assert.True(t, o.IsNone())
			})
		}
	})
}

func TestOption_String(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		value := fake.Int()
		o := newSome[int](value)

		expected := fmt.Sprintf("Some(%v)", value)

		assert.Equal(t, expected, o.String())
	})

	t.Run("None", func(t *testing.T) {
		o := newNone[int]()

		expected := "None"

		assert.Equal(t, expected, o.String())
	})
}

func newSome[T any](val T) *option[T] {
	o, ok := Some(val).(*option[T])
	if !ok {
		panic("expected *option[T]")
	}

	return o
}

func newNone[T any]() *option[T] {
	o, ok := None[T]().(*option[T])
	if !ok {
		panic("expected *option[T]")
	}

	return o
}
