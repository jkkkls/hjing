package cache

import (
	"testing"
	"time"

	"github.com/gookit/goutil/dump"
)

type Data struct {
	name string
	age  int
}

func TestCache(t *testing.T) {
	// cache := New[*Data](0)
	// cache.Set("1001", &Data{"guangbo", 12})
	// cache.Set("1002", &Data{"guangxia", 13})

	// dump.Println(cache.Get("1001"))
	// dump.Println(cache.Get("1002"))
	// dump.Println(cache.Get("1003"))

	cache := NewGCache[*Data](CacheList, 1024, 3*time.Second)
	cache.Set("1001", &Data{"guangbo", 12})
	cache.Set("1002", &Data{"guangxia", 13})
	cache.Set("1003", &Data{"guangjing", 14})

	time.Sleep(2 * time.Second)
	dump.Println(cache.Get("1001"))
	time.Sleep(2 * time.Second)

	dump.Println(cache.Get("1001"))
	dump.Println(cache.Get("1002"))
	dump.Println(cache.Get("1003"))
}
