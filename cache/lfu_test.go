package cache

import (
	"math/rand"
	"testing"
	"time"

	viktor "github.com/myron934/go-viktor"
)

func TestNewLFUCache(t *testing.T) {
	cache := NewLFUCache[int](3)
	cache.Put("1", viktor.Ptr(1))
	cache.Put("2", viktor.Ptr(2))
	cache.Put("3", viktor.Ptr(3))
	cache.Print()
	cache.Get("1")
	cache.Print()
	cache.Get("2")
	cache.Get("2")
	cache.Get("2")
	cache.Get("2")
	cache.Get("2")
	cache.Print()
	cache.Put("4", viktor.Ptr(4))
	cache.Print()
}

func TestParallelRunLFU(t *testing.T) {
	cache := NewLFUCache[int](20)
	go func() {
		for {
			time.Sleep(time.Millisecond * 10)
			key := rand.Intn(50)
			cache.Get(key)
			printf("get key=%d, val=%v\n", key, cache.Get(key))
		}
	}()
	go func() {
		for {
			time.Sleep(time.Millisecond * 10)
			key := rand.Intn(50)
			cache.Put(key, &key)
			printf("set key=%d, val=%d\n", key, key)
		}
	}()
	go func() {
		for {
			time.Sleep(time.Millisecond * 10)
			c := rand.Intn(10) + 10
			cache.Resize(c)
			printf("reset capcity=%d\n", c)
		}
	}()
	go func() {
		for {
			time.Sleep(time.Millisecond * 10)
			cache.Print()
		}
	}()
	time.Sleep(time.Second * 3600)
}
