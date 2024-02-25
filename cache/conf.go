package cache

import "time"

type Config[K, V any] struct {
	capacity         int
	expireAfterWrite time.Duration // 过期时间
	clearInterval    time.Duration // 定时清理过期key的间隔
	minClearInterval time.Duration // 为了防止缓存满了以后频繁触发清理, 定义最小触发间隔, 该时间内如果已经清理过,则不再清理
	keyToString      func(key K) string
	getterFunc       func(key K) (*V, error) //缓存不存在时的获取方法
}
type Option[K, V any] func(conf *Config[K, V]) *Config[K, V]

func NewDefaultConf[K, V any]() *Config[K, V] {
	return &Config[K, V]{
		capacity:         100,
		expireAfterWrite: time.Minute,
		clearInterval:    time.Minute * 5,
		minClearInterval: time.Second * 3,
		keyToString:      nil,
	}
}

func WithConfig[K, V any](conf *Config[K, V]) Option[K, V] {
	return func(oldConf *Config[K, V]) *Config[K, V] {
		return conf
	}
}

func WithKeyEncoder[K, V any](encoder func(K) string) Option[K, V] {
	return func(conf *Config[K, V]) *Config[K, V] {
		conf.keyToString = encoder
		return conf
	}
}

func WithCapacity[K, V any](capacity int) Option[K, V] {
	if capacity < 0 {
		panic("capacity less than 0")
	}
	return func(conf *Config[K, V]) *Config[K, V] {
		conf.capacity = capacity
		return conf
	}
}

func WithExpireAfterWrite[K, V any](expireAfterWrite time.Duration) Option[K, V] {
	if expireAfterWrite < 0 {
		panic("expireAfterWrite less than 0")
	}
	return func(conf *Config[K, V]) *Config[K, V] {
		conf.expireAfterWrite = expireAfterWrite
		return conf
	}
}

func WithClearInterval[K, V any](clearInterval time.Duration) Option[K, V] {
	if clearInterval < 0 {
		panic("clearInterval less than 0")
	}
	return func(conf *Config[K, V]) *Config[K, V] {
		conf.clearInterval = clearInterval
		return conf
	}
}

func WithMinClearInterval[K, V any](minClearInterval time.Duration) Option[K, V] {
	if minClearInterval < 0 {
		panic("minClearInterval less than 0")
	}
	return func(conf *Config[K, V]) *Config[K, V] {
		conf.minClearInterval = minClearInterval
		return conf
	}
}
func WithGetterFunc[K, V any](getterFunc func(key K) (*V, error)) Option[K, V] {
	return func(conf *Config[K, V]) *Config[K, V] {
		conf.getterFunc = getterFunc
		return conf
	}
}
