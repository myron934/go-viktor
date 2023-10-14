package optional

type Optional[T any] struct {
	t            *T
	isPresentFun func(t T) bool
}

// OrElse 如果存在该值，返回值， 否则返回 other。
func (o Optional[T]) OrElse(other *T) *T {
	if o.IsPresent() {
		return o.t
	}
	return other
}

// OrElseGet 如果存在该值，返回值， 否则触发 other，并返回 other 调用的结果
func (o Optional[T]) OrElseGet(other func() *T) *T {
	if o.IsPresent() {
		return o.t
	}
	return other()
}

// OrElseV 如果存在该值，返回值， 否则返回 other。
func (o Optional[T]) OrElseV(other T) T {
	if o.IsPresent() {
		return *o.t
	}
	return other
}

// OrElseGetV 如果存在该值，返回值， 否则触发 other，并返回 other 调用的结果
func (o Optional[T]) OrElseGetV(other func() T) T {
	if o.IsPresent() {
		return *o.t
	}
	return other()
}

// IsPresent 如果值存在返回 true，否则返回 false。
func (o Optional[T]) IsPresent() bool {
	if o.isPresentFun != nil {
		return o.isPresentFun(*o.t)
	}
	return o.t != nil
}

// IfPresent 如果值存在，使用该值调用指定的fun。
func (o Optional[T]) IfPresent(fun func(t *T)) {
	if !o.IsPresent() {
		return
	}
	fun(o.t)
	return
}

// Map 如果返回值不为 null，则创建包含映射返回值的Optional[T]作为map方法返回值，否则返回空Optional[T]。
func (o Optional[T]) Map(mapper func(t *T) *T) Optional[T] {
	if !o.IsPresent() {
		return Empty[T]()
	}
	return Optional[T]{t: mapper(o.t)}
}

// FlatMap 如果返回值不为 null，则创建包含映射返回值的Optional作为map方法返回值，否则返回空Optional。
// go的泛型真的鸡肋, FlatMap做不了成员方法
func FlatMap[T1, T2 any](o Optional[T1], mapper func(t *T1) *T2) Optional[T2] {
	if !o.IsPresent() {
		return Empty[T2]()
	}
	return Optional[T2]{t: mapper(o.t)}
}

// Of 返回一个Optional，其值可能是空，也可能包含给定的非空值。
func Of[T any](t *T) Optional[T] {
	return Optional[T]{t: t}
}

// OfV 返回一个Optional, 传值, 而不是指针, 传值需要提供一个判断0值的方法来判断值是否存在
func OfV[T any](t T, isPresentFun func(t T) bool) Optional[T] {
	return Optional[T]{
		t:            &t,
		isPresentFun: isPresentFun,
	}
}

// Empty 返回一个空Optional
func Empty[T any]() Optional[T] {
	return Optional[T]{t: nil}
}

// OrElse condition==true 返回 conditionVal; 否则返回other
func OrElse[T any](condition bool, conditionVal T, other T) T {
	if condition {
		return conditionVal
	}
	return other
}

// OrElseGet condition==true 返回 conditionVal; 否则返回 getFunc获取的值
func OrElseGet[T any](condition bool, conditionVal T, other func() T) T {
	if condition {
		return conditionVal
	}
	return other()
}
