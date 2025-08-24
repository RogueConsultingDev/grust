package st

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		expected := Err[int](err)

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

		expected := Ok[int](value)

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

		expected := Err[any](newE)

		assert.Equal(t, expected, MapResultErr(e, f))
	})
}

func TestResult_IsOk(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		val := fake.Int()
		r := newOk[int](val)

		assert.True(t, r.IsOk())
	})

	t.Run("Err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		r := newErr[any](err)

		assert.False(t, r.IsOk())
	})
}

func TestResult_IsOkAnd(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		val := fake.Int()
		r := newOk[int](val)

		for _, res := range []bool{true, false} {
			name := fmt.Sprintf("predicate returns %v", res)
			t.Run(name, func(t *testing.T) {
				called := false

				f := func(v int) bool {
					called = true

					assert.Equal(t, val, v)

					return res
				}

				assert.Equal(t, res, r.IsOkAnd(f))
				assert.True(t, called, "predicate should have been called")
			})
		}
	})

	t.Run("Err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		r := newErr[int](err)

		for _, res := range []bool{true, false} {
			name := fmt.Sprintf("predicate returns %v", res)
			t.Run(name, func(t *testing.T) {
				f := func(int) bool {
					assert.Fail(t, "predicate should not have been called")

					return res
				}

				assert.False(t, r.IsOkAnd(f))
			})
		}
	})
}

func TestResult_IsErr(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		val := fake.Int()
		r := newOk[int](val)

		assert.False(t, r.IsErr())
	})

	t.Run("Err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		r := newErr[any](err)

		assert.True(t, r.IsErr())
	})
}

func TestResult_IsErrAnd(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		val := fake.Int()
		r := newOk[int](val)

		for _, res := range []bool{true, false} {
			name := fmt.Sprintf("predicate returns %v", res)
			t.Run(name, func(t *testing.T) {
				f := func(error) bool {
					assert.Fail(t, "predicate should not have been called")

					return res
				}

				assert.False(t, r.IsErrAnd(f))
			})
		}
	})

	t.Run("Err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		r := newErr[int](err)

		for _, res := range []bool{true, false} {
			name := fmt.Sprintf("predicate returns %v", res)
			t.Run(name, func(t *testing.T) {
				called := false

				f := func(e error) bool {
					called = true

					assert.Equal(t, err, e)

					return res
				}

				assert.Equal(t, res, r.IsErrAnd(f))
				assert.True(t, called, "predicate should have been called")
			})
		}
	})
}

func TestResult_Expect(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		val := fake.Int()
		r := newOk[int](val)

		assert.Equal(t, val, r.Expect(fake.RandomStringWithLength(8)))
	})

	t.Run("Err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		r := newErr[any](err)

		msg := fake.RandomStringWithLength(8)
		expectedError := fmt.Errorf("%s: %w", msg, err)
		assert.PanicsWithError(t, expectedError.Error(), func() {
			r.Expect(msg)
		})
	})
}

func TestResult_ExpectErr(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		val := fake.Int()
		r := newOk[int](val)

		msg := fake.RandomStringWithLength(8)
		expectedError := fmt.Errorf("%s: %v", msg, val)
		assert.PanicsWithError(t, expectedError.Error(), func() {
			_ = r.ExpectErr(msg)
		})
	})

	t.Run("Err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		r := newErr[any](err)

		assert.Equal(t, err, r.ExpectErr(fake.RandomStringWithLength(8)))
	})
}

func TestResult_Unwrap(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		val := fake.Int()
		r := newOk[int](val)

		assert.Equal(t, val, r.Unwrap())
	})

	t.Run("Err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		r := newErr[any](err)

		expected := fmt.Errorf("called `Result.Unwrap()` on an `Err` value: %w", err)
		assert.PanicsWithError(t, expected.Error(), func() {
			_ = r.Unwrap()
		})
	})
}

func TestResult_UnwrapOr(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		val := fake.Int()
		r := newOk[int](val)

		assert.Equal(t, val, r.UnwrapOr(fake.Int()))
	})

	t.Run("Err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		r := newErr[any](err)

		def := fake.Int()

		assert.Equal(t, def, r.UnwrapOr(def))
	})
}

func TestResult_UnwrapOrElse(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		val := fake.Int()
		r := newOk[int](val)

		f := func() int {
			assert.Fail(t, "should not have been called")

			return fake.Int()
		}

		assert.Equal(t, val, r.UnwrapOrElse(f))
	})

	t.Run("Err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		r := newErr[int](err)

		def := fake.Int()
		f := func() int {
			return def
		}

		assert.Equal(t, def, r.UnwrapOrElse(f))
	})
}

func TestResult_UnwrapOrDefault(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		val := fake.Int()
		r := newOk[int](val)

		assert.Equal(t, val, r.UnwrapOrDefault())
	})

	t.Run("Err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))

		ri := newErr[int](err)
		assert.Equal(t, 0, ri.UnwrapOrDefault())

		rs := newErr[string](err)
		assert.Empty(t, rs.UnwrapOrDefault())
	})
}

func TestResult_UnwrapErr(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		val := fake.Int()
		r := newOk[int](val)

		expected := fmt.Errorf("called `Result.UnwrapErr()` on an `Ok` value: %v", val)
		assert.PanicsWithError(t, expected.Error(), func() {
			_ = r.UnwrapErr()
		})
	})

	t.Run("Err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		r := newErr[any](err)

		assert.Equal(t, err, r.UnwrapErr())
	})
}

func TestResult_Inspect(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		val := fake.Int()
		r := newOk[int](val)

		called := false
		p := func(v *int) {
			called = true

			assert.Equal(t, val, *v)
		}

		res := r.Inspect(p)

		assert.Same(t, res, r)
		assert.True(t, called, "predicate should have been called")
	})

	t.Run("Err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		r := newErr[int](err)

		p := func(*int) {
			assert.Fail(t, "predicate should not have been called")
		}

		res := r.Inspect(p)

		assert.Same(t, res, r)
	})
}

func TestResult_InspectErr(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		val := fake.Int()
		r := newOk[int](val)

		p := func(error) {
			assert.Fail(t, "predicate should not have been called")
		}

		res := r.InspectErr(p)

		assert.Same(t, res, r)
	})

	t.Run("Err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		r := newErr[int](err)

		called := false
		p := func(e error) {
			called = true

			assert.Equal(t, err, e)
		}

		res := r.InspectErr(p)

		assert.Same(t, res, r)
		assert.True(t, called, "predicate should have been called")
	})
}

func TestResult_AsOptionValue(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		val := fake.Int()
		r := newOk[int](val)

		assert.Equal(t, Some(val), r.AsOptionValue())
	})

	t.Run("Err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		r := newErr[int](err)

		assert.Equal(t, None[int](), r.AsOptionValue())
	})
}

func TestResult_AsOptionErr(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		val := fake.Int()
		r := newOk[int](val)

		assert.Equal(t, None[error](), r.AsOptionErr())
	})

	t.Run("Err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		r := newErr[int](err)

		assert.Equal(t, Some(err), r.AsOptionErr())
	})
}

func TestResult_WrapErr(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		val := fake.Int()
		r := newOk[int](val)

		assert.Same(t, r, r.WrapErr(fake.RandomStringWithLength(8)))
	})

	t.Run("Err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		r := newErr[int](err)

		msg := fake.RandomStringWithLength(8)
		expected := newErr[int](fmt.Errorf("%s: %w", msg, err))
		assert.Equal(t, expected, r.WrapErr(msg))
	})
}

func TestResult_Expand(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		val := fake.Int()
		r := newOk[int](val)

		v, err := r.Expand()
		require.NoError(t, err)
		assert.Equal(t, val, v)
	})

	t.Run("Err", func(t *testing.T) {
		sourceErr := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
		r := newErr[int](sourceErr)

		v, err := r.Expand()
		require.ErrorContains(t, err, sourceErr.Error())
		assert.Zero(t, v)
	})
}

func TestResult_String(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		value := fake.Int()
		r := newOk[int](value)

		expected := fmt.Sprintf("Ok(%v)", value)

		assert.Equal(t, expected, r.String())
	})

	t.Run("Err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		r := newErr[int](err)

		expected := fmt.Sprintf("Err(%v)", err)
		assert.Equal(t, expected, r.String())
	})
}

func newOk[T any](val T) *result[T] {
	r, ok := Ok(val).(*result[T])
	if !ok {
		panic("expected *result[T]")
	}

	return r
}

func newErr[T any](err error) *result[T] {
	r, ok := Err[T](err).(*result[T])
	if !ok {
		panic("expected *result[T]")
	}

	return r
}
