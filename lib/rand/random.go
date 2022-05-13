package rand

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func RandInt(max, min int64) int64 {
	return rand.Int63n(max-min) + min
}

func RandIntWithSeed(seed, max, min int64) int64 {
	return rand.New(rand.NewSource(seed + time.Now().Unix())).Int63n(max-min) + min
}
