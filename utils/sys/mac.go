// +build darwin

package sys

import (
	"math/rand"
)

func GetThreadId() uint32 {
	return uint32(rand.Intn(1000))
}
