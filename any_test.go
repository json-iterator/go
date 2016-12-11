package jsoniter

import (
	"testing"
	"fmt"
)

func Test_get_from_map(t *testing.T) {
	any := Any{val: map[string]interface{}{
		"hello": "world",
	}}
	if any.ToString("hello") != "world" {
		t.FailNow()
	}
}

func Test_get_from_array(t *testing.T) {
	any := Any{val: []interface{}{
		"hello", "world",
	}}
	if any.ToString(1) != "world" {
		t.FailNow()
	}
}

func Test_get_int(t *testing.T) {
	any := Any{val: []interface{}{
		1, 2, 3,
	}}
	if any.ToInt(1) != 2 {
		t.FailNow()
	}
}

func Test_is_null(t *testing.T) {
	any := Any{val: []interface{}{
		1, 2, 3,
	}}
	if any.IsNull() != false {
		t.FailNow()
	}
}

func Test_get_bool(t *testing.T) {
	any := Any{val: []interface{}{
		true, true, false,
	}}
	if any.ToBool(1) != true {
		t.FailNow()
	}
}

func Test_nested_read(t *testing.T) {
	any := Any{val: []interface{}{
		true, map[string]interface{}{
			"hello": "world",
		}, false,
	}}
	if any.ToString(1, "hello") != "world" {
		fmt.Println(any.Error)
		t.FailNow()
	}
}

func Test_int_to_string(t *testing.T) {
	any := Any{val: []interface{}{
		true, 5, false,
	}}
	if any.ToString(1) != "5" {
		t.FailNow()
	}
}
