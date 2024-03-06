package utils

import (
	"runtime/debug"
	"sync"

	"github.com/panjf2000/ants/v2"
)

var (
	defaultAntsPool     *ants.Pool
	defaultAntsPoolOnce sync.Once
	DefaultErrCall      func(string, ...interface{})
)

func getPool() *ants.Pool {
	if defaultAntsPool == nil {
		defaultAntsPoolOnce.Do(func() {
			defaultAntsPool, _ = ants.NewPool(ants.DefaultAntsPoolSize, ants.WithPanicHandler(func(err interface{}) {
				Error("错误信息", "err", err)
				Error("错误调用堆栈信息", "stack", string(debug.Stack()))

				if DefaultErrCall != nil {
					DefaultErrCall("程序出错", "err", err, "stack", string(debug.Stack()))
				}
			}))
		})
	}

	return defaultAntsPool
}

func SetDefaultErrCall(f func(string, ...interface{})) {
	DefaultErrCall = f
}

// Submit submits a task to pool.
func Go(task func()) error {
	return getPool().Submit(task)
}

// // Running returns the number of the currently running goroutines.
// func Running() int {
// 	return getPool().Running()
// }

// // Cap returns the capacity of this default pool.
// func Cap() int {
// 	return getPool().Cap()
// }

// // Free returns the available goroutines to work.
// func Free() int {
// 	return getPool().Free()
// }

// // Release Closes the default pool.
// func Release() {
// 	getPool().Release()
// }

// // Reboot reboots the default pool.
// func Reboot() {
// 	getPool().Reboot()
// }
