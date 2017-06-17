package jsoniter

import (
	"testing"
	"github.com/json-iterator/go/require"
)

func Test_read_object_as_any(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString(`{"a":"b","c":"d"}`)
	should.Nil(err)
	should.Equal(`{"a":"b","c":"d"}`, any.ToString())
	// partial parse
	should.Equal("b", any.Get("a").ToString())
	should.Equal("d", any.Get("c").ToString())
	should.Equal(2, len(any.Keys()))
	any, err = UnmarshalAnyFromString(`{"a":"b","c":"d"}`)
	// full parse
	should.Equal(2, len(any.Keys()))
	should.Equal(2, any.Size())
	should.True(any.ToBool())
	should.Equal(1, any.ToInt())
}

func Test_object_any_lazy_iterator(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString(`{"a":"b","c":"d"}`)
	should.Nil(err)
	// iterator parse
	vals := map[string]string{}
	var k string
	var v Any
	next, hasNext := any.IterateObject()
	should.True(hasNext)

	k, v, hasNext = next()
	should.True(hasNext)
	vals[k] = v.ToString()

	// trigger full parse
	should.Equal(2, len(any.Keys()))

	k, v, hasNext = next()
	should.False(hasNext)
	vals[k] = v.ToString()

	should.Equal(map[string]string{"a": "b", "c": "d"}, vals)
	vals = map[string]string{}
	for next, hasNext := any.IterateObject(); hasNext; {
		k, v, hasNext = next()
		if v.ValueType() == String {
			vals[k] = v.ToString()
		}
	}
	should.Equal(map[string]string{"a": "b", "c": "d"}, vals)
}

func Test_object_any_with_two_lazy_iterators(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString(`{"a":"b","c":"d","e":"f"}`)
	should.Nil(err)
	var k string
	var v Any
	next1, hasNext1 := any.IterateObject()
	next2, hasNext2 := any.IterateObject()
	should.True(hasNext1)
	k, v, hasNext1 = next1()
	should.True(hasNext1)
	should.Equal("a", k)
	should.Equal("b", v.ToString())

	should.True(hasNext2)
	k, v, hasNext2 = next2()
	should.True(hasNext2)
	should.Equal("a", k)
	should.Equal("b", v.ToString())

	k, v, hasNext1 = next1()
	should.True(hasNext1)
	should.Equal("c", k)
	should.Equal("d", v.ToString())

	k, v, hasNext2 = next2()
	should.True(hasNext2)
	should.Equal("c", k)
	should.Equal("d", v.ToString())
}

func Test_object_lazy_any_get(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString(`{"a":{"b":{"c":"d"}}}`)
	should.Nil(err)
	should.Equal("d", any.Get("a", "b", "c").ToString())
}

func Test_object_lazy_any_get_all(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString(`{"a":[0],"b":[1]}`)
	should.Nil(err)
	should.Contains(any.Get('*', 0).ToString(), `"a":0`)
}

func Test_object_lazy_any_get_invalid(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString(`{}`)
	should.Nil(err)
	should.Equal(Invalid, any.Get("a", "b", "c").ValueType())
	should.Equal(Invalid, any.Get(1).ValueType())
}

func Test_object_lazy_any_set(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString(`{"a":{"b":{"c":"d"}}}`)
	should.Nil(err)
	any.GetObject()["a"] = WrapInt64(1)
	str, err := MarshalToString(any)
	should.Nil(err)
	should.Equal(`{"a":1}`, str)
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
	any = Wrap(TestObject{"hello", "world"})
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
