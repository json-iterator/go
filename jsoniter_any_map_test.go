package jsoniter

import (
	"testing"
	"github.com/json-iterator/go/require"
)

func Test_wrap_map(t *testing.T) {
	should := require.New(t)
	any := Wrap(map[string]string{"Field1": "hello"})
	should.Equal("hello", any.Get("Field1").ToString())
	any = Wrap(map[string]string{"Field1": "hello"})
	should.Equal(1, any.Size())
}
