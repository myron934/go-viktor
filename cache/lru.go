package cache

import (
	"container/list"
	"fmt"
	"reflect"
	"sync"
)

// LRUCache (Least Recently Used，最近最少使用) 淘汰策略缓存
type LRUCache[V any] struct {
	capacity    int
	cache       map[string]*list.Element
	list        *list.List
	mutex       sync.Mutex
	keyToString func(key any) string
}

type Entry[V any] struct {
	key   any
	value *V
}

// NewLRUCache 新建lru缓存(并发安全)
// capacity 最大容量, 超过会根据 最近最少使用 淘汰最后的数据
// V 为value类型
func NewLRUCache[V any](capacity int) *LRUCache[V] {
	return NewLRUCacheWithCustomKey[V](capacity, nil)
}

// NewLRUCacheWithCustomKey 新建lru缓存(并发安全). 如果key不是基础类型, 可以指定使用该方法创建, 指定key转化成字符串的方法
// capacity 最大容量, 超过会根据 最近最少使用 淘汰最后的数据
// customKeyFunc key转化成字符串的方法
// // V 为value类型
func NewLRUCacheWithCustomKey[V any](capacity int, customKeyFunc func(key any) string) *LRUCache[V] {
	if capacity < 0 {
		panic("capacity less than 0")
	}
	return &LRUCache[V]{
		capacity:    capacity,
		cache:       make(map[string]*list.Element),
		list:        list.New(),
		keyToString: customKeyFunc,
	}
}

// Get 获取数据
// Key可以支持基础类型: string | int | int8 | int16 | int32 | int64 | float32 | float64 | uint8 | uint16 | uint32 | uint64 | bool
// 以及实现了 String() string 接口的类
func (lru *LRUCache[V]) Get(key any) *V {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	keyStr := lru.stringKey(key)
	if elem, ok := lru.cache[keyStr]; ok {
		lru.list.MoveToFront(elem)
		return elem.Value.(*Entry[V]).value
	}
	return nil
}

// Put 设置缓存数据
// Key可以支持基础类型: string | int | int8 | int16 | int32 | int64 | float32 | float64 | uint8 | uint16 | uint32 | uint64 | bool
// 以及实现了 String() string 接口的类
func (lru *LRUCache[V]) Put(key any, value *V) {
	if lru.capacity == 0 {
		return
	}
	lru.mutex.Lock()
	defer lru.mutex.Unlock()
	strKey := lru.stringKey(key)

	if elem, ok := lru.cache[strKey]; ok {
		lru.list.MoveToFront(elem)
		elem.Value.(*Entry[V]).value = value
		return
	}

	if lru.list.Len() >= lru.capacity {
		// Remove the least recently used entry
		lru.deleteLast()
	}

	newEntry := &Entry[V]{key, value}
	newElem := lru.list.PushFront(newEntry)
	lru.cache[strKey] = newElem
}

// Clear 清空缓存
func (lru *LRUCache[V]) Clear() {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	lru.cache = make(map[string]*list.Element)
	lru.list.Init()
}

// Resize 重设缓存大小
func (lru *LRUCache[V]) Resize(capacity int) {
	if capacity < 0 {
		panic("capacity less than 0")
	}
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	for lru.list.Len() >= capacity {
		lru.deleteLast()
	}
	lru.capacity = capacity
}

func (lru *LRUCache[V]) Print() {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()
	printf("capacity=%v\n", lru.capacity)
	for key, elem := range lru.cache {
		val := elem.Value.(*Entry[V]).value
		printf("key=%v, val=%v\n", key, val)
	}
}

// deleteLast 删除最后一个元素
func (lru *LRUCache[V]) deleteLast() {
	lastElem := lru.list.Back()
	if lastElem == nil {
		return
	}
	delete(lru.cache, lru.stringKey(lastElem.Value.(*Entry[V]).key))
	lru.list.Remove(lastElem)
}

func (lru *LRUCache[V]) stringKey(key any) string {

	switch data := key.(type) {
	case fmt.Stringer:
		return data.String()
	case string:
		return data
	case int, int8, int16, int32, int64, float32, float64, uint8, uint16, uint32, uint64, bool:
		return fmt.Sprint(key)
	}
	if lru.keyToString == nil {
		panic("unknown key type and key to string function is nil")
	}
	return lru.keyToString(key)
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
