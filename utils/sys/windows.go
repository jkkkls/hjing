// +build windows

package sys

import (
	"golang.org/x/sys/windows"
)

func GetThreadId() uint32 {
	return windows.GetCurrentThreadId()
}
