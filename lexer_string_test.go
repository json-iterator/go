package jsoniter

import (
	"testing"
	"bytes"
	"encoding/json"
)

func Test_string_empty(t *testing.T) {
	lexer := NewLexer(bytes.NewBufferString(`""`), 4096)
	val, err := lexer.LexString()
	if err != nil {
		t.Fatal(err)
	}
	if val != "" {
		t.Fatal(val)
	}
}

func Test_string_hello(t *testing.T) {
	lexer := NewLexer(bytes.NewBufferString(`"hello"`), 4096)
	val, err := lexer.LexString()
	if err != nil {
		t.Fatal(err)
	}
	if val != "hello" {
		t.Fatal(val)
	}
}

func Test_string_escape_quote(t *testing.T) {
	lexer := NewLexer(bytes.NewBufferString(`"hel\"lo"`), 4096)
	val, err := lexer.LexString()
	if err != nil {
		t.Fatal(err)
	}
	if val != `hel"lo` {
		t.Fatal(val)
	}
}

func Test_string_escape_newline(t *testing.T) {
	lexer := NewLexer(bytes.NewBufferString(`"hel\nlo"`), 4096)
	val, err := lexer.LexString()
	if err != nil {
		t.Fatal(err)
	}
	if val != "hel\nlo" {
		t.Fatal(val)
	}
}

func Test_string_escape_unicode(t *testing.T) {
	lexer := NewLexer(bytes.NewBufferString(`"\u4e2d\u6587"`), 4096)
	val, err := lexer.LexString()
	if err != nil {
		t.Fatal(err)
	}
	if val != "中文" {
		t.Fatal(val)
	}
}

func Test_string_escape_unicode_with_surrogate(t *testing.T) {
	lexer := NewLexer(bytes.NewBufferString(`"\ud83d\udc4a"`), 4096)
	val, err := lexer.LexString()
	if err != nil {
		t.Fatal(err)
	}
	if val != "\xf0\x9f\x91\x8a" {
		t.Fatal(val)
	}
}

func Benchmark_jsoniter_unicode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		lexer := NewLexerWithArray([]byte(`"\ud83d\udc4a"`))
		lexer.LexString()
	}
}

func Benchmark_jsoniter_ascii(b *testing.B) {
	for n := 0; n < b.N; n++ {
		lexer := NewLexerWithArray([]byte(`"hello"`))
		lexer.LexString()
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