package result

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErr_ReturnsNewErrResult(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))

	res := Err[any, error](err)

	expected := errT[any, error]{
		err: err,
	}
	assert.Equal(t, expected, res)
}

func TestErr_IsOk(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err[any, error](err)

	assert.False(t, e.IsOk())
}

func TestErr_IsOkAnd(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err[any, error](err)

	for _, res := range []bool{true, false} {
		name := fmt.Sprintf("predicate returns %v", res)
		t.Run(name, func(t *testing.T) {
			f := func(_ any) bool {
				assert.Fail(t, "predicate should not have been called")

				return res
			}

			assert.False(t, e.IsOkAnd(f))
		})
	}
}

func TestErr_IsErr(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err[any, error](err)

	assert.True(t, e.IsErr())
}

func TestErr_IsErrAnd(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err[any, error](err)

	for _, res := range []bool{true, false} {
		name := fmt.Sprintf("predicate returns %v", res)
		t.Run(name, func(t *testing.T) {
			called := false

			f := func(e error) bool {
				called = true

				assert.Equal(t, err, e)

				return res
			}

			assert.Equal(t, res, e.IsErrAnd(f))
			assert.True(t, called, "predicate should have been called")
		})
	}
}

func TestErr_Expect(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err[any, error](err)

	msg := fake.RandomStringWithLength(8)
	expectedError := fmt.Errorf("%s: %w", msg, err)
	assert.PanicsWithError(t, expectedError.Error(), func() {
		e.Expect(msg)
	})
}

func TestErr_ExpectErr(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err[any, error](err)

	msg := fake.RandomStringWithLength(8)

	assert.Equal(t, err, e.ExpectErr(msg))
}

func TestErr_Inspect(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err[any, error](err)

	f := func(_ *any) {
		assert.Fail(t, "inspector should not be called")
	}

	assert.Equal(t, e, e.Inspect(f))
}

func TestErr_InspectErr(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err[any, error](err)

	called := false
	f := func(ep *error) {
		called = true

		assert.Equal(t, err, *ep)
	}

	assert.Equal(t, e, e.InspectErr(f))
	assert.True(t, called, "inspector should have been called")
}

func TestErr_Unwrap(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err[any, error](err)

	expected := fmt.Errorf("called `Result.Unwrap()` on an `Err` value: %w", err)
	assert.PanicsWithError(t, expected.Error(), func() {
		e.Unwrap()
	})
}

func TestErr_UnwrapOr(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err[any, error](err)

	def := fake.RandomStringWithLength(8)

	assert.Equal(t, def, e.UnwrapOr(def))
}

func TestErr_UnwrapOrElse(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err[any, error](err)

	def := fake.RandomStringWithLength(8)
	f := func() any {
		return def
	}

	assert.Equal(t, def, e.UnwrapOrElse(f))
}

func TestErr_UnwrapOrDefault(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	errStr := Err[string, error](err)
	errInt := Err[int, error](err)
	errFloat := Err[float64, error](err)

	assert.Zero( //nolint:testifylint  // Using assert.Zero for consistency
		t,
		errStr.UnwrapOrDefault(),
	)
	assert.Zero(t, errInt.UnwrapOrDefault())
	assert.Zero(t, errFloat.UnwrapOrDefault())
}

func TestErr_UnwrapErr(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err[any, error](err)

	assert.Equal(t, err, e.UnwrapErr())
}

func TestErr_String(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err[any, error](err)

	expected := fmt.Sprintf("Err(%v)", err)
	assert.Equal(t, expected, e.String())
}
