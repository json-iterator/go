package jsoniter

import (
	"fmt"
	"testing"
)

func Test_bind_api_demo(t *testing.T) {
	iter := ParseString(`[0,1,2,3]`)
	val := []int{}
	iter.ReadVal(&val)
	fmt.Println(val[3])
}

func Test_iterator_api_demo(t *testing.T) {
	iter := ParseString(`[0,1,2,3]`)
	total := 0
	for iter.ReadArray() {
		total += iter.ReadInt()
	}
	fmt.Println(total)
}

type User struct {
	userID int
	name   string
	tags   []string
}

func Test_iterator_and_bind_api(t *testing.T) {
	iter := ParseString(`[123, {"name": "taowen", "tags": ["crazy", "hacker"]}]`)
	user := User{}
	iter.ReadArray()
	user.userID = iter.ReadInt()
	iter.ReadArray()
	iter.ReadVal(&user)
	iter.ReadArray() // array end
	fmt.Println(user)
}
