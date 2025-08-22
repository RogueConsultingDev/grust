package safetypes

import (
	"errors"
	"testing"
	"time"

	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"

	"github.com/mathieu-lemay/grust/safetypes/option"
	"github.com/mathieu-lemay/grust/safetypes/result"
)

var fake = faker.NewWithSeedInt64(time.Now().UnixNano())

func TestAsOkOr(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		value := fake.Int()
		s := option.Some(value)

		expected := result.Ok[int, error](value)

		assert.Equal(t, expected, AsOkOr(s, errors.New(fake.RandomStringWithLength(8))))
	})

	t.Run("none", func(t *testing.T) {
		n := option.Of(0)
		err := errors.New(fake.RandomStringWithLength(8))

		expected := result.Of(0, err)

		assert.Equal(t, expected, AsOkOr(n, err))
	})
}

func TestAsOkOrElse(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		value := fake.Int()
		s := option.Some(value)

		f := func() error {
			assert.Fail(t, "should not be called")

			return errors.New(fake.RandomStringWithLength(8))
		}

		expected := result.Ok[int, error](value)

		assert.Equal(t, expected, AsOkOrElse(s, f))
	})

	t.Run("none", func(t *testing.T) {
		n := option.Of(0)
		err := errors.New(fake.RandomStringWithLength(8))

		f := func() error {
			return err
		}

		expected := result.Of(0, err)

		assert.Equal(t, expected, AsOkOrElse(n, f))
	})
}

func TestAsOptionValue(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := result.Ok[int, error](value)

		expected := option.Some(value)

		assert.Equal(t, expected, AsOptionValue(o))
	})

	t.Run("err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		e := result.Of(0, err)

		expected := option.Of(0)

		assert.Equal(t, expected, AsOptionValue(e))
	})
}

func TestAsOptionErr(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := result.Ok[int, error](value)

		expected := option.Of((error)(nil))

		assert.Equal(t, expected, AsOptionErr(o))
	})

	t.Run("err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		e := result.Of(0, err)

		expected := option.Of(err)

		assert.Equal(t, expected, AsOptionErr(e))
	})
}
