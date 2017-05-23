package jsoniter

import (
	"fmt"
	"testing"
	"github.com/json-iterator/go/require"
)

func Test_bind_api_demo(t *testing.T) {
	should := require.New(t)
	val := []int{}
	err := UnmarshalFromString(`[0,1,2,3]  `, &val)
	should.Nil(err)
	should.Equal([]int{0, 1, 2, 3}, val)
}

func Test_iterator_api_demo(t *testing.T) {
	iter := ParseString(`[0,1,2,3]`)
	total := 0
	for iter.ReadArray() {
		total += iter.ReadInt()
	}
	fmt.Println(total)
}
