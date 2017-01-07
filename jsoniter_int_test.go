package jsoniter

import (
	"bytes"
	"encoding/json"
	"testing"
	"github.com/json-iterator/go/require"
	"fmt"
	"strconv"
	"io/ioutil"
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
			should.Equal(strconv.FormatUint(uint64(val), 10), buf.String())
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
			should.Equal(strconv.FormatInt(int64(val), 10), buf.String())
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
			should.Equal(strconv.FormatUint(uint64(val), 10), buf.String())
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

func Test_write_int16(t *testing.T) {
	vals := []int16{0, 1, 11, 111, 255, 0xfff, 0x7fff, -0x7fff}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(buf, 4096)
			stream.WriteInt16(val)
			stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatInt(int64(val), 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(buf, 6)
	stream.WriteString("a")
	stream.WriteInt16(-10000) // should clear buffer
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a-10000", buf.String())
}

func Test_write_uint32(t *testing.T) {
	vals := []uint32{0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0xffffffff}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(buf, 4096)
			stream.WriteUint32(val)
			stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatUint(uint64(val), 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(buf, 10)
	stream.WriteString("a")
	stream.WriteUint32(0xffffffff) // should clear buffer
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a4294967295", buf.String())
}

func Test_write_int32(t *testing.T) {
	vals := []int32{0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0x7fffffff, -0x7fffffff}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(buf, 4096)
			stream.WriteInt32(val)
			stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatInt(int64(val), 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(buf, 11)
	stream.WriteString("a")
	stream.WriteInt32(-0x7fffffff) // should clear buffer
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a-2147483647", buf.String())
}

func Test_write_uint64(t *testing.T) {
	vals := []uint64{0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0xffffffff,
		0xfffffffff,0xffffffffff,0xfffffffffff,0xffffffffffff,0xfffffffffffff,0xffffffffffffff,
		0xfffffffffffffff,0xffffffffffffffff}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(buf, 4096)
			stream.WriteUint64(val)
			stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatUint(uint64(val), 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(buf, 10)
	stream.WriteString("a")
	stream.WriteUint64(0xffffffff) // should clear buffer
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a4294967295", buf.String())
}

func Test_write_int64(t *testing.T) {
	vals := []int64{0, 1, 11, 111, 255, 999999, 0xfff, 0xffff, 0xfffff, 0xffffff, 0xfffffff, 0xffffffff,
		0xfffffffff,0xffffffffff,0xfffffffffff,0xffffffffffff,0xfffffffffffff,0xffffffffffffff,
		0xfffffffffffffff,0x7fffffffffffffff,-0x7fffffffffffffff}
	for _, val := range vals {
		t.Run(fmt.Sprintf("%v", val), func(t *testing.T) {
			should := require.New(t)
			buf := &bytes.Buffer{}
			stream := NewStream(buf, 4096)
			stream.WriteInt64(val)
			stream.Flush()
			should.Nil(stream.Error)
			should.Equal(strconv.FormatInt(val, 10), buf.String())
		})
	}
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(buf, 10)
	stream.WriteString("a")
	stream.WriteInt64(0xffffffff) // should clear buffer
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("a4294967295", buf.String())
}

func Benchmark_jsoniter_encode_int(b *testing.B) {
	stream := NewStream(ioutil.Discard, 64)
	for n := 0; n < b.N; n++ {
		stream.n = 0
		stream.WriteUint64(0xffffffff)
	}
}

func Benchmark_itoa(b *testing.B) {
	for n := 0; n < b.N; n++ {
		strconv.FormatInt(0xffffffff, 10)
	}
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
