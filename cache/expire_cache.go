package cache

import (
	"context"
	"fmt"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type (
	Serializer[T1, T2 any]   func(context.Context, *T1) *T2
	Deserializer[T1, T2 any] func(context.Context, *T1) *T2
	ExpireFunc[T any]        func(context.Context, *T) time.Duration
	LoadFunc[T any]          func(context.Context, []any) (map[any]*T, error)
)

type ICache[T any] interface {
	Get(ctx context.Context, key any) (*T, error)
	Put(ctx context.Context, key any, v *T) error
	MultiGet(ctx context.Context, keys []any) map[any]*T
	//Reload(ctx context.Context, key any) (*T, error)
	Remove(ctx context.Context, keys ...any) error
}

type LoadingCache[V any] struct {
}

func (c *LoadingCache[V]) Get(ctx context.Context, key any) *V {
	gc := gocache.New(5*time.Minute, 10*time.Minute)
	val, exist := gc.Get("")
	println(val, exist)
	return nil
}

func (c *LoadingCache[V]) stringKey(key any) string {

	switch data := key.(type) {
	case fmt.Stringer:
		return data.String()
	case string:
		return data
	case int, int8, int16, int32, int64, float32, float64, uint8, uint16, uint32, uint64, bool:
		return fmt.Sprint(key)
	}
	//if lfu.keyToString == nil {
	//	panic("unknown key type and key to string function is nil")
	//}
	//return lfu.keyToString(key)
	return ""
}
