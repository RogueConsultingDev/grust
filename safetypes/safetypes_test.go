package st

import (
	"errors"
	"testing"
	"time"

	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
)

var fake = faker.NewWithSeedInt64(time.Now().UnixNano())

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
