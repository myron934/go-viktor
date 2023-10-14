package optional

import (
	"fmt"
	"testing"

	viktor "github.com/myron934/go-viktor"
)

func TestOf(t *testing.T) {
	fmt.Println(Of(viktor.Ptr(1)).OrElse(viktor.Ptr(2)))
}

func TestOfV(t *testing.T) {
	fmt.Println(OfV(1, func(t int) bool { return t != 0 }).OrElseV(2))
	fmt.Println(OfV("a", func(t string) bool { return t != "" }).OrElseV("b"))
}

func TestOrElseGet(t *testing.T) {
	var pstr *string
	fmt.Println(OrElseGet(pstr == nil, "val1", func() string { return *pstr }))
	var str string
	fmt.Println(OrElseGet(str == "", "val1", func() string { return str + "val2" }))

	type Obj struct {
		Id   *int
		Name *string
	}

	obj := &Obj{}
	name := OrElseGet(obj == nil, "obj is nil", func() string { return Of(obj.Name).OrElseV("obj.name is nil, too") })
	fmt.Println(name)

	obj.Name = viktor.Ptr("hello world")
	name = OrElseGet(obj == nil, "obj is nil", func() string { return Of(obj.Name).OrElseV("obj.name is nil, too") })
	fmt.Println(name)

	obj = nil
	name = OrElseGet(obj == nil, "obj is nil", func() string { return Of(obj.Name).OrElseV("obj.name is nil, too") })
	fmt.Println(name)

	obj2 := &Obj{
		Id:   OrElseGet(obj == nil, nil, func() *int { return obj.Id }),
		Name: OrElseGet(obj == nil, viktor.Ptr("obj is nil"), func() *string { return obj.Name }),
	}
	fmt.Println(obj2)
}

func TestOrElse(t *testing.T) {
	var pstr *string
	fmt.Println(OrElse(pstr == nil, "val1", "val2"))
}
