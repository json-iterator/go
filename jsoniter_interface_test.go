package jsoniter

import (
	"testing"
	"github.com/json-iterator/go/require"
	"unsafe"
)

func Test_write_array_of_interface(t *testing.T) {
	should := require.New(t)
	array := []interface{}{"hello"}
	str, err := MarshalToString(array)
	should.Nil(err)
	should.Equal(`["hello"]`, str)
}

func Test_write_map_of_interface(t *testing.T) {
	should := require.New(t)
	val := map[string]interface{}{"hello":"world"}
	str, err := MarshalToString(val)
	should.Nil(err)
	should.Equal(`{"hello":"world"}`, str)
}

func Test_write_map_of_interface_in_struct(t *testing.T) {
	type TestObject struct {
		Field map[string]interface{}
	}
	should := require.New(t)
	val := TestObject{map[string]interface{}{"hello":"world"}}
	str, err := MarshalToString(val)
	should.Nil(err)
	should.Equal(`{"Field":{"hello":"world"}}`, str)
}

func Test_write_map_of_interface_in_struct_with_two_fields(t *testing.T) {
	type TestObject struct {
		Field map[string]interface{}
		Field2 string
	}
	should := require.New(t)
	val := TestObject{map[string]interface{}{"hello":"world"}, ""}
	str, err := MarshalToString(val)
	should.Nil(err)
	should.Equal(`{"Field":{"hello":"world"},"Field2":""}`, str)
}

type MyInterface interface {
	Hello() string
}

type MyString string

func (ms MyString) Hello() string {
	return string(ms)
}

func Test_write_map_of_custom_interface(t *testing.T) {
	should := require.New(t)
	myStr := MyString("world")
	should.Equal("world", myStr.Hello())
	val := map[string]MyInterface{"hello":myStr}
	str, err := MarshalToString(val)
	should.Nil(err)
	should.Equal(`{"hello":"world"}`, str)
}

func Test_write_interface(t *testing.T) {
	should := require.New(t)
	var val interface{}
	val = "hello"
	str, err := MarshalToString(val)
	should.Nil(err)
	should.Equal(`"hello"`, str)
}

func Test_read_interface(t *testing.T) {
	should := require.New(t)
	var val interface{}
	err := UnmarshalFromString(`"hello"`, &val)
	should.Nil(err)
	should.Equal("hello", val)
}

func Test_read_custom_interface(t *testing.T) {
	should := require.New(t)
	var val MyInterface
	RegisterTypeDecoder("jsoniter.MyInterface", func(ptr unsafe.Pointer, iter *Iterator) {
		*((*MyInterface)(ptr)) = MyString(iter.ReadString())
	})
	err := UnmarshalFromString(`"hello"`, &val)
	should.Nil(err)
	should.Equal("hello", val.Hello())
}