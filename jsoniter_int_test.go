package jsoniter

import (
	"testing"
	"bytes"
	"encoding/json"
)

func Test_uint64_0(t *testing.T) {
	iter := Parse(bytes.NewBufferString("0"), 4096)
	val := iter.ReadUint64()
	if iter.Error != nil {
		t.Fatal(iter.Error)
	}
	if val != 0 {
		t.Fatal(val)
	}
}

func Test_uint64_1(t *testing.T) {
	iter := Parse(bytes.NewBufferString("1"), 4096)
	val := iter.ReadUint64()
	if val != 1 {
		t.Fatal(val)
	}
}

func Test_uint64_100(t *testing.T) {
	iter := Parse(bytes.NewBufferString("100"), 4096)
	val := iter.ReadUint64()
	if val != 100 {
		t.Fatal(val)
	}
}

func Test_uint64_100_comma(t *testing.T) {
	iter := Parse(bytes.NewBufferString("100,"), 4096)
	val := iter.ReadUint64()
	if iter.Error != nil {
		t.Fatal(iter.Error)
	}
	if val != 100 {
		t.Fatal(val)
	}
}

func Test_uint64_invalid(t *testing.T) {
	iter := Parse(bytes.NewBufferString(","), 4096)
	iter.ReadUint64()
	if iter.Error == nil {
		t.FailNow()
	}
}

func Test_int64_100(t *testing.T) {
	iter := Parse(bytes.NewBufferString("100"), 4096)
	val := iter.ReadInt64()
	if val != 100 {
		t.Fatal(val)
	}
}

func Test_int64_minus_100(t *testing.T) {
	iter := Parse(bytes.NewBufferString("-100"), 4096)
	val := iter.ReadInt64()
	if val != -100 {
		t.Fatal(val)
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