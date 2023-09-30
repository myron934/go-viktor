package optional

import (
	"fmt"
	"testing"

	viktor "github.com/myron934/go-viktor"
)

func TestNullable(t *testing.T) {

	arr := []string{"1"}
	var strArr *[]string
	fmt.Println(len(*OfNullable(strArr).OrElse(&arr)))
	var str *string

	fmt.Println(*OfNullable(str).OrElse(viktor.Ptr("13")))

}

func TestOrElseGet(t *testing.T) {
	var strArr *string
	OrElseGet[string](strArr == nil, "val1", func() string { return *strArr })
}

func TestOrElse(t *testing.T) {
	var strArr *string
	OrElse[string](strArr == nil, "val1", "val2")
}
