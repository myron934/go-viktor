package cache

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	viktor "github.com/myron934/go-viktor"
)

func TestNewLoadingCache(t *testing.T) {
	ctx := context.Background()
	c := NewLoadingCache[string, int](WithCapacity[string, int](3))
	printf("key=a, value=%v\n", c.MustGet(ctx, "a"))

	c = NewLoadingCache[string, int](
		WithCapacity[string, int](3),
		WithKeyEncoder[string, int](func(k string) string { return k + "*" }),
		WithGetterFunc[string, int](func(key string) (*int, error) {
			println("refresh key " + key)
			return viktor.Ptr(1), nil
		}),
	)
	printf("key=a, value=%v\n", c.MustGet(ctx, "a"))
}

func TestNewLoadingCache2(t *testing.T) {
	c := NewLoadingCache[int, int](
		WithCapacity[int, int](10000000),
		WithExpireAfterWrite[int, int](time.Second*5),
		WithKeyEncoder[int, int](func(k int) string { return fmt.Sprint(k) }),
		WithGetterFunc[int, int](func(key int) (*int, error) {
			//fmt.Printf("refresh key %d \n", key)
			return viktor.Ptr(key), nil
		}),
	)
	for i := 0; i < 5; i++ {
		go func() {
			for {
				now := time.Now()
				c.MustGet(context.Background(), rand.Intn(10000000))
				cost := time.Now().Sub(now)
				if cost.Milliseconds() > 20 {
					fmt.Printf("get cache, cost: %d ms\n", time.Now().Sub(now).Milliseconds())

				}
				//printf("key=a, value=%v\n", c.MustGet(context.Background(), rand.Int()))
			}
		}()
	}
	for {
		time.Sleep(time.Second * 3)
		fmt.Printf("size:%d\n", c.Size())
	}

}
