package jsoniter

import (
	"testing"
	"github.com/json-iterator/go/require"
)

func Test_encode_optional_int_pointer(t *testing.T) {
	should := require.New(t)
	var ptr *int
	str, err := MarshalToString(ptr)
	should.Nil(err)
	should.Equal("null", str)
	val := 100
	ptr = &val
	str, err = MarshalToString(ptr)
	should.Nil(err)
	should.Equal("100", str)
}

func Test_decode_struct_with_optional_field(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		field1 *string
		field2 *string
	}
	obj := TestObject{}
	UnmarshalFromString(`{"field1": null, "field2": "world"}`, &obj)
	should.Nil(obj.field1)
	should.Equal("world", *obj.field2)
}

func Test_encode_struct_with_optional_field(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		field1 *string
		field2 *string
	}
	obj := TestObject{}
	world := "world"
	obj.field2 = &world
	str, err := MarshalToString(obj)
	should.Nil(err)
	should.Equal(`{"field1":null,"field2":"world"}`, str)
}