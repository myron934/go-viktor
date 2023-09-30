package optional

type GetPtrFunc[T any] func() *T
type GetFunc[T any] func() T

type Optional[T any] struct {
	t *T
}

func (o Optional[T]) OrElse(other *T) *T {
	if o.isPresent() {
		return o.t
	}
	return other
}

func (o Optional[T]) OrElseValue(other T) T {
	if o.isPresent() {
		return *o.t
	}
	return other
}

func (o Optional[T]) isPresent() bool {
	return o.t != nil
}

func (o Optional[T]) OrElseGet(other GetPtrFunc[T]) *T {
	if o.isPresent() {
		return o.t
	}
	return other()
}

func (o Optional[T]) OrElseGetValue(other GetFunc[T]) T {
	if o.isPresent() {
		return *o.t
	}
	return other()
}

// OfNullable t可为空.
func OfNullable[T any](t *T) Optional[T] {
	return Optional[T]{t: t}
}

func Of[T any](t *T) Optional[T] {
	if t == nil {
		panic("[Of] sli is nil")
	}
	return Optional[T]{t: t}
}

// OrElse condition==true conditionVal; 否则返回other
func OrElse[T any](condition bool, conditionVal T, other T) T {
	if condition {
		return conditionVal
	}
	return other
}

// OrElseGet condition==true conditionVal; 否则返回 getFunc获取的值
func OrElseGet[T any](condition bool, conditionVal T, getFunc GetFunc[T]) T {
	if condition {
		return conditionVal
	}
	return getFunc()
}
