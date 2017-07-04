package jsoniter

import (
	"fmt"
	"io"
	"testing"

	"github.com/json-iterator/go/require"
)

var intConvertMap = map[string]int{
	"321.1":      321,
	"-321.1":     -321,
	`"1.1"`:      1,
	`"-1.1"`:     -1,
	"0.0":        0,
	"0":          0,
	`"0"`:        0,
	`"0.0"`:      0,
	"-1.1":       -1,
	"true":       1,
	"false":      0,
	`"true"`:     0,
	`"false"`:    0,
	`"true123"`:  0,
	`"123true"`:  123,
	`"1.2332e6"`: 1,
	`""`:         0,
	"+":          0,
	"-":          0,
	"[]":         0,
	"[1,2]":      1,
	// object in php cannot convert to int
	"{}": 0,
}

func Test_read_any_to_int(t *testing.T) {
	should := require.New(t)

	// int
	for k, v := range intConvertMap {
		any := Get([]byte(k))
		should.Equal(v, any.ToInt(), fmt.Sprintf("origin val %v", k))
	}

	// int32
	for k, v := range intConvertMap {
		any := Get([]byte(k))
		should.Equal(int32(v), any.ToInt32(), fmt.Sprintf("original val is %v", k))
	}

	// int64
	for k, v := range intConvertMap {
		any := Get([]byte(k))
		should.Equal(int64(v), any.ToInt64(), fmt.Sprintf("original val is %v", k))
	}

}

var uintConvertMap = map[string]int{
	"321.1":      321,
	`"1.1"`:      1,
	`"-1.1"`:     1,
	"0.0":        0,
	"0":          0,
	`"0"`:        0,
	`"0.0"`:      0,
	"true":       1,
	"false":      0,
	`"true"`:     0,
	`"false"`:    0,
	`"true123"`:  0,
	`"123true"`:  123,
	`"1.2332e6"`: 1,
	`""`:         0,
	"+":          0,
	"-":          0,
	"[]":         0,
	"[1,2]":      1,
	"{}":         0,
	// TODO need to solve
	//"-1.1":       1,
	//"-321.1": 321,
}

func Test_read_any_to_uint(t *testing.T) {
	should := require.New(t)

	for k, v := range uintConvertMap {
		any := Get([]byte(k))
		should.Equal(uint64(v), any.ToUint64(), fmt.Sprintf("origin val %v", k))
	}

	for k, v := range uintConvertMap {
		any := Get([]byte(k))
		should.Equal(uint32(v), any.ToUint32(), fmt.Sprintf("origin val %v", k))
	}

	for k, v := range uintConvertMap {
		any := Get([]byte(k))
		should.Equal(uint32(v), any.ToUint32(), fmt.Sprintf("origin val %v", k))
	}

}

func Test_read_int64_as_any(t *testing.T) {
	should := require.New(t)
	any := Get([]byte("1234"))
	should.Equal(1234, any.ToInt())
	should.Equal(io.EOF, any.LastError())
	should.Equal("1234", any.ToString())
	should.True(any.ToBool())
}

func Test_int_lazy_any_get(t *testing.T) {
	should := require.New(t)
	any := Get([]byte("1234"))
	should.Equal(Invalid, any.Get(1, "2").ValueType())
}
