package jsoniter

import (
	"testing"
	"github.com/json-iterator/go/require"
)

func Test_read_string_as_any(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString(`"hello"`)
	should.Nil(err)
	should.Equal("hello", any.ToString())
	should.True(any.ToBool())
	any, err = UnmarshalAnyFromString(`" "`)
	should.False(any.ToBool())
	any, err = UnmarshalAnyFromString(`"false"`)
	should.False(any.ToBool())
	any, err = UnmarshalAnyFromString(`"123"`)
	should.Equal(123, any.ToInt())
}

func Test_wrap_string(t *testing.T) {
	should := require.New(t)
	any := WrapString("123")
	should.Equal(123, any.ToInt())
}