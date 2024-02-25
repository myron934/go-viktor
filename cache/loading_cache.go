package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type (
	Serializer[T1, T2 any]   func(context.Context, *T1) *T2
	Deserializer[T1, T2 any] func(context.Context, *T1) *T2
	ExpireFunc[T any]        func(context.Context, *T) time.Duration
	LoadFunc[T any]          func(context.Context, []any) (map[any]*T, error)
)

type ICache[K, V any] interface {
	Get(ctx context.Context, key K) (*V, error)
	Put(ctx context.Context, key K, val *V) error
	GetAll(ctx context.Context, keys []K) map[any]*V
	refresh(ctx context.Context, key any) (*V, error)
	Remove(ctx context.Context, keys ...any) error
}

type LoadingItem[V any] struct {
	expire time.Time
	value  *V
}

type LoadingCache[K, V any] struct {
	lruCache      *LRUCache[K, LoadingItem[V]]
	mutex         sync.Mutex
	conf          *Config[K, V]
	lastClearTime time.Time
}

func NewLoadingCache[K, V any](opts ...Option[K, V]) *LoadingCache[K, V] {
	c := &LoadingCache[K, V]{
		lruCache: NewLRUCache[K, LoadingItem[V]](),
		conf:     NewDefaultConf[K, V](),
	}
	for _, opt := range opts {
		c.conf = opt(c.conf)
	}
	c.lruCache = NewLRUCache[K, LoadingItem[V]](
		WithCapacity[K, LoadingItem[V]](c.conf.capacity),
		WithKeyEncoder[K, LoadingItem[V]](c.conf.keyToString),
	)
	return c
}

func (c *LoadingCache[K, V]) Get(_ context.Context, key K) (*V, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	val, err := c.lruCache.Get(key)
	now := time.Now()
	if err == nil && val != nil && now.Before(val.expire) {
		return val.value, nil
	}
	if newVal, err := c.refresh(key); err == nil {
		return newVal, nil
	}
	return nil, ErrorKeyNotFound
}

func (c *LoadingCache[K, V]) MustGet(ctx context.Context, key K) *V {
	val, _ := c.Get(ctx, key)
	return val
}

func (c *LoadingCache[K, V]) Refresh(_ context.Context, key K) error {
	_, err := c.refresh(key)
	return err
}

func (c *LoadingCache[K, V]) refresh(key K) (*V, error) {
	if c.conf.getterFunc == nil {
		return nil, ErrorKeyNotFound
	}
	val, err := c.conf.getterFunc(key)
	if err != nil {
		return nil, err
	}
	if err = c.put(key, val); err != nil {
		return nil, err
	}
	return val, nil
}

func (c *LoadingCache[K, V]) Put(_ context.Context, key K, val *V) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.put(key, val)
}

func (c *LoadingCache[K, V]) put(key K, val *V) error {
	item := &LoadingItem[V]{
		expire: time.Now().Add(c.conf.expireAfterWrite),
		value:  val,
	}
	//if c.lruCache.IsFull() {
	//	c.clearExpireItem(false)
	//}
	c.lruCache.Put(key, item)
	return nil
}

func (c *LoadingCache[K, V]) Size() int {
	return c.lruCache.Size()
}

func (c *LoadingCache[K, V]) clearExpireItem(force bool) {
	now := time.Now()
	size := c.Size()
	if !force && c.lastClearTime.Add(c.conf.minClearInterval).After(now) {
		return
	}
	c.lruCache.RemoveIf(func(key K, val *LoadingItem[V]) bool {
		return val.expire.Before(now)
	})
	c.lastClearTime = now
	fmt.Printf("clear expire item, num: %d, cost: %d ms\n", size-c.Size(), time.Now().Sub(now).Milliseconds())
}

func (c *LoadingCache[K, V]) stringKey(key any) string {

	if c.conf.keyToString != nil {
		k, ok := key.(K)
		if !ok {
			panic("key type error " + fmt.Sprint(key))
		}
		return c.conf.keyToString(k)
	}
	switch data := key.(type) {
	case string:
		return data
	case int, int8, int16, int32, int64, float32, float64, uint8, uint16, uint32, uint64, bool:
		return fmt.Sprint(key)
	case fmt.Stringer:
		return data.String()
	default:
		panic("unsupported key type " + fmt.Sprint(key))
	}
}
