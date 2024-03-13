package utils

import "sync"

// WheelArray 是一个环形数组，用于存储最近的数据
type WheelArray[T any] struct {
	sync.Mutex
	data  []T
	index int
	size  int
}

func NewWheelArray[T any](size int) *WheelArray[T] {
	return &WheelArray[T]{
		data:  make([]T, size),
		index: 0,
		size:  size,
	}
}

func (wa *WheelArray[T]) Add(value T) {
	wa.Lock()
	defer wa.Unlock()

	wa.data[wa.index] = value
	wa.index = (wa.index + 1) % wa.size
}

func (wa *WheelArray[T]) Get() []T {
	wa.Lock()
	defer wa.Unlock()

	return append(wa.data[wa.index:], wa.data[:wa.index]...)
}
