package st

import (
	"time"

	"github.com/jaswdr/faker/v2"
)

var fake = faker.NewWithSeedInt64(time.Now().UnixNano())
