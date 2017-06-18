package jsoniter

import (
	"testing"
	"github.com/json-iterator/go/require"
)

func Test_read_string_as_any(t *testing.T) {
	should := require.New(t)
	any := Get([]byte(`"hello"`))
	should.Equal("hello", any.ToString())
	should.True(any.ToBool())
	any = Get([]byte(`" "`))
	should.False(any.ToBool())
	any = Get([]byte(`"false"`))
	should.False(any.ToBool())
	any = Get([]byte(`"123"`))
	should.Equal(123, any.ToInt())
}

func Test_wrap_string(t *testing.T) {
	should := require.New(t)
	any := WrapString("123")
	should.Equal(123, any.ToInt())
}