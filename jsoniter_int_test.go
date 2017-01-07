package jsoniter

import (
	"bytes"
	"encoding/json"
	"testing"
	"github.com/json-iterator/go/require"
	"fmt"
	"strconv"
)

func Test_decode_decode_uint64_0(t *testing.T) {
	iter := Parse(bytes.NewBufferString("0"), 4096)
	val := iter.ReadUint64()
	if iter.Error != nil {
		t.Fatal(iter.Error)
	}
	if val != 0 {
		t.Fatal(val)
	}
}

func Test_decode_uint64_1(t *testing.T) {
	iter := Parse(bytes.NewBufferString("1"), 4096)
	val := iter.ReadUint64()
	if val != 1 {
		t.Fatal(val)
	}
}

func Test_decode_uint64_100(t *testing.T) {
	iter := Parse(bytes.NewBufferString("100"), 4096)
	val := iter.ReadUint64()
	if val != 100 {
		t.Fatal(val)
	}
}

func Test_decode_uint64_100_comma(t *testing.T) {
	iter := Parse(bytes.NewBufferString("100,"), 4096)
	val := iter.ReadUint64()
	if iter.Error != nil {
		t.Fatal(iter.Error)
	}
	if val != 100 {
		t.Fatal(val)
	}
}

func Test_decode_uint64_invalid(t *testing.T) {
	iter := Parse(bytes.NewBufferString(","), 4096)
	iter.ReadUint64()
	if iter.Error == nil {
		t.FailNow()
	}
}

func Test_decode_int64_100(t *testing.T) {
	iter := Parse(bytes.NewBufferString("100"), 4096)
	val := iter.ReadInt64()
	if val != 100 {
		t.Fatal(val)
	}
}

func Test_decode_int64_minus_100(t *testing.T) {
	iter := Parse(bytes.NewBufferString("-100"), 4096)
	val := iter.ReadInt64()
	if val != -100 {
		t.Fatal(val)
	}
}

func Test_write_uint8(t *testing.T) {
	vals := []uint8{0, 1, 11, 111, 255}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(buf, 4096)
			stream.WriteUint8(val)
			stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.Itoa(int(val)), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(buf, 3)
	stream.WriteString("a")
	stream.WriteUint8(100) // should clear buffer
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a100", buf.String())
}

func Test_write_int8(t *testing.T) {
	vals := []int8{0, 1, -1, 99, 0x7f, -0x7f}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(buf, 4096)
			stream.WriteInt8(val)
			stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.Itoa(int(val)), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(buf, 4)
	stream.WriteString("a")
	stream.WriteInt8(-100) // should clear buffer
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a-100", buf.String())
}

func Test_write_uint16(t *testing.T) {
	vals := []uint16{0, 1, 11, 111, 255, 0xfff, 0xffff}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(buf, 4096)
			stream.WriteUint16(val)
			stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.Itoa(int(val)), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(buf, 5)
	stream.WriteString("a")
	stream.WriteUint16(10000) // should clear buffer
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a10000", buf.String())
}

func Benchmark_jsoniter_int(b *testing.B) {
	for n := 0; n < b.N; n++ {
		iter := ParseString(`-100`)
		iter.ReadInt64()
	}
}

func Benchmark_json_int(b *testing.B) {
	for n := 0; n < b.N; n++ {
		result := int64(0)
		json.Unmarshal([]byte(`-100`), &result)
	}
}
