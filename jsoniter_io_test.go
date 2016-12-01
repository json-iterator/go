package jsoniter

import (
	"testing"
	"bytes"
	"io"
)

func Test_read_by_one(t *testing.T) {
	iter := Parse(bytes.NewBufferString("abc"), 1)
	b := iter.readByte()
	if iter.Error != nil {
		t.Fatal(iter.Error)
	}
	if b != 'a' {
		t.Fatal(b)
	}
	iter.unreadByte()
	if iter.Error != nil {
		t.Fatal(iter.Error)
	}
	iter.unreadByte()
	if iter.Error == nil {
		t.FailNow()
	}
	iter.Error = nil
	b = iter.readByte()
	if iter.Error != nil {
		t.Fatal(iter.Error)
	}
	if b != 'a' {
		t.Fatal(b)
	}
}

func Test_read_by_two(t *testing.T) {
	iter := Parse(bytes.NewBufferString("abc"), 2)
	b := iter.readByte()
	if iter.Error != nil {
		t.Fatal(iter.Error)
	}
	if b != 'a' {
		t.Fatal(b)
	}
	b = iter.readByte()
	if iter.Error != nil {
		t.Fatal(iter.Error)
	}
	if b != 'b' {
		t.Fatal(b)
	}
	iter.unreadByte()
	if iter.Error != nil {
		t.Fatal(iter.Error)
	}
	iter.unreadByte()
	if iter.Error != nil {
		t.Fatal(iter.Error)
	}
	b = iter.readByte()
	if iter.Error != nil {
		t.Fatal(iter.Error)
	}
	if b != 'a' {
		t.Fatal(b)
	}
}

func Test_read_until_eof(t *testing.T) {
	iter := Parse(bytes.NewBufferString("abc"), 2)
	iter.readByte()
	iter.readByte()
	b := iter.readByte()
	if iter.Error != nil {
		t.Fatal(iter.Error)
	}
	if b != 'c' {
		t.Fatal(b)
	}
	iter.readByte()
	if iter.Error != io.EOF {
		t.Fatal(iter.Error)
	}
}