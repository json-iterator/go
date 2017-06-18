package jsoniter

import (
	"testing"
	"github.com/json-iterator/go/require"
)

func Test_read_object_as_any(t *testing.T) {
	should := require.New(t)
	any := Get([]byte(`{"a":"b","c":"d"}`))
	should.Equal(`{"a":"b","c":"d"}`, any.ToString())
	// partial parse
	should.Equal("b", any.Get("a").ToString())
	should.Equal("d", any.Get("c").ToString())
	should.Equal(2, len(any.Keys()))
	any = Get([]byte(`{"a":"b","c":"d"}`))
	// full parse
	should.Equal(2, len(any.Keys()))
	should.Equal(2, any.Size())
	should.True(any.ToBool())
	should.Equal(1, any.ToInt())
	should.Equal(Object, any.ValueType())
	should.Nil(any.LastError())
	should.Equal("b", any.GetObject()["a"].ToString())
}

func Test_object_lazy_any_get(t *testing.T) {
	should := require.New(t)
	any := Get([]byte(`{"a":{"b":{"c":"d"}}}`))
	should.Equal("d", any.Get("a", "b", "c").ToString())
}

func Test_object_lazy_any_get_all(t *testing.T) {
	should := require.New(t)
	any := Get([]byte(`{"a":[0],"b":[1]}`))
	should.Contains(any.Get('*', 0).ToString(), `"a":0`)
}

func Test_object_lazy_any_get_invalid(t *testing.T) {
	should := require.New(t)
	any := Get([]byte(`{}`))
	should.Equal(Invalid, any.Get("a", "b", "c").ValueType())
	should.Equal(Invalid, any.Get(1).ValueType())
}

func Test_wrap_object(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 string
		field2 string
	}
	any := Wrap(TestObject{"hello", "world"})
	should.Equal("hello", any.Get("Field1").ToString())
	any = Wrap(TestObject{"hello", "world"})
	should.Equal(2, any.Size())
}

func Test_any_within_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 Any
		Field2 Any
	}
	obj := TestObject{}
	err := UnmarshalFromString(`{"Field1": "hello", "Field2": [1,2,3]}`, &obj)
	should.Nil(err)
	should.Equal("hello", obj.Field1.ToString())
	should.Equal("[1,2,3]", obj.Field2.ToString())
}
