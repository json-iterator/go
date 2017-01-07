package jsoniter

import (
	"encoding/json"
	"fmt"
	"testing"
	"github.com/json-iterator/go/require"
	"bytes"
	"strconv"
)

func Test_float64_0(t *testing.T) {
	iter := ParseString(`0`)
	val := iter.ReadFloat64()
	if val != 0 {
		t.Fatal(val)
	}
}

func Test_float64_1_dot_1(t *testing.T) {
	iter := ParseString(`1.1`)
	val := iter.ReadFloat64()
	if val != 1.1 {
		t.Fatal(val)
	}
}

func Test_float32_1_dot_1_comma(t *testing.T) {
	iter := ParseString(`1.1,`)
	val := iter.ReadFloat32()
	if val != 1.1 {
		fmt.Println(iter.Error)
		t.Fatal(val)
	}
}

func Test_write_float32(t *testing.T) {
	vals := []float32{0, 1, -1, 99, 0xff, 0xfff, 0xffff, 0xfffff, 0xffffff, 0x4ffffff, 0xfffffff,
	-0x4ffffff, -0xfffffff, 1.2345, 1.23456, 1.234567, 1.001}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(buf, 4096)
			stream.WriteFloat32(val)
			stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatFloat(float64(val), 'f', -1, 32), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(buf, 10)
	stream.WriteRaw("abcdefg")
	stream.WriteFloat32(1.123456)
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("abcdefg1.123456", buf.String())
}

func Benchmark_jsoniter_float(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		iter := ParseString(`1.1111111111`)
		iter.ReadFloat64()
	}
}

func Benchmark_json_float(b *testing.B) {
	for n := 0; n < b.N; n++ {
		result := float64(0)
		json.Unmarshal([]byte(`1.1`), &result)
	}
}
