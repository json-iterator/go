package jsoniter

import (
	"bytes"
	"encoding/json"
	"github.com/json-iterator/go/require"
	"io"
	"testing"
)

func Test_empty_array(t *testing.T) {
	should := require.New(t)
	iter := ParseString(DEFAULT_CONFIG, `[]`)
	cont := iter.ReadArray()
	should.False(cont)
	iter = ParseString(DEFAULT_CONFIG, `[]`)
	iter.ReadArrayCB(func(iter *Iterator) bool {
		should.FailNow("should not call")
		return true
	})
}

func Test_one_element(t *testing.T) {
	should := require.New(t)
	iter := ParseString(DEFAULT_CONFIG, `[1]`)
	should.True(iter.ReadArray())
	should.Equal(1, iter.ReadInt())
	should.False(iter.ReadArray())
	iter = ParseString(DEFAULT_CONFIG, `[1]`)
	iter.ReadArrayCB(func(iter *Iterator) bool {
		should.Equal(1, iter.ReadInt())
		return true
	})
}

func Test_two_elements(t *testing.T) {
	should := require.New(t)
	iter := ParseString(DEFAULT_CONFIG, `[1,2]`)
	should.True(iter.ReadArray())
	should.Equal(int64(1), iter.ReadInt64())
	should.True(iter.ReadArray())
	should.Equal(int64(2), iter.ReadInt64())
	should.False(iter.ReadArray())
	iter = ParseString(DEFAULT_CONFIG, `[1,2]`)
	should.Equal([]interface{}{float64(1), float64(2)}, iter.Read())
}

func Test_read_empty_array_as_any(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString("[]")
	should.Nil(err)
	should.Equal(0, any.Size())
}

func Test_read_one_element_array_as_any(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString("[1]")
	should.Nil(err)
	should.Equal(1, any.Size())
}

func Test_read_two_element_array_as_any(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString("[1,2]")
	should.Nil(err)
	should.Equal(1, any.Get(0).ToInt())
	should.Equal(2, any.Size())
	should.True(any.ToBool())
	should.Equal(1, any.ToInt())
}

func Test_read_array_with_any_iterator(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString("[1,2]")
	should.Nil(err)
	var element Any
	var elements []int
	for next, hasNext := any.IterateArray(); hasNext; {
		element, hasNext = next()
		elements = append(elements, element.ToInt())
	}
	should.Equal([]int{1, 2}, elements)
}

func Test_wrap_array(t *testing.T) {
	should := require.New(t)
	any := Wrap([]int{1, 2, 3})
	should.Equal("[1,2,3]", any.ToString())
	var element Any
	var elements []int
	for next, hasNext := any.IterateArray(); hasNext; {
		element, hasNext = next()
		elements = append(elements, element.ToInt())
	}
	should.Equal([]int{1, 2, 3}, elements)
	any = Wrap([]int{1, 2, 3})
	should.Equal(3, any.Size())
	any = Wrap([]int{1, 2, 3})
	should.Equal(2, any.Get(1).ToInt())
}

func Test_array_lazy_any_get(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString("[1,[2,3],4]")
	should.Nil(err)
	should.Equal(3, any.Get(1, 1).ToInt())
	should.Equal("[1,[2,3],4]", any.ToString())
}

func Test_array_lazy_any_get_all(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString("[[1],[2],[3,4]]")
	should.Nil(err)
	should.Equal("[1,2,3]", any.Get('*', 0).ToString())
}

func Test_array_wrapper_any_get_all(t *testing.T) {
	should := require.New(t)
	any := wrapArray([][]int{
		{1, 2},
		{3, 4},
		{5, 6},
	})
	should.Equal("[1,3,5]", any.Get('*', 0).ToString())
}

func Test_array_lazy_any_get_invalid(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString("[]")
	should.Nil(err)
	should.Equal(Invalid, any.Get(1, 1).ValueType())
	should.NotNil(any.Get(1, 1).LastError())
	should.Equal(Invalid, any.Get("1").ValueType())
	should.NotNil(any.Get("1").LastError())
}

func Test_array_lazy_any_set(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString("[1,[2,3],4]")
	should.Nil(err)
	any.GetArray()[0] = WrapInt64(2)
	str, err := MarshalToString(any)
	should.Nil(err)
	should.Equal("[2,[2,3],4]", str)
}

func Test_invalid_array(t *testing.T) {
	_, err := UnmarshalAnyFromString("[")
	if err == nil || err == io.EOF {
		t.FailNow()
	}
}

func Test_whitespace_in_head(t *testing.T) {
	iter := ParseString(DEFAULT_CONFIG, ` [1]`)
	cont := iter.ReadArray()
	if cont != true {
		t.FailNow()
	}
	if iter.ReadUint64() != 1 {
		t.FailNow()
	}
}

func Test_whitespace_after_array_start(t *testing.T) {
	iter := ParseString(DEFAULT_CONFIG, `[ 1]`)
	cont := iter.ReadArray()
	if cont != true {
		t.FailNow()
	}
	if iter.ReadUint64() != 1 {
		t.FailNow()
	}
}

func Test_whitespace_before_array_end(t *testing.T) {
	iter := ParseString(DEFAULT_CONFIG, `[1 ]`)
	cont := iter.ReadArray()
	if cont != true {
		t.FailNow()
	}
	if iter.ReadUint64() != 1 {
		t.FailNow()
	}
	cont = iter.ReadArray()
	if cont != false {
		t.FailNow()
	}
}

func Test_whitespace_before_comma(t *testing.T) {
	iter := ParseString(DEFAULT_CONFIG, `[1 ,2]`)
	cont := iter.ReadArray()
	if cont != true {
		t.FailNow()
	}
	if iter.ReadUint64() != 1 {
		t.FailNow()
	}
	cont = iter.ReadArray()
	if cont != true {
		t.FailNow()
	}
	if iter.ReadUint64() != 2 {
		t.FailNow()
	}
	cont = iter.ReadArray()
	if cont != false {
		t.FailNow()
	}
}

func Test_write_array(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	newCfg := &Config{IndentionStep: 2}
	initConfig(newCfg)
	stream := NewStream(newCfg, buf, 4096)
	stream.WriteArrayStart()
	stream.WriteInt(1)
	stream.WriteMore()
	stream.WriteInt(2)
	stream.WriteArrayEnd()
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("[\n  1,\n  2\n]", buf.String())
}

func Test_write_val_array(t *testing.T) {
	should := require.New(t)
	val := []int{1, 2, 3}
	str, err := MarshalToString(&val)
	should.Nil(err)
	should.Equal("[1,2,3]", str)
}

func Test_write_val_empty_array(t *testing.T) {
	should := require.New(t)
	val := []int{}
	str, err := MarshalToString(val)
	should.Nil(err)
	should.Equal("[]", str)
}

func Test_write_array_of_interface_in_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field  []interface{}
		Field2 string
	}
	val := TestObject{[]interface{}{1, 2}, ""}
	str, err := MarshalToString(val)
	should.Nil(err)
	should.Contains(str, `"Field":[1,2]`)
	should.Contains(str, `"Field2":""`)
}

func Test_json_RawMessage(t *testing.T) {
	should := require.New(t)
	var data json.RawMessage
	should.Nil(Unmarshal([]byte(`[1,2,3]`), &data))
	should.Equal(`[1,2,3]`, string(data))
	str, err := MarshalToString(data)
	should.Nil(err)
	should.Equal(`[1,2,3]`, str)
}

func Test_encode_byte_array(t *testing.T) {
	should := require.New(t)
	bytes, err := json.Marshal([]byte{1, 2, 3})
	should.Nil(err)
	should.Equal(`"AQID"`, string(bytes))
	bytes, err = Marshal([]byte{1, 2, 3})
	should.Nil(err)
	should.Equal(`"AQID"`, string(bytes))
}

func Test_decode_byte_array(t *testing.T) {
	should := require.New(t)
	data := []byte{}
	err := json.Unmarshal([]byte(`"AQID"`), &data)
	should.Nil(err)
	should.Equal([]byte{1, 2, 3}, data)
	err = Unmarshal([]byte(`"AQID"`), &data)
	should.Nil(err)
	should.Equal([]byte{1, 2, 3}, data)
}

func Benchmark_jsoniter_array(b *testing.B) {
	b.ReportAllocs()
	input := []byte(`[1,2,3,4,5,6,7,8,9]`)
	iter := ParseBytes(DEFAULT_CONFIG, input)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		iter.ResetBytes(input)
		for iter.ReadArray() {
			iter.ReadUint64()
		}
	}
}

func Benchmark_json_array(b *testing.B) {
	for n := 0; n < b.N; n++ {
		result := []interface{}{}
		json.Unmarshal([]byte(`[1,2,3]`), &result)
	}
}
