package tool

import (
	"fmt"
	"strings"
	"testing"
)

func TestSplitStringArr(t *testing.T) {
	arr := make([]string, 0)
	for i := 0; i < 99; i++ {
		arr = append(arr, fmt.Sprintf("%d", i))
	}
	list := SplitStringArr(arr, 200)
	if len(list) != 1 {
		t.Fatalf("want 1 get %d", len(list))
	}
	if len(list[0]) != 99 {
		t.Fatalf("want 99 get %d", len(list[0]))
	}
	//
	arr = make([]string, 0)
	for i := 0; i < 199; i++ {
		arr = append(arr, fmt.Sprintf("%d", i))
	}
	list = SplitStringArr(arr, 200)
	if len(list) != 1 {
		t.Fatalf("want 1 get %d", len(list))
	}
	if len(list[0]) != 199 {
		t.Fatalf("want 199 get %d", len(list[0]))
	}
	//
	arr = make([]string, 0)
	for i := 0; i < 201; i++ {
		arr = append(arr, fmt.Sprintf("%d", i))
	}
	list = SplitStringArr(arr, 200)
	if len(list) != 2 {
		t.Fatalf("want 2 get %d", len(list))
	}
	if len(list[0]) != 200 {
		t.Fatalf("want 200 get %d", len(list[0]))
	}
	if len(list[1]) != 1 {
		t.Fatalf("want 1 get %d", len(list[1]))
	}
	//
	arr = make([]string, 0)
	for i := 0; i < 399; i++ {
		arr = append(arr, fmt.Sprintf("%d", i))
	}
	list = SplitStringArr(arr, 200)
	if len(list) != 2 {
		t.Fatalf("want 2 get %d", len(list))
	}
	if len(list[0]) != 200 {
		t.Fatalf("want 200 get %d", len(list[0]))
	}
	if len(list[1]) != 199 {
		t.Fatalf("want 199 get %d", len(list[1]))
	}
	//
	arr = make([]string, 0)
	for i := 0; i < 999; i++ {
		arr = append(arr, fmt.Sprintf("%d", i))
	}
	list = SplitStringArr(arr, 200)
	if len(list) != 5 {
		t.Fatalf("want 5 get %d", len(list))
	}
	if len(list[len(list)-1]) != 199 {
		t.Fatalf("want 199 get %d", len(list[len(list)-1]))
	}
	inStr := fmt.Sprintf("('%s')", strings.Join(list[0], "','"))
	fmt.Println(inStr)
}

func TestMd5(t *testing.T) {
	str := Md5("123456")
	fmt.Println(str)
	if str != "E10ADC3949BA59ABBE56E057F20F883E" {
		t.Fail()
	}
}
