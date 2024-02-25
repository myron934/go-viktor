package cache

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	viktor "github.com/myron934/go-viktor"
)

func TestNewLRUCache(t *testing.T) {
	cache1 := NewLRUCache(WithCapacity[string, int](100))
	cache1.Put("false", viktor.Ptr(0))
	cache1.Put("1", viktor.Ptr(1))
	cache1.Put("k2", viktor.Ptr(2))
	cache1.Put("k3", viktor.Ptr(3))
	cache1.Put("k4", viktor.Ptr(4))
	fmt.Println("=======print all")
	cache1.Print()
	cache1.Put("k5", viktor.Ptr(5))
	fmt.Println("=======get value")
	printf("key=false, value=%v\n", cache1.MustGet("false"))
	printf("key=1, value=%v\n", cache1.MustGet("1"))
	printf("key=k2, value=%v\n", cache1.MustGet("k2"))
	printf("key=k3, value=%v\n", cache1.MustGet("k3"))
	printf("key=k4, value=%v\n", cache1.MustGet("k4"))
	printf("key=k5, value=%v\n", cache1.MustGet("k5"))

	fmt.Println("=======print all")
	cache1.Print()
}

func TestNewLRUCacheWithCustomKeyEncoder(t *testing.T) {
	type KeyStruct struct {
		a string
	}
	cache1 := NewLRUCache(
		WithCapacity[KeyStruct, int](100),
		WithKeyEncoder[KeyStruct, int](func(k KeyStruct) string {
			return k.a
		}))
	cache1.Put(KeyStruct{a: "1"}, viktor.Ptr(1))
	cache1.Put(KeyStruct{a: "2"}, viktor.Ptr(2))
	cache1.Put(KeyStruct{a: "3"}, viktor.Ptr(3))
	cache1.Put(KeyStruct{a: "4"}, viktor.Ptr(4))
	cache1.Put(KeyStruct{a: "5"}, viktor.Ptr(5))
	fmt.Println("=======print all")
	cache1.Print()
	cache1.Put(KeyStruct{a: "6"}, viktor.Ptr(6))

	fmt.Println("=======get value")
	printf("key=k1, value=%v\n", cache1.MustGet(KeyStruct{a: "1"}))
	printf("key=k2, value=%v\n", cache1.MustGet(KeyStruct{a: "2"}))
	printf("key=k3, value=%v\n", cache1.MustGet(KeyStruct{a: "3"}))
	printf("key=k4, value=%v\n", cache1.MustGet(KeyStruct{a: "4"}))
	printf("key=k5, value=%v\n", cache1.MustGet(KeyStruct{a: "5"}))
	printf("key=k6, value=%v\n", cache1.MustGet(KeyStruct{a: "6"}))

	fmt.Println("=======print all")
	cache1.Print()
}

func TestParallelRunLRU(t *testing.T) {
	cache := NewLRUCache[int, int](WithCapacity[int, int](20))
	go func() {
		for {
			time.Sleep(time.Millisecond * 10)
			key := rand.Intn(50)
			cache.MustGet(key)
			//printf("get key=%d, val=%v\n", key, cache.MustGet(key))
		}
	}()
	go func() {
		for {
			time.Sleep(time.Millisecond * 10)
			key := rand.Intn(50)
			cache.Put(key, &key)
			//printf("set key=%d, val=%d\n", key, key)
		}
	}()
	go func() {
		for {
			time.Sleep(time.Millisecond * 10)
			//c := rand.Intn(10) + 10
			//cache.Resize(c)
			//printf("reset capacity=%d\n", c)
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
