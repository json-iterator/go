package jsoniter

import (
	"testing"
	"bytes"
	"encoding/json"
)

func Test_string_empty(t *testing.T) {
	iter := Parse(bytes.NewBufferString(`""`), 4096)
	val := iter.ReadString()
	if iter.Error != nil {
		t.Fatal(iter.Error)
	}
	if val != "" {
		t.Fatal(val)
	}
}

func Test_string_hello(t *testing.T) {
	iter := Parse(bytes.NewBufferString(`"hello"`), 4096)
	val := iter.ReadString()
	if iter.Error != nil {
		t.Fatal(iter.Error)
	}
	if val != "hello" {
		t.Fatal(val)
	}
}

func Test_string_escape_quote(t *testing.T) {
	iter := Parse(bytes.NewBufferString(`"hel\"lo"`), 4096)
	val := iter.ReadString()
	if iter.Error != nil {
		t.Fatal(iter.Error)
	}
	if val != `hel"lo` {
		t.Fatal(val)
	}
}

func Test_string_escape_newline(t *testing.T) {
	iter := Parse(bytes.NewBufferString(`"hel\nlo"`), 4096)
	val := iter.ReadString()
	if iter.Error != nil {
		t.Fatal(iter.Error)
	}
	if val != "hel\nlo" {
		t.Fatal(val)
	}
}

func Test_string_escape_unicode(t *testing.T) {
	iter := Parse(bytes.NewBufferString(`"\u4e2d\u6587"`), 4096)
	val := iter.ReadString()
	if iter.Error != nil {
		t.Fatal(iter.Error)
	}
	if val != "中文" {
		t.Fatal(val)
	}
}

func Test_string_escape_unicode_with_surrogate(t *testing.T) {
	iter := Parse(bytes.NewBufferString(`"\ud83d\udc4a"`), 4096)
	val := iter.ReadString()
	if iter.Error != nil {
		t.Fatal(iter.Error)
	}
	if val != "\xf0\x9f\x91\x8a" {
		t.Fatal(val)
	}
}

func Test_string_as_bytes(t *testing.T) {
	iter := Parse(bytes.NewBufferString(`"hello""world"`), 4096)
	val := string(iter.readStringAsBytes())
	if val != "hello" {
		t.Fatal(val)
	}
	val = string(iter.readStringAsBytes())
	if val != "world" {
		t.Fatal(val)
	}
}

func Benchmark_jsoniter_unicode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		iter := ParseString(`"\ud83d\udc4a"`)
		iter.ReadString()
	}
}

func Benchmark_jsoniter_ascii(b *testing.B) {
	iter := ParseString(`"hello, world!"`)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		iter.ResetBytes(iter.buf)
		iter.ReadString()
	}
}

func Benchmark_jsoniter_string_as_bytes(b *testing.B) {
	iter := ParseString(`"hello, world!"`)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		iter.ResetBytes(iter.buf)
		iter.readStringAsBytes()
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