package concurrency

import (
	"fmt"
	"testing"

	viktor "github.com/myron934/go-viktor"
)

func TestSubmit(t *testing.T) {
	f1 := Submit(func() (*int64, error) {
		var data int64 = 1
		fmt.Printf("f1==== data=%v\n", &data)
		return &data, nil
	})
	if data, err := f1.Get(); err != nil {
		fmt.Printf("f1 err=%v\n", err)
	} else if data != nil {
		fmt.Printf("f1 data=%v\n", *data)
	}

	f2 := Submit[[]string](func() (*[]string, error) {
		data := []string{"1"}
		return &data, nil
	})
	if data, err := f2.Get(); err != nil {
		fmt.Printf("f2 err=%v\n", err)
	} else if data != nil {
		fmt.Printf("f2 data=%v\n", *data)
	}

}

func TestSubmitSlice(t *testing.T) {
	f1 := SubmitSlice(func() ([]int64, error) {
		return []int64{1, 2, 3}, nil
	})
	if list, err := f1.Get(); err != nil {
		fmt.Printf("f1 err=%v\n", err)
	} else {
		fmt.Printf("f1 list=%v\n", list)
	}

	// 元素是指针
	f2 := SubmitSlice(func() ([]*int64, error) {
		var data1, data2, data3 int64 = 1, 2, 3
		return []*int64{&data1, &data2, &data3}, nil
	})
	if list, err := f2.Get(); err != nil {
		fmt.Printf("f2 err=%v\n", err)
	} else {
		fmt.Printf("f2 list=%v\n", list)
	}
}

func TestSubmitMap(t *testing.T) {
	f1 := SubmitMap(func() (map[string]int64, error) {
		return map[string]int64{
			"k1": 1,
			"k2": 2,
		}, nil
	})
	if mp, err := f1.Get(); err != nil {
		fmt.Printf("f1 err=%v\n", err)
	} else {
		fmt.Printf("f1 map=%v\n", mp)
	}

	// 元素是指针
	f2 := SubmitMap(func() (map[string]*int64, error) {
		return map[string]*int64{
			"k1": viktor.Ptr[int64](1),
			"k2": viktor.Ptr[int64](2),
		}, nil
	})
	if mp, err := f2.Get(); err != nil {
		fmt.Printf("f2 err=%v\n", err)
	} else {
		fmt.Printf("f2 map=%v\n", mp)
	}
}
