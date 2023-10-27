package bytes_cache

import (
	"sync"
)

const (
	n = 64
)

var (
	cache [n]sync.Pool
)

func init() {
	for i := 0; i < n; i++ {
		size := 1 << i
		cache[i].New = func() interface{} {
			return make([]byte, 0, size)
		}
	}
}

func Get(l int, capa ...int) []byte {
	size := l
	if len(capa) > 0 {
		if l > capa[0] {
			panic("capacity less to length ")
		}
		size = capa[0]
	}

	return cache[getIndex(size)].Get().([]byte)[:l]
}

func Put(b []byte) {
	capa := cap(b)
	if (capa & -capa) != capa {
		return
	}
	cache[getIndex(capa)].Put(b)
}

func getIndex(size int) int {
	if size == 0 {
		return 0
	}
	tp := (size & -size) == size
	var index int
	for {
		if size == 0 {
			break
		}
		size >>= 1
		index++
	}
	if tp {
		index--
	}

	return index
}
