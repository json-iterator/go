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
	any = Wrap(map[string]string{"Field1": "hello"})
	vals := map[string]string{}
	var k string
	var v Any
	for next, hasNext := any.IterateObject(); hasNext; {
		k, v, hasNext = next()
		if v.ValueType() == String {
			vals[k] = v.ToString()
		}
	}
	should.Equal(map[string]string{"Field1": "hello"}, vals)
}
