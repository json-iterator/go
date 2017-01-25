package jsoniter

import (
	"testing"
	"github.com/json-iterator/go/require"
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

type MyInterface interface {
}

func Test_write_map_of_custom_interface(t *testing.T) {
	should := require.New(t)
	val := map[string]MyInterface{"hello":"world"}
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