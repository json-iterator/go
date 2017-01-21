package jsoniter

import (
	"encoding/json"
	"fmt"
	"testing"
	"github.com/json-iterator/go/require"
	"bytes"
	"strconv"
)

func Test_read_float(t *testing.T) {
	inputs := []string{`1.1`, `1000`, `9223372036854775807`, `12.3`, `-12.3`, `720368.54775807`, `720368.547758075`}
	for _, input := range inputs {
		// non-streaming
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := ParseString(input + ",")
			expected, err := strconv.ParseFloat(input, 32)
			should.Nil(err)
			should.Equal(float32(expected), iter.ReadFloat32())
		})
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := ParseString(input + ",")
			expected, err := strconv.ParseFloat(input, 64)
			should.Nil(err)
			should.Equal(expected, iter.ReadFloat64())
		})
		// streaming
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := Parse(bytes.NewBufferString(input + ","), 2)
			expected, err := strconv.ParseFloat(input, 32)
			should.Nil(err)
			should.Equal(float32(expected), iter.ReadFloat32())
		})
		t.Run(fmt.Sprintf("%v", input), func(t *testing.T) {
			should := require.New(t)
			iter := Parse(bytes.NewBufferString(input + ","), 2)
			expected, err := strconv.ParseFloat(input, 64)
			should.Nil(err)
			should.Equal(expected, iter.ReadFloat64())
		})
	}
}

func Test_read_float_via_read(t *testing.T) {
	should := require.New(t)
	iter := ParseString(`12.3`)
	should.Equal(float64(12.3), iter.Read())
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
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(buf, 4096)
			stream.WriteVal(val)
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

func Test_write_float64(t *testing.T) {
	vals := []float64{0, 1, -1, 99, 0xff, 0xfff, 0xffff, 0xfffff, 0xffffff, 0x4ffffff, 0xfffffff,
	-0x4ffffff, -0xfffffff, 1.2345, 1.23456, 1.234567, 1.001}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(buf, 4096)
			stream.WriteFloat64(val)
			stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatFloat(val, 'f', -1, 64), buf.String())
		})
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(buf, 4096)
			stream.WriteVal(val)
			stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatFloat(val, 'f', -1, 64), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(buf, 10)
	stream.WriteRaw("abcdefg")
	stream.WriteFloat64(1.123456)
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("abcdefg1.123456", buf.String())
}

func Benchmark_jsoniter_float(b *testing.B) {
	b.ReportAllocs()
	input := []byte(`1.1123,`)
	iter := NewIterator()
	for n := 0; n < b.N; n++ {
		iter.ResetBytes(input)
		iter.ReadFloat64()
	}
}

func Benchmark_json_float(b *testing.B) {
	for n := 0; n < b.N; n++ {
		result := float64(0)
		json.Unmarshal([]byte(`1.1`), &result)
	}
}
