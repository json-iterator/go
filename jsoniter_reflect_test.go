package jsoniter

import (
	"testing"
	"fmt"
	"encoding/json"
)

func Test_reflect_str(t *testing.T) {
	iter := ParseString(`"hello"`)
	str := ""
	iter.Read(&str)
	if str != "hello" {
		fmt.Println(iter.Error)
		t.Fatal(str)
	}
}

func Test_reflect_ptr_str(t *testing.T) {
	iter := ParseString(`"hello"`)
	var str *string
	iter.Read(&str)
	if *str != "hello" {
		t.Fatal(str)
	}
}

type StructOfString struct {
	field1 string
	field2 string
}

func Test_reflect_struct_string(t *testing.T) {
	iter := ParseString(`{"field1": "hello", "field2": "world"}`)
	struct_ := StructOfString{}
	iter.Read(&struct_)
	if struct_.field1 != "hello" {
		fmt.Println(iter.Error)
		t.Fatal(struct_.field1)
	}
	if struct_.field2 != "world" {
		fmt.Println(iter.Error)
		t.Fatal(struct_.field1)
	}
}

type StructOfStringPtr struct {
	field1 *string
	field2 *string
}

func Test_reflect_struct_string_ptr(t *testing.T) {
	iter := ParseString(`{"field1": null, "field2": "world"}`)
	struct_ := StructOfStringPtr{}
	iter.Read(&struct_)
	if struct_.field1 != nil {
		fmt.Println(iter.Error)
		t.Fatal(struct_.field1)
	}
	if *struct_.field2 != "world" {
		fmt.Println(iter.Error)
		t.Fatal(struct_.field1)
	}
}

func Test_reflect_slice(t *testing.T) {
	iter := ParseString(`["hello", "world"]`)
	array := make([]string, 0, 1)
	iter.Read(&array)
	if len(array) != 2 {
		fmt.Println(iter.Error)
		t.Fatal(len(array))
	}
	if array[0] != "hello" {
		fmt.Println(iter.Error)
		t.Fatal(array[0])
	}
	if array[1] != "world" {
		fmt.Println(iter.Error)
		t.Fatal(array[1])
	}
}

func Benchmark_jsoniter_reflect(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		//iter := ParseString(`{"field1": "hello", "field2": "world"}`)
		//struct_ := StructOfString{}
		//iter.Read(&struct_)
		iter := ParseString(`["hello", "world"]`)
		array := make([]string, 0, 1)
		iter.Read(&array)
	}
}

func Benchmark_jsoniter_direct(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		//iter := ParseString(`{"field1": "hello", "field2": "world"}`)
		//struct_ := StructOfString{}
		//for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
		//	switch field {
		//	case "field1":
		//		struct_.field1 = iter.ReadString()
		//	case "field2":
		//		struct_.field2 = iter.ReadString()
		//	default:
		//		iter.Skip()
		//	}
		//}
		iter := ParseString(`["hello", "world"]`)
		array := make([]string, 0, 2)
		for iter.ReadArray() {
			array = append(array, iter.ReadString())
		}
	}
}

func Benchmark_json_reflect(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		//struct_ := StructOfString{}
		//json.Unmarshal([]byte(`{"field1": "hello", "field2": "world"}`), &struct_)
		array := make([]string, 0, 2)
		json.Unmarshal([]byte(`["hello", "world"]`), &array)
	}
}