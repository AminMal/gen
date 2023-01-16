package gen

import (
	"math/rand"
	"time"
)

var random *rand.Rand

func Seed(seed int64) {
	random = rand.New(rand.NewSource(seed))
}

func init() {
	Seed(time.Now().UTC().UnixMilli())
}
