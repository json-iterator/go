package jsoniter

import (
	"testing"
	"reflect"
	"fmt"
)

func Test_read_map(t *testing.T) {
	iter := ParseString(`{"hello": "world"}`)
	m := map[string]string{"1": "2"}
	iter.Read(&m)
	if !reflect.DeepEqual(map[string]string{"1": "2", "hello": "world"}, m) {
		fmt.Println(iter.Error)
		t.Fatal(m)
	}
}
