package jsoniter

import (
	"testing"
)

func Test_null(t *testing.T) {
	iter := ParseString(`null`)
	if iter.ReadNull() != true {
		t.FailNow()
	}
}

func Test_null_object(t *testing.T) {
	iter := ParseString(`[null,"a"]`)
	iter.ReadArray()
	if iter.ReadObject() != "" {
		t.FailNow()
	}
	iter.ReadArray()
	if iter.ReadString() != "a" {
		t.FailNow()
	}
}

func Test_null_array(t *testing.T) {
	iter := ParseString(`[null,"a"]`)
	iter.ReadArray()
	if iter.ReadArray() != false {
		t.FailNow()
	}
	iter.ReadArray()
	if iter.ReadString() != "a" {
		t.FailNow()
	}
}

func Test_null_string(t *testing.T) {
	iter := ParseString(`[null,"a"]`)
	iter.ReadArray()
	if iter.ReadString() != "" {
		t.FailNow()
	}
	iter.ReadArray()
	if iter.ReadString() != "a" {
		t.FailNow()
	}
}

func Test_null_skip(t *testing.T) {
	iter := ParseString(`[null,"a"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "a" {
		t.FailNow()
	}
}