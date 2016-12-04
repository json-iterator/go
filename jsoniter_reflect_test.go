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
		t.Fatal(str)
	}
}

type StructOfString struct {
	field1 string
	field2 string
}

func Test_reflect_struct(t *testing.T) {
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

func Benchmark_jsoniter_reflect(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		iter := ParseString(`{"field1": "hello", "field2": "world"}`)
		struct_ := StructOfString{}
		iter.Read(&struct_)
	}
}

func Benchmark_jsoniter_direct(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		iter := ParseString(`{"field1": "hello", "field2": "world"}`)
		struct_ := StructOfString{}
		for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
			switch field {
			case "field1":
				struct_.field1 = iter.ReadString()
			case "field2":
				struct_.field2 = iter.ReadString()
			default:
				iter.Skip()
			}
		}
	}
}

func Benchmark_json_reflect(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		struct_ := StructOfString{}
		json.Unmarshal([]byte(`{"field1": "hello", "field2": "world"}`), &struct_)
	}
}