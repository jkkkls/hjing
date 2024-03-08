package cache

import (
	"sync"
	"time"
)

type RCacheNode[T any] struct {
	v          T
	expireTime int64
}

func (n *RCacheNode[T]) Get() T {
	return n.v
}

type RCache[T any] struct {
	sync.RWMutex

	expireTime time.Duration

	m map[string]*CacheNode[T]
}

func NewRCache[T any](expireTime time.Duration) *RCache[T] {
	return &RCache[T]{expireTime: expireTime, m: make(map[string]*CacheNode[T])}
}

func (t *RCache[T]) checkExpire() {
	i := 4
	for k, v := range t.m {
		if time.Now().Unix() >= v.expireTime {
			delete(t.m, k)
		}
		i--
		if i == 0 {
			break
		}
	}
}

func (t *RCache[T]) refreshNodeLock(n *CacheNode[T]) {
	t.Lock()
	defer t.Unlock()
	t.refreshNode(n)
}

func (t *RCache[T]) refreshNode(n *CacheNode[T]) {
	if t.expireTime <= 0 {
		return
	}
	n.expireTime = time.Now().Add(t.expireTime).Unix()
	t.checkExpire()
}

func (t *RCache[T]) Set(key string, value T) ICacheNode[T] {
	t.Lock()
	defer t.Unlock()

	//如果存在
	if n, ok := t.m[key]; ok {
		n.v = value
		t.refreshNode(n)

		return n
	}

	n := &CacheNode[T]{
		k: key,
		v: value,
	}
	if t.expireTime > 0 {
		n.expireTime = time.Now().Add(t.expireTime).Unix()
	}

	t.m[key] = n

	if t.expireTime > 0 {
		t.checkExpire()
	}
	return n
}

func (t *RCache[T]) Get(key string) ICacheNode[T] {
	t.RLock()

	if n, ok := t.m[key]; ok {
		t.RUnlock()
		t.refreshNodeLock(n)
		return n
	}
	t.RUnlock()

	if t.expireTime > 0 {
		t.checkExpire()
	}
	return nil
}

func (t *RCache[T]) GetOrStore(key string, f func() (T, error)) ICacheNode[T] {
	t.RLock()

	if n, ok := t.m[key]; ok {
		t.RUnlock()
		t.refreshNodeLock(n)
		return n
	}
	t.RUnlock()

	value, err := f()
	if err != nil {
		return nil
	}

	n := t.Set(key, value)

	return n
}
