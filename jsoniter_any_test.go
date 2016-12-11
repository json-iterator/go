package jsoniter

import "testing"

func Test_read_string_as_any(t *testing.T) {
	iter := ParseString(`[1, {"hello": "world"}, 2]`)
	any := iter.ReadAny()
	if any.ToString(1, "hello") != "world" {
		t.FailNow()
	}
}

func Test_read_float64_as_any(t *testing.T) {
	iter := ParseString(`1.23`)
	any := iter.ReadAny()
	if any.ToFloat32() != 1.23 {
		t.FailNow()
	}
}

func Test_read_int_as_any(t *testing.T) {
	iter := ParseString(`123`)
	any := iter.ReadAny()
	if any.ToFloat32() != 123 {
		t.FailNow()
	}
}
