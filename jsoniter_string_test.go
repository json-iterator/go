package jsoniter

import (
	"bytes"
	"encoding/json"
	"testing"
	"github.com/json-iterator/go/require"
	"fmt"
)

func Test_read_normal_string(t *testing.T) {
	cases := map[string]string{
		`"0123456789012345678901234567890123456789"`: `0123456789012345678901234567890123456789`,
		`""`: ``,
		`"hello"`: `hello`,
	}
	for input, output := range cases {
		t.Run(fmt.Sprintf("%v:%v", input, output), func(t *testing.T) {
			should := require.New(t)
			iter := ParseString(input)
			should.Equal(output, iter.ReadString())
		})
		t.Run(fmt.Sprintf("%v:%v", input, output), func(t *testing.T) {
			should := require.New(t)
			iter := Parse(bytes.NewBufferString(input), 2)
			should.Equal(output, iter.ReadString())
		})
		t.Run(fmt.Sprintf("%v:%v", input, output), func(t *testing.T) {
			should := require.New(t)
			iter := ParseString(input)
			should.Equal(output, string(iter.ReadStringAsSlice()))
		})
		t.Run(fmt.Sprintf("%v:%v", input, output), func(t *testing.T) {
			should := require.New(t)
			iter := Parse(bytes.NewBufferString(input), 2)
			should.Equal(output, string(iter.ReadStringAsSlice()))
		})
	}
}

func Test_read_exotic_string(t *testing.T) {
	cases := map[string]string{
		`"hel\"lo"`: `hel"lo`,
		`"hel\nlo"`: "hel\nlo",
		`"\u4e2d\u6587"`: "中文",
		`"\ud83d\udc4a"`: "\xf0\x9f\x91\x8a", // surrogate
	}
	for input, output := range cases {
		t.Run(fmt.Sprintf("%v:%v", input, output), func(t *testing.T) {
			should := require.New(t)
			iter := ParseString(input)
			should.Equal(output, iter.ReadString())
		})
		t.Run(fmt.Sprintf("%v:%v", input, output), func(t *testing.T) {
			should := require.New(t)
			iter := Parse(bytes.NewBufferString(input), 2)
			should.Equal(output, iter.ReadString())
		})
	}
}

func Test_write_string(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(buf, 4096)
	stream.WriteString("hello")
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal(`"hello"`, buf.String())
}

func Test_write_val_string(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(buf, 4096)
	stream.WriteVal("hello")
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal(`"hello"`, buf.String())
}

func Benchmark_jsoniter_unicode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		iter := ParseString(`"\ud83d\udc4a"`)
		iter.ReadString()
	}
}

func Benchmark_jsoniter_ascii(b *testing.B) {
	iter := NewIterator()
	input := []byte(`"hello, world! hello, world!"`)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		iter.ResetBytes(input)
		iter.ReadString()
	}
}

func Benchmark_jsoniter_string_as_bytes(b *testing.B) {
	iter := ParseString(`"hello, world!"`)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		iter.ResetBytes(iter.buf)
		iter.ReadStringAsSlice()
	}
}

func Benchmark_json_unicode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		result := ""
		json.Unmarshal([]byte(`"\ud83d\udc4a"`), &result)
	}
}

func Benchmark_json_ascii(b *testing.B) {
	for n := 0; n < b.N; n++ {
		result := ""
		json.Unmarshal([]byte(`"hello"`), &result)
	}
}
