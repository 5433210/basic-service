package util

import (
	"math/rand"
	"time"
)

func RandomInt(max int, min int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return min + r.Intn(max-min)
}
