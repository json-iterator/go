package jsoniter

import (
	"encoding/json"
	"testing"
	"github.com/json-iterator/go/require"
	"bytes"
)

func Test_empty_array(t *testing.T) {
	iter := ParseString(`[]`)
	cont := iter.ReadArray()
	if cont != false {
		t.FailNow()
	}
}

func Test_one_element(t *testing.T) {
	iter := ParseString(`[1]`)
	cont := iter.ReadArray()
	if cont != true {
		t.FailNow()
	}
	if iter.ReadInt64() != 1 {
		t.FailNow()
	}
	cont = iter.ReadArray()
	if cont != false {
		t.FailNow()
	}
}

func Test_two_elements(t *testing.T) {
	iter := ParseString(`[1,2]`)
	cont := iter.ReadArray()
	if cont != true {
		t.FailNow()
	}
	if iter.ReadInt64() != 1 {
		t.FailNow()
	}
	cont = iter.ReadArray()
	if cont != true {
		t.FailNow()
	}
	if iter.ReadInt64() != 2 {
		t.FailNow()
	}
	cont = iter.ReadArray()
	if cont != false {
		t.FailNow()
	}
}

func Test_invalid_array(t *testing.T) {
	iter := ParseString(`[`)
	iter.ReadArray()
	if iter.Error == nil {
		t.FailNow()
	}
}

func Test_whitespace_in_head(t *testing.T) {
	iter := ParseString(` [1]`)
	cont := iter.ReadArray()
	if cont != true {
		t.FailNow()
	}
	if iter.ReadUint64() != 1 {
		t.FailNow()
	}
}

func Test_whitespace_after_array_start(t *testing.T) {
	iter := ParseString(`[ 1]`)
	cont := iter.ReadArray()
	if cont != true {
		t.FailNow()
	}
	if iter.ReadUint64() != 1 {
		t.FailNow()
	}
}

func Test_whitespace_before_array_end(t *testing.T) {
	iter := ParseString(`[1 ]`)
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
	iter := ParseString(`[1 ,2]`)
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
	stream := NewStream(buf, 4096)
	stream.IndentionStep = 2
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
	val := []int{1,2,3}
	str, err := MarshalToString(val)
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

func Benchmark_jsoniter_array(b *testing.B) {
	b.ReportAllocs()
	input := []byte(`[1,2,3,4,5,6,7,8,9]`)
	iter := ParseBytes(input)
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
