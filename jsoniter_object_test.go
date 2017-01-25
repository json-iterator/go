package jsoniter

import (
	"encoding/json"
	"testing"
	"github.com/json-iterator/go/require"
	"bytes"
)

func Test_empty_object(t *testing.T) {
	should := require.New(t)
	iter := ParseString(`{}`)
	field := iter.ReadObject()
	should.Equal("", field)
	iter = ParseString(`{}`)
	iter.ReadObjectCB(func(iter *Iterator, field string) bool {
		should.FailNow("should not call")
		return true
	})
}

func Test_one_field(t *testing.T) {
	should := require.New(t)
	iter := ParseString(`{"a": "b"}`)
	field := iter.ReadObject()
	should.Equal("a", field)
	value := iter.ReadString()
	should.Equal("b", value)
	field = iter.ReadObject()
	should.Equal("", field)
	iter = ParseString(`{"a": "b"}`)
	should.True(iter.ReadObjectCB(func(iter *Iterator, field string) bool {
		should.Equal("a", field)
		return true
	}))
}

func Test_two_field(t *testing.T) {
	should := require.New(t)
	iter := ParseString(`{ "a": "b" , "c": "d" }`)
	field := iter.ReadObject()
	should.Equal("a", field)
	value := iter.ReadString()
	should.Equal("b", value)
	field = iter.ReadObject()
	should.Equal("c", field)
	value = iter.ReadString()
	should.Equal("d", value)
	field = iter.ReadObject()
	should.Equal("", field)
	iter = ParseString(`{"field1": "1", "field2": 2}`)
	for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
		switch field {
		case "field1":
			iter.ReadString()
		case "field2":
			iter.ReadInt64()
		default:
			iter.reportError("bind object", "unexpected field")
		}
	}
}

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

	should.Equal(map[string]string{"a":"b", "c":"d"}, vals)
	vals = map[string]string{}
	for next, hasNext := any.IterateObject(); hasNext; k, v, hasNext = next() {
		vals[k] = v.ToString()
	}
	should.Equal(map[string]string{"a":"b", "c":"d"}, vals)
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

func Test_object_lazy_any_set(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString(`{"a":{"b":{"c":"d"}}}`)
	should.Nil(err)
	any.GetObject()["a"] = WrapInt64(1)
	str, err := MarshalToString(any)
	should.Nil(err)
	should.Equal(`{"a":1}`, str)
}

func Test_write_object(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(buf, 4096)
	stream.IndentionStep = 2
	stream.WriteObjectStart()
	stream.WriteObjectField("hello")
	stream.WriteInt(1)
	stream.WriteMore()
	stream.WriteObjectField("world")
	stream.WriteInt(2)
	stream.WriteObjectEnd()
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("{\n  \"hello\":1,\n  \"world\":2\n}", buf.String())
}

type TestObj struct {
	Field1 string
	Field2 uint64
}

func Benchmark_jsoniter_object(b *testing.B) {
	for n := 0; n < b.N; n++ {
		iter := ParseString(`{"field1": "1", "field2": 2}`)
		obj := TestObj{}
		for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
			switch field {
			case "field1":
				obj.Field1 = iter.ReadString()
			case "field2":
				obj.Field2 = iter.ReadUint64()
			default:
				iter.reportError("bind object", "unexpected field")
			}
		}
	}
}

func Benchmark_json_object(b *testing.B) {
	for n := 0; n < b.N; n++ {
		result := TestObj{}
		json.Unmarshal([]byte(`{"field1": "1", "field2": 2}`), &result)
	}
}
