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
	copy(iter.buf, []byte{0,0,0,0,0,0})
	if !reflect.DeepEqual(map[string]string{"1": "2", "hello": "world"}, m) {
		fmt.Println(iter.Error)
		t.Fatal(m)
	}
}

func Test_read_map_of_interface(t *testing.T) {
	iter := ParseString(`{"hello": "world"}`)
	m := map[string]interface{}{"1": "2"}
	iter.Read(&m)
	if !reflect.DeepEqual(map[string]interface{}{"1": "2", "hello": "world"}, m) {
		fmt.Println(iter.Error)
		t.Fatal(m)
	}
}

func Test_read_map_of_any(t *testing.T) {
	iter := ParseString(`{"hello": "world"}`)
	m := map[string]Any{"1": *MakeAny("2")}
	iter.Read(&m)
	if !reflect.DeepEqual(map[string]Any{"1": *MakeAny("2"), "hello": *MakeAny("world")}, m) {
		fmt.Println(iter.Error)
		t.Fatal(m)
	}
}
