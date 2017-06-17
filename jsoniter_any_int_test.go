package jsoniter

import (
	"testing"
	"github.com/json-iterator/go/require"
	"io"
)

func Test_read_int64_as_any(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString("1234")
	should.Nil(err)
	should.Equal(1234, any.ToInt())
	should.Equal(io.EOF, any.LastError())
	should.Equal("1234", any.ToString())
	should.True(any.ToBool())
}

func Test_int_lazy_any_get(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString("1234")
	should.Nil(err)
	should.Equal(Invalid, any.Get(1, "2").ValueType())
}
