package result

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOK_ReturnsNewOkResult(t *testing.T) {
	value := fake.Lorem().Word()

	res := Ok[string, error](value)

	expected := ok[string, error]{
		val: value,
	}
	assert.Equal(t, expected, res)
}

func TestOk_IsOk(t *testing.T) {
	o := Ok[int, error](fake.Int())

	assert.True(t, o.IsOk())
}

func TestOk_IsOkAnd(t *testing.T) {
	value := fake.Int()
	o := Ok[int, error](value)

	for _, res := range []bool{true, false} {
		name := fmt.Sprintf("predicate returns %v", res)
		t.Run(name, func(t *testing.T) {
			called := false

			f := func(v int) bool {
				called = true

				assert.Equal(t, value, v)

				return res
			}

			assert.Equal(t, res, o.IsOkAnd(f))
			assert.True(t, called, "predicate should have been called")
		})
	}
}

func TestOk_IsErr(t *testing.T) {
	o := Ok[int, error](fake.Int())

	assert.False(t, o.IsErr())
}

func TestOk_IsErrAnd(t *testing.T) {
	o := Ok[int, error](fake.Int())

	for _, res := range []bool{true, false} {
		name := fmt.Sprintf("predicate returns %v", res)
		t.Run(name, func(t *testing.T) {
			f := func(_ error) bool {
				assert.Fail(t, "predicate should not have been called")

				return res
			}

			assert.False(t, o.IsErrAnd(f))
		})
	}
}

func TestOk_Expect(t *testing.T) {
	value := fake.Int()
	o := Ok[int, error](value)

	msg := fake.RandomStringWithLength(8)
	assert.Equal(t, value, o.Expect(msg))
}

func TestOk_ExpectErr(t *testing.T) {
	value := fake.Int()
	o := Ok[int, error](value)

	msg := fake.RandomStringWithLength(8)
	expectedError := fmt.Errorf("%s: %v", msg, value)
	assert.PanicsWithError(t, expectedError.Error(), func() {
		_ = o.ExpectErr(msg)
	})
}

func TestOk_Inspect(t *testing.T) {
	value := fake.Int()
	o := Ok[int, error](value)

	called := false
	f := func(v *int) {
		called = true

		assert.Equal(t, value, *v)
	}

	assert.Equal(t, o, o.Inspect(f))
	assert.True(t, called, "inspector should have been called")
}

func TestOk_InspectErr(t *testing.T) {
	value := fake.Int()
	o := Ok[int, error](value)

	f := func(_ *error) {
		assert.Fail(t, "inspector should not be called")
	}

	assert.Equal(t, o, o.InspectErr(f))
}

func TestOk_Unwrap(t *testing.T) {
	value := fake.Int()
	s := Ok[int, error](value)

	assert.Equal(t, value, s.Unwrap())
}

func TestOk_UnwrapOr(t *testing.T) {
	value := fake.Int()
	s := Ok[int, error](value)

	def := fake.Int()

	assert.Equal(t, value, s.UnwrapOr(def))
}

func TestOk_UnwrapOrElse(t *testing.T) {
	value := fake.Int()
	s := Ok[int, error](value)

	def := fake.Int()
	f := func() int {
		assert.Fail(t, "should not be called")

		return def
	}

	assert.Equal(t, value, s.UnwrapOrElse(f))
}

func TestOk_UnwrapOrDefault(t *testing.T) {
	valStr := fake.RandomStringWithLength(8)
	valInt := fake.IntBetween(1, 100)
	valFloat := fake.Float64(2, 1, 100)

	assert.Equal(t, valStr, Ok[string, error](valStr).UnwrapOrDefault())
	assert.Equal(t, valInt, Ok[int, error](valInt).UnwrapOrDefault())
	assert.Equal( //nolint:testifylint  // Value should be _exactly_ equal
		t,
		valFloat,
		Ok[float64, error](valFloat).UnwrapOrDefault(),
	)
}

func TestOk_UnwrapErr(t *testing.T) {
	value := fake.Int()
	o := Ok[int, error](value)

	expected := fmt.Errorf("called `Result.UnwrapErr()` on an `Ok` value: %v", value)
	assert.PanicsWithError(t, expected.Error(), func() {
		_ = o.UnwrapErr()
	})
}

func TestOk_String(t *testing.T) {
	value := fake.Int()
	o := Ok[int, error](value)

	expected := fmt.Sprintf("Ok(%v)", value)
	assert.Equal(t, expected, o.String())
}
