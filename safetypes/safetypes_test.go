package st

import (
	"errors"
	"testing"
	"time"

	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
)

var fake = faker.NewWithSeedInt64(time.Now().UnixNano())

func TestAsOkOr(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		value := fake.Int()
		s := Some(value)

		expected := Ok[int](value)

		assert.Equal(t, expected, AsOkOr(s, errors.New(fake.RandomStringWithLength(8))))
	})

	t.Run("none", func(t *testing.T) {
		n := OptionOf(0)
		err := errors.New(fake.RandomStringWithLength(8))

		expected := Err[int](err)

		assert.Equal(t, expected, AsOkOr(n, err))
	})
}

func TestAsOkOrElse(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		value := fake.Int()
		s := Some(value)

		f := func() error {
			assert.Fail(t, "should not be called")

			return errors.New(fake.RandomStringWithLength(8))
		}

		expected := Ok[int](value)

		assert.Equal(t, expected, AsOkOrElse(s, f))
	})

	t.Run("none", func(t *testing.T) {
		n := OptionOf(0)
		err := errors.New(fake.RandomStringWithLength(8))

		f := func() error {
			return err
		}

		expected := Err[int](err)

		assert.Equal(t, expected, AsOkOrElse(n, f))
	})
}

func TestAsOptionValue(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := Ok[int](value)

		assert.Equal(t, Some(value), AsOptionValue(o))
	})

	t.Run("err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		e := ResultOf(0, err)

		expected := None[int]()

		assert.Equal(t, expected, AsOptionValue(e))
	})
}

func TestAsOptionErr(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := Ok[int](value)

		expected := None[error]()

		assert.Equal(t, expected, AsOptionErr(o))
	})

	t.Run("err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		e := ResultOf(0, err)

		expected := Some(err)

		assert.Equal(t, expected, AsOptionErr(e))
	})
}
