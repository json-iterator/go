package jsoniter

import (
	"testing"
	"fmt"
)

func Test_bind_api_demo(t *testing.T) {
	iter := ParseString(`[0,1,2,3]`)
	val := []int{}
	iter.Read(&val)
	fmt.Println(val[3])
}

func Test_any_api_demo(t *testing.T) {
	iter := ParseString(`[0,1,2,3]`)
	val := iter.ReadAny()
	fmt.Println(val.Get(3))
}

func Test_iterator_api_demo(t *testing.T) {
	iter := ParseString(`[0,1,2,3]`)
	total := 0
	for iter.ReadArray() {
		total += iter.ReadInt()
	}
	fmt.Println(total)
}

type ABC struct {
	a Any
}

func Test_deep_nested_any_api(t *testing.T) {
	iter := ParseString(`{"a": {"b": {"c": "d"}}}`)
	abc := &ABC{}
	iter.Read(&abc)
	fmt.Println(abc.a.Get("b", "c"))
}

type User struct {
	userId int
	name string
	tags []string
}

func Test_iterator_and_bind_api(t *testing.T) {
	iter := ParseString(`[123, {"name": "taowen", "tags": ["crazy", "hacker"]}]`)
	user := User{}
	iter.ReadArray()
	user.userId = iter.ReadInt()
	iter.ReadArray()
	iter.Read(&user)
	iter.ReadArray() // array end
	fmt.Println(user)
}
