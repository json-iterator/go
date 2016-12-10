package jsoniter

import (
	"testing"
	"fmt"
)

func Test_get_from_map(t *testing.T) {
	any := Any{Val: map[string]interface{}{
		"hello": "world",
	}}
	if any.GetString("hello") != "world" {
		t.FailNow()
	}
}

func Test_get_from_array(t *testing.T) {
	any := Any{Val: []interface{}{
		"hello", "world",
	}}
	if any.GetString(1) != "world" {
		t.FailNow()
	}
}

func Test_get_int(t *testing.T) {
	any := Any{Val: []interface{}{
		1, 2, 3,
	}}
	if any.GetInt(1) != 2 {
		t.FailNow()
	}
}

func Test_is_null(t *testing.T) {
	any := Any{Val: []interface{}{
		1, 2, 3,
	}}
	if any.IsNull() != false {
		t.FailNow()
	}
}

func Test_get_bool(t *testing.T) {
	any := Any{Val: []interface{}{
		true, true, false,
	}}
	if any.GetBool(1) != true {
		t.FailNow()
	}
}

func Test_nested_read(t *testing.T) {
	any := Any{Val: []interface{}{
		true, map[string]interface{}{
			"hello": "world",
		}, false,
	}}
	if any.GetString(1, "hello") != "world" {
		fmt.Println(any.Error)
		t.FailNow()
	}
}
