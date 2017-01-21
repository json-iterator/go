package jsoniter

import (
	"fmt"
	"reflect"
	"testing"
	"github.com/json-iterator/go/require"
)

func Test_read_map(t *testing.T) {
	should := require.New(t)
	iter := ParseString(`{"hello": "world"}`)
	m := map[string]string{"1": "2"}
	iter.ReadVal(&m)
	copy(iter.buf, []byte{0, 0, 0, 0, 0, 0})
	should.Equal(map[string]string{"1": "2", "hello": "world"}, m)
}

func Test_read_map_of_interface(t *testing.T) {
	should := require.New(t)
	iter := ParseString(`{"hello": "world"}`)
	m := map[string]interface{}{"1": "2"}
	iter.ReadVal(&m)
	should.Equal(map[string]interface{}{"1": "2", "hello": "world"}, m)
	iter = ParseString(`{"hello": "world"}`)
	should.Equal(map[string]interface{}{"hello": "world"}, iter.Read())
}

func Test_read_map_of_any(t *testing.T) {
	iter := ParseString(`{"hello": "world"}`)
	m := map[string]Any{"1": *MakeAny("2")}
	iter.ReadVal(&m)
	if !reflect.DeepEqual(map[string]Any{"1": *MakeAny("2"), "hello": *MakeAny("world")}, m) {
		fmt.Println(iter.Error)
		t.Fatal(m)
	}
}

func Test_write_val_map(t *testing.T) {
	should := require.New(t)
	val := map[string]string{"1": "2"}
	str, err := MarshalToString(val)
	should.Nil(err)
	should.Equal(`{"1":"2"}`, str)
}
