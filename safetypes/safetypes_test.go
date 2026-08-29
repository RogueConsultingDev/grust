package st

import (
	"github.com/jaswdr/faker/v2"
)

var fake = faker.New()

func ptr[T any](v T) *T {
	return &v
}
