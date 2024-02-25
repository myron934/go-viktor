package cache

import (
	"container/list"
	"fmt"
	"reflect"
)

// LRUCache (Least Recently Used，最近最少使用) 淘汰策略缓存
type LRUCache[K, V any] struct {
	cache map[string]*list.Element
	list  *list.List
	//mutex       sync.Mutex
	conf *Config[K, V]
}

type Entry[K, V any] struct {
	key   K
	value *V
}

// NewLRUCache 新建lru缓存(并发安全)
// capacity 最大容量, 超过会根据 最近最少使用 淘汰最后的数据
// V 为value类型
func NewLRUCache[K, V any](opts ...Option[K, V]) *LRUCache[K, V] {
	c := &LRUCache[K, V]{
		cache: make(map[string]*list.Element),
		list:  list.New(),
		conf:  NewDefaultConf[K, V](),
	}
	for _, opt := range opts {
		c.conf = opt(c.conf)
	}

	return c
}

// Get 获取数据
// Key可以支持基础类型: string | int | int8 | int16 | int32 | int64 | float32 | float64 | uint8 | uint16 | uint32 | uint64 | bool
// 以及实现了 String() string 接口的类
func (lru *LRUCache[K, V]) Get(key K) (*V, error) {
	//lru.mutex.Lock()
	//defer lru.mutex.Unlock()

	keyStr := lru.stringKey(key)
	if elem, ok := lru.cache[keyStr]; ok {
		lru.list.MoveToFront(elem)
		return elem.Value.(*Entry[K, V]).value, nil
	}
	return nil, ErrorKeyNotFound
}

// MustGet 同 Get, 如果key不存在返回nil
func (lru *LRUCache[K, V]) MustGet(key K) *V {
	//lru.mutex.Lock()
	//defer lru.mutex.Unlock()

	keyStr := lru.stringKey(key)
	if elem, ok := lru.cache[keyStr]; ok {
		lru.list.MoveToFront(elem)
		return elem.Value.(*Entry[K, V]).value
	}
	return nil
}

// Put 设置缓存数据
// Key可以支持基础类型: string | int | int8 | int16 | int32 | int64 | float32 | float64 | uint8 | uint16 | uint32 | uint64 | bool
// 以及实现了 String() string 接口的类
func (lru *LRUCache[K, V]) Put(key K, value *V) {
	if lru.conf.capacity == 0 {
		return
	}
	//lru.mutex.Lock()
	//defer lru.mutex.Unlock()
	strKey := lru.stringKey(key)

	if elem, ok := lru.cache[strKey]; ok {
		lru.list.MoveToFront(elem)
		elem.Value.(*Entry[K, V]).value = value
		return
	}

	if lru.list.Len() >= lru.conf.capacity {
		// Remove the least recently used entry
		lru.deleteLast()
	}

	newEntry := &Entry[K, V]{key, value}
	newElem := lru.list.PushFront(newEntry)
	lru.cache[strKey] = newElem
}

// Clear 清空缓存
func (lru *LRUCache[K, V]) Clear() {
	//lru.mutex.Lock()
	//defer lru.mutex.Unlock()

	lru.cache = make(map[string]*list.Element)
	lru.list.Init()
}

// Size 获取当前元素数量
func (lru *LRUCache[K, V]) Size() int {
	return lru.list.Len()
}

func (lru *LRUCache[K, V]) IsFull() bool {
	return lru.Size() >= lru.conf.capacity
}

// Resize 重设缓存大小
func (lru *LRUCache[K, V]) Resize(capacity int) {
	if capacity < 0 {
		panic("capacity less than 0")
	}
	//lru.mutex.Lock()
	//defer lru.mutex.Unlock()

	for lru.list.Len() >= capacity {
		lru.deleteLast()
	}
	lru.conf.capacity = capacity
}

func (lru *LRUCache[K, V]) Print() {
	//lru.mutex.Lock()
	//defer lru.mutex.Unlock()
	printf("capacity=%v\n", lru.conf.capacity)
	for key, elem := range lru.cache {
		val := elem.Value.(*Entry[K, V]).value
		printf("key=%v, val=%v\n", key, val)
	}
}

// Remove 删除元素
func (lru *LRUCache[K, V]) Remove(key K) {
	strKey := lru.stringKey(key)
	elem, ok := lru.cache[strKey]
	if !ok {
		return
	}
	delete(lru.cache, strKey)
	lru.list.Remove(elem)
}

func (lru *LRUCache[K, V]) RemoveIf(condition func(K, *V) bool) {
	for _, elem := range lru.cache {
		entry := elem.Value.(*Entry[K, V])
		if condition(entry.key, entry.value) {
			lru.Remove(entry.key)
		}
	}
}

// deleteLast 删除最后一个元素
func (lru *LRUCache[K, V]) deleteLast() {
	lastElem := lru.list.Back()
	if lastElem == nil {
		return
	}
	lru.Remove(lastElem.Value.(*Entry[K, V]).key)
}

func (lru *LRUCache[K, V]) stringKey(key any) string {
	if lru.conf.keyToString != nil {
		k, ok := key.(K)
		if !ok {
			panic("key type error " + fmt.Sprint(key))
		}
		return lru.conf.keyToString(k)
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

func printf(format string, a ...any) {
	if len(a) == 0 {
		fmt.Printf(format, a...)
		return
	}
	temp := make([]any, 0, len(a))
	for _, item := range a {
		if item == nil {
			temp = append(temp, item)
			continue
		}
		itemValue := reflect.ValueOf(item)
		if itemValue.Kind() == reflect.Ptr {
			if itemValue.IsNil() {
				temp = append(temp, nil)
				continue
			}
			itemValue = itemValue.Elem()
			temp = append(temp, itemValue.Interface())
			continue
		}
		temp = append(temp, item)

	}
	fmt.Printf(format, temp...)

}
