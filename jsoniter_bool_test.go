package jsoniter

import "testing"

func Test_true(t *testing.T) {
	iter := ParseString(`true`)
	if iter.ReadBool() != true {
		t.FailNow()
	}
}

func Test_false(t *testing.T) {
	iter := ParseString(`false`)
	if iter.ReadBool() != false {
		t.FailNow()
	}
}
