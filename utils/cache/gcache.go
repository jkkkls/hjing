package cache

import "time"

const (
	// offset64 FNVa offset basis. See https://en.wikipedia.org/wiki/Fowler–Noll–Vo_hash_function#FNV-1a_hash
	offset64 = 14695981039346656037
	// prime64 FNVa prime value. See https://en.wikipedia.org/wiki/Fowler–Noll–Vo_hash_function#FNV-1a_hash
	prime64 = 1099511628211
)

const (
	CacheList = 0
	CacheRand = 1
)

// Sum64 gets the string and returns its uint64 hash value.
func fnv64a(key string) uint64 {
	var hash uint64 = offset64
	for i := 0; i < len(key); i++ {
		hash ^= uint64(key[i])
		hash *= prime64
	}

	return hash
}

type GCache[T any] struct {
	shards    []ICache[T]
	shardMask uint64
}

// NewGCache 分片cache
// t 过期机制类型，CacheList-双向链表，每次删除链头过期数据，CacheRand-随机模式，每次随机找到四条记录检查过期
// n must be a power of 2
func NewGCache[T any](t int, n int, expireTime time.Duration) *GCache[T] {
	if n == 0 {
		n = 1024
	}

	c := &GCache[T]{
		shards:    make([]ICache[T], n),
		shardMask: uint64(n) - 1,
	}
	if t == CacheList {
		for i := 0; i < n; i++ {
			c.shards[i] = New[T](expireTime)
		}
	} else {
		for i := 0; i < n; i++ {
			c.shards[i] = NewRCache[T](expireTime)
		}
	}
	return c
}

func (t *GCache[T]) getShard(hashedKey uint64) ICache[T] {
	return t.shards[hashedKey&t.shardMask]
}

func (t *GCache[T]) Set(key string, value T) ICacheNode[T] {
	hashedKey := fnv64a(key)
	shard := t.getShard(hashedKey)
	return shard.Set(key, value)
}

func (t *GCache[T]) Get(key string) ICacheNode[T] {
	hashedKey := fnv64a(key)
	shard := t.getShard(hashedKey)
	return shard.Get(key)
}

func (t *GCache[T]) GetOrStore(key string, f func() (T, error)) ICacheNode[T] {
	hashedKey := fnv64a(key)
	shard := t.getShard(hashedKey)
	return shard.GetOrStore(key, f)
}
