package utils

import (
	"sync"
)

//
//id生成器逻辑
//

// IDGen id生成器
type IDGen struct {
	Time  uint64
	GenID uint64 //生成器实例id
	Index uint64
	Mutex sync.Mutex
}

// NewIDGen 创建一个id生成器
func NewIDGen(id uint64) *IDGen {
	return &IDGen{
		GenID: id,
	}
}

// NewID 生成一个id
func (gen *IDGen) NewID(moduleID uint64) uint64 {
	now := NowNano() / 1000000

	gen.Mutex.Lock()
	defer gen.Mutex.Unlock()
	if now != gen.Time {
		gen.Time = now
		gen.Index = 0
	}

	gen.Index++
	if gen.Index&0xFFF == 0 {
		panic("id generate error")
	}

	return (gen.Time&0x1FFFFFFFFFF)<<20 | (gen.GenID&0xFF)<<16 | (moduleID&0xFF)<<12 | gen.Index&0xFFF
}

// GetTimeFromId 从id拿到时间
func GetTimeFromId(id uint64) uint64 {
	return ((id >> 20) & 0x1FFFFFFFFFF) / 1000
}
