package concurrency

import (
	"context"
	"time"
)

type FutureBase struct {
	done chan struct{}
}

// IsDone 判断异步执行是否完成(立即返回结果)
func (f *FutureBase) IsDone() bool {
	select {
	case <-f.done:
		return true
	default:
		return false
	}
}

func newBaseFuture() FutureBase {
	return FutureBase{
		done: make(chan struct{}),
	}
}

// Future 异步执行的结果
type Future[T any] struct {
	FutureBase
	data *T
	err  error
}

// Get 获取执行结果(阻塞)
func (f *Future[T]) Get() (t *T, err error) {
	<-f.done
	return f.data, f.err
}

// GetWithTimeout 获取执行结果. 超时则return error
func (f *Future[T]) GetWithTimeout(timeout time.Duration) (t *T, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer func() {
		cancel()
	}()
	for {
		select {
		case <-f.done:
			return f.data, f.err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (f *Future[T]) setResult(t *T, err error) {
	f.data = t
	f.err = err
	close(f.done)
}

// SliceFuture slice异步执行的结果
type SliceFuture[T any] struct {
	FutureBase
	data []T
	err  error
}

func (f *SliceFuture[T]) setResult(t []T, err error) {
	f.data = t
	f.err = err
	close(f.done)
}

// Get 获取执行结果(阻塞).
func (f *SliceFuture[T]) Get() (t []T, err error) {
	<-f.done
	return f.data, f.err
}

// GetWithTimeout 获取执行结果. 超时则return error
func (f *SliceFuture[T]) GetWithTimeout(timeout time.Duration) (t []T, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer func() {
		cancel()
	}()
	for {
		select {
		case <-f.done:
			return f.data, f.err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// MapFuture map 异步执行的结果
type MapFuture[K string | int | int16 | int32 | int64 | float32 | float64 | bool, V any] struct {
	FutureBase
	mp  map[K]V
	err error
}

func (f *MapFuture[K, V]) setResult(mp map[K]V, err error) {
	f.mp = mp
	f.err = err
	close(f.done)
}

// Get 获取执行结果(阻塞).
func (f *MapFuture[K, V]) Get() (mp map[K]V, err error) {
	<-f.done
	return f.mp, f.err
}

// GetWithTimeout 获取执行结果. 超时则return error
func (f *MapFuture[K, V]) GetWithTimeout(timeout time.Duration) (mp map[K]V, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer func() {
		cancel()
	}()
	for {
		select {
		case <-f.done:
			return f.mp, f.err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func NewFuture[T any]() *Future[T] {
	ft := &Future[T]{
		FutureBase: newBaseFuture(),
	}
	return ft
}

func NewSliceFuture[T any]() *SliceFuture[T] {
	ft := &SliceFuture[T]{
		FutureBase: newBaseFuture(),
	}
	return ft
}

func NewMapFuture[K string | int | int16 | int32 | int64 | float32 | float64 | bool, V any]() *MapFuture[K, V] {
	ft := &MapFuture[K, V]{
		FutureBase: newBaseFuture(),
	}
	return ft
}

// Submit 新启协程执行函数f, 执行完以后将结果放入返回的Futrue中
func Submit[T any](f func() (*T, error)) *Future[T] {
	ft := NewFuture[T]()
	go func() {
		data, err := f()
		ft.setResult(data, err)
	}()
	return ft
}

// SubmitSlice 新启协程执行函数f, 执行完以后将Slice结果放入返回的Futrue中
func SubmitSlice[T any](f func() ([]T, error)) *SliceFuture[T] {
	ft := NewSliceFuture[T]()
	go func() {
		data, err := f()
		ft.setResult(data, err)
	}()
	return ft
}

// SubmitMap 新启协程执行函数f, 执行完以后将Map结果放入返回的Futrue中
func SubmitMap[K string | int | int16 | int32 | int64 | float32 | float64 | bool, V any](f func() (map[K]V, error)) *MapFuture[K, V] {
	ft := NewMapFuture[K, V]()
	go func() {
		mp, err := f()
		ft.setResult(mp, err)
	}()
	return ft
}
