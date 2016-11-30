package jsoniter

import (
	"testing"
	"bytes"
	"io"
)

func Test_read_by_one(t *testing.T) {
	lexer := NewLexer(bytes.NewBufferString("abc"), 1)
	b, err := lexer.readByte()
	if err != nil {
		t.Fatal(err)
	}
	if b != 'a' {
		t.Fatal(b)
	}
	err = lexer.unreadByte()
	if err != nil {
		t.Fatal(err)
	}
	err = lexer.unreadByte()
	if err == nil {
		t.FailNow()
	}
	b, err = lexer.readByte()
	if err != nil {
		t.Fatal(err)
	}
	if b != 'a' {
		t.Fatal(b)
	}
}

func Test_read_by_two(t *testing.T) {
	lexer := NewLexer(bytes.NewBufferString("abc"), 2)
	b, err := lexer.readByte()
	if err != nil {
		t.Fatal(err)
	}
	if b != 'a' {
		t.Fatal(b)
	}
	b, err = lexer.readByte()
	if err != nil {
		t.Fatal(err)
	}
	if b != 'b' {
		t.Fatal(b)
	}
	err = lexer.unreadByte()
	if err != nil {
		t.Fatal(err)
	}
	err = lexer.unreadByte()
	if err != nil {
		t.Fatal(err)
	}
	b, err = lexer.readByte()
	if err != nil {
		t.Fatal(err)
	}
	if b != 'a' {
		t.Fatal(b)
	}
}

func Test_read_until_eof(t *testing.T) {
	lexer := NewLexer(bytes.NewBufferString("abc"), 2)
	lexer.readByte()
	lexer.readByte()
	b, err := lexer.readByte()
	if err != nil {
		t.Fatal(err)
	}
	if b != 'c' {
		t.Fatal(b)
	}
	_, err = lexer.readByte()
	if err != io.EOF {
		t.Fatal(err)
	}
}