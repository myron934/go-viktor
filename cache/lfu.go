package cache

import (
	"container/heap"
	"fmt"
	"sync"
)

// LFUCache (Least Recently Used，最近最少使用) 淘汰策略缓存
type LFUCache[V any] struct {
	capacity    int
	cache       map[string]*LFUItem
	pq          PriorityQueue
	mutex       sync.Mutex
	keyToString func(key any) string
}

type LFUItem struct {
	key       any
	value     any
	frequency int
	index     int
}

// NewLFUCache 新建lfu缓存(并发安全)
// capacity 最大容量, 超过会根据 最近最少使用 淘汰最后的数据
// V 为value类型
func NewLFUCache[V any](capacity int) *LFUCache[V] {
	return NewLFUCacheWithCustomKey[V](capacity, nil)
}

// NewLFUCacheWithCustomKey 新建lfu缓存(并发安全). 如果key不是基础类型, 可以指定使用该方法创建, 指定key转化成字符串的方法
// capacity 最大容量, 超过会根据 最近最少使用 淘汰最后的数据
// customKeyFunc key转化成字符串的方法
// // V 为value类型
func NewLFUCacheWithCustomKey[V any](capacity int, customKeyFunc func(key any) string) *LFUCache[V] {
	if capacity < 0 {
		panic("capacity less than 0")
	}
	return &LFUCache[V]{
		capacity:    capacity,
		cache:       make(map[string]*LFUItem),
		pq:          make(PriorityQueue, 0),
		keyToString: customKeyFunc,
	}
}

// Get 获取数据
// Key可以支持基础类型: string | int | int8 | int16 | int32 | int64 | float32 | float64 | uint8 | uint16 | uint32 | uint64 | bool
// 以及实现了 String() string 接口的类
func (lfu *LFUCache[V]) Get(key any) *V {
	lfu.mutex.Lock()
	defer lfu.mutex.Unlock()

	keyStr := lfu.stringKey(key)
	if item, ok := lfu.cache[keyStr]; ok {
		lfu.updateFrequency(item)
		return item.value.(*V)
	}
	return nil
}

// Put 设置缓存数据
// Key可以支持基础类型: string | int | int8 | int16 | int32 | int64 | float32 | float64 | uint8 | uint16 | uint32 | uint64 | bool
// 以及实现了 String() string 接口的类
func (lfu *LFUCache[V]) Put(key any, value *V) {
	if lfu.capacity == 0 {
		return
	}
	lfu.mutex.Lock()
	defer lfu.mutex.Unlock()
	strKey := lfu.stringKey(key)

	if item, ok := lfu.cache[strKey]; ok {
		item.value = value
		lfu.updateFrequency(item)
		return
	}

	if len(lfu.cache) >= lfu.capacity {
		// Remove the least frequently used item
		lfu.deleteLeastUsed()
	}

	newItem := &LFUItem{
		key:       key,
		value:     value,
		frequency: 1,
	}
	heap.Push(&lfu.pq, newItem)
	lfu.cache[strKey] = newItem
}

// Clear 清空缓存
func (lfu *LFUCache[V]) Clear() {
	lfu.mutex.Lock()
	defer lfu.mutex.Unlock()

	lfu.cache = make(map[string]*LFUItem)
	lfu.pq = make(PriorityQueue, 0)
}

// Resize 重设缓存大小
func (lfu *LFUCache[V]) Resize(capacity int) {
	if capacity < 0 {
		panic("capacity less than 0")
	}
	lfu.mutex.Lock()
	defer lfu.mutex.Unlock()

	for len(lfu.cache) >= capacity {
		lfu.deleteLeastUsed()
	}
	lfu.capacity = capacity
}

func (lfu *LFUCache[V]) Print() {
	lfu.mutex.Lock()
	defer lfu.mutex.Unlock()
	printf("capacity=%v\n", lfu.capacity)
	for key, item := range lfu.cache {
		printf("key=%v, val=%v, frequency=%v, index=%v\n", key, item.value, item.frequency, item.index)
	}
}

// deleteLeastUsed 删除优先级最低的一个元素
func (lfu *LFUCache[V]) deleteLeastUsed() {
	removedItem := heap.Pop(&lfu.pq).(*LFUItem)
	delete(lfu.cache, lfu.stringKey(removedItem.key))
}

func (lfu *LFUCache[V]) updateFrequency(item *LFUItem) {
	item.frequency++
	heap.Fix(&lfu.pq, item.index)
}

func (lfu *LFUCache[V]) stringKey(key any) string {

	switch data := key.(type) {
	case fmt.Stringer:
		return data.String()
	case string:
		return data
	case int, int8, int16, int32, int64, float32, float64, uint8, uint16, uint32, uint64, bool:
		return fmt.Sprint(key)
	}
	if lfu.keyToString == nil {
		panic("unknown key type and key to string function is nil")
	}
	return lfu.keyToString(key)
}

// PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*LFUItem

func (pq *PriorityQueue) Len() int { return len((*pq)) }

func (pq *PriorityQueue) Less(i, j int) bool {
	if (*pq)[i].frequency == (*pq)[j].frequency {
		// If frequencies are the same, prioritize by index
		return (*pq)[i].index < (*pq)[j].index
	}
	// Otherwise, prioritize by frequency (lower is better)
	return (*pq)[i].frequency < (*pq)[j].frequency
}

func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	(*pq)[i].index = i
	(*pq)[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*LFUItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
