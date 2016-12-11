package jsoniter

import (
	"testing"
	"fmt"
)

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

func Test_read_any_from_nested(t *testing.T) {
	iter := ParseString(`{"numbers": ["1", "2", ["3", "4"]]}`)
	val := iter.ReadAny()
	if val.ToInt("numbers", 2, 0) != 3 {
		fmt.Println(val.Error)
		t.FailNow()
	}
}
