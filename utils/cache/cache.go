package cache

import (
	"sync"
	"time"
)

type ICacheNode[T any] interface {
	Get() T
}

type ICache[T any] interface {
	Set(string, T) ICacheNode[T]
	Get(string) ICacheNode[T]
	GetOrStore(string, func() (T, error)) ICacheNode[T]
}

type CacheNode[T any] struct {
	k          string
	v          T
	pre, next  *CacheNode[T]
	expireTime int64
}

func (n *CacheNode[T]) Get() T {
	return n.v
}

type Cache[T any] struct {
	sync.RWMutex

	expireTime time.Duration

	m          map[string]*CacheNode[T]
	head, tail *CacheNode[T]
}

func New[T any](expireTime time.Duration) *Cache[T] {
	return &Cache[T]{expireTime: expireTime, m: make(map[string]*CacheNode[T])}
}

func (t *Cache[T]) checkExpire() {
	if t.head == nil {
		return
	}

	if time.Now().Unix() < t.head.expireTime {
		return
	}

	n := t.head

	if t.head.next == nil {
		t.head, t.tail = nil, nil
	} else {
		t.head = t.head.next
		t.head.pre = nil
		n.next = nil
	}

	delete(t.m, n.k)
}
func (t *Cache[T]) refreshNodeLock(n *CacheNode[T]) {
	t.Lock()
	defer t.Unlock()
	t.refreshNode(n)
}

func (t *Cache[T]) refreshNode(n *CacheNode[T]) {
	if t.expireTime <= 0 {
		return
	}
	n.expireTime = time.Now().Add(t.expireTime).Unix()

	if n.next != nil && n.pre != nil {
		n.pre.next, n.next.pre = n.next, n.pre

		//放到尾部
		n.next = nil
		t.tail.next, n.pre = n, t.tail
		t.tail = n
	} else if n.next != nil && n.pre == nil {
		n.next.pre = nil
		t.head = n.next
		n.next = nil
		//放到尾部
		t.tail.next, n.pre = n, t.tail
		t.tail = n
	}
	t.checkExpire()
}

func (t *Cache[T]) Set(key string, value T) ICacheNode[T] {
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
		if t.tail == nil {
			t.head, t.tail = n, n
		} else {
			t.tail.next, n.pre = n, t.tail
			t.tail = n
		}
	}

	t.m[key] = n

	if t.expireTime > 0 {
		t.checkExpire()
	}
	return n
}

func (t *Cache[T]) Get(key string) ICacheNode[T] {
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

func (t *Cache[T]) GetOrStore(key string, f func() (T, error)) ICacheNode[T] {
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
