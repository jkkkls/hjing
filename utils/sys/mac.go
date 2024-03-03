//go:build darwin
// +build darwin

package sys

import (
	"math/rand/v2"
)

func GetThreadId() uint32 {
	return uint32(rand.IntN(1000))
}
