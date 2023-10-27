// +build linux

package sys

import (
	"golang.org/x/sys/unix"
)

func GetThreadId() uint32 {
	return uint32(unix.Gettid())
}
