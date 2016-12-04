package jsoniter

import (
	"testing"
)

func Test_reflect_str(t *testing.T) {
	iter := ParseString(`"hello"`)
	str := ""
	iter.Read(&str)
	if str != "hello" {
		t.FailNow()
	}
}
