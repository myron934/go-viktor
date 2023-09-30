package common

import "unsafe"

func Ptr[T any](t T) *T {
	return &t
}

type interfaceStructure struct {
	pt uintptr // 到值类型的指针
	pv uintptr // 到值内容的指针
}

// IsNilPtr 判断指针是否为nil
func IsNilPtr(i interface{}) bool {
	if i == nil {
		return true
	}
	is := (*interfaceStructure)(unsafe.Pointer(&i))
	return is.pv == 0 || is.pt == 0
}
