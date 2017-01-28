package jsoniter

import (
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
	should.Equal(map[string]string{"Field1":"hello"}, vals)
}

func Test_write_val_map(t *testing.T) {
	should := require.New(t)
	val := map[string]string{"1": "2"}
	str, err := MarshalToString(val)
	should.Nil(err)
	should.Equal(`{"1":"2"}`, str)
}
