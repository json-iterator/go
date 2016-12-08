package jsoniter

import (
	"testing"
	"fmt"
	"encoding/json"
	"unsafe"
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

func Test_reflect_int(t *testing.T) {
	iter := ParseString(`123`)
	val := int(0)
	iter.Read(&val)
	if val != 123 {
		t.Fatal(val)
	}
}

func Test_reflect_int8(t *testing.T) {
	iter := ParseString(`123`)
	val := int8(0)
	iter.Read(&val)
	if val != 123 {
		t.Fatal(val)
	}
}

func Test_reflect_int16(t *testing.T) {
	iter := ParseString(`123`)
	val := int16(0)
	iter.Read(&val)
	if val != 123 {
		t.Fatal(val)
	}
}

func Test_reflect_int32(t *testing.T) {
	iter := ParseString(`123`)
	val := int32(0)
	iter.Read(&val)
	if val != 123 {
		t.Fatal(val)
	}
}

func Test_reflect_int64(t *testing.T) {
	iter := ParseString(`123`)
	val := int64(0)
	iter.Read(&val)
	if val != 123 {
		t.Fatal(val)
	}
}

func Test_reflect_uint(t *testing.T) {
	iter := ParseString(`123`)
	val := uint(0)
	iter.Read(&val)
	if val != 123 {
		t.Fatal(val)
	}
}

func Test_reflect_uint8(t *testing.T) {
	iter := ParseString(`123`)
	val := uint8(0)
	iter.Read(&val)
	if val != 123 {
		t.Fatal(val)
	}
}

func Test_reflect_uint16(t *testing.T) {
	iter := ParseString(`123`)
	val := uint16(0)
	iter.Read(&val)
	if val != 123 {
		t.Fatal(val)
	}
}

func Test_reflect_uint32(t *testing.T) {
	iter := ParseString(`123`)
	val := uint32(0)
	iter.Read(&val)
	if val != 123 {
		t.Fatal(val)
	}
}

func Test_reflect_uint64(t *testing.T) {
	iter := ParseString(`123`)
	val := uint64(0)
	iter.Read(&val)
	if val != 123 {
		t.Fatal(val)
	}
}

func Test_reflect_byte(t *testing.T) {
	iter := ParseString(`123`)
	val := byte(0)
	iter.Read(&val)
	if val != 123 {
		t.Fatal(val)
	}
}

func Test_reflect_float32(t *testing.T) {
	iter := ParseString(`1.23`)
	val := float32(0)
	iter.Read(&val)
	if val != 1.23 {
		fmt.Println(iter.Error)
		t.Fatal(val)
	}
}

func Test_reflect_float64(t *testing.T) {
	iter := ParseString(`1.23`)
	val := float64(0)
	iter.Read(&val)
	if val != 1.23 {
		fmt.Println(iter.Error)
		t.Fatal(val)
	}
}

func Test_reflect_bool(t *testing.T) {
	iter := ParseString(`true`)
	val := false
	iter.Read(&val)
	if val != true {
		fmt.Println(iter.Error)
		t.Fatal(val)
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
		t.Fatal(struct_.field2)
	}
}

type StructOfTag struct {
	field1 string `json:"field-1"`
	field2 string `json:"-"`
	field3 int `json:",string"`
}

func Test_reflect_struct_tag_field(t *testing.T) {
	iter := ParseString(`{"field-1": "hello", "field2": "", "field3": "100"}`)
	struct_ := StructOfTag{field2: "world"}
	iter.Read(&struct_)
	if struct_.field1 != "hello" {
		fmt.Println(iter.Error)
		t.Fatal(struct_.field1)
	}
	if struct_.field2 != "world" {
		fmt.Println(iter.Error)
		t.Fatal(struct_.field2)
	}
	if struct_.field3 != 100 {
		fmt.Println(iter.Error)
		t.Fatal(struct_.field3)
	}
}

func Test_reflect_slice(t *testing.T) {
	iter := ParseString(`["hello", "world"]`)
	slice := make([]string, 0, 5)
	iter.Read(&slice)
	if len(slice) != 2 {
		fmt.Println(iter.Error)
		t.Fatal(len(slice))
	}
	if slice[0] != "hello" {
		fmt.Println(iter.Error)
		t.Fatal(slice[0])
	}
	if slice[1] != "world" {
		fmt.Println(iter.Error)
		t.Fatal(slice[1])
	}
}

func Test_reflect_large_slice(t *testing.T) {
	iter := ParseString(`[1,2,3,4,5,6,7,8,9]`)
	slice := make([]int, 0, 1)
	iter.Read(&slice)
	if len(slice) != 9 {
		fmt.Println(iter.Error)
		t.Fatal(len(slice))
	}
	if slice[0] != 1 {
		fmt.Println(iter.Error)
		t.Fatal(slice[0])
	}
	if slice[8] != 9 {
		fmt.Println(iter.Error)
		t.Fatal(slice[1])
	}
}

func Test_reflect_nested(t *testing.T) {
	iter := ParseString(`[{"field1": "hello"}, null, {"field2": "world"}]`)
	slice := []*StructOfString{}
	iter.Read(&slice)
	if len(slice) != 3 {
		fmt.Println(iter.Error)
		t.Fatal(len(slice))
	}
	if slice[0].field1 != "hello" {
		fmt.Println(iter.Error)
		t.Fatal(slice[0])
	}
	if slice[1] != nil {
		fmt.Println(iter.Error)
		t.Fatal(slice[1])
	}
	if slice[2].field2 != "world" {
		fmt.Println(iter.Error)
		t.Fatal(slice[1])
	}
}

func Test_reflect_base64(t *testing.T) {
	iter := ParseString(`"YWJj"`)
	val := []byte{}
	RegisterTypeDecoder("[]uint8", func(ptr unsafe.Pointer, iter *Iterator) {
		*((*[]byte)(ptr)) = iter.ReadBase64()
	})
	defer ClearDecoders()
	iter.Read(&val)
	if "abc" != string(val) {
		t.Fatal(string(val))
	}
}

type StructOfTagOne struct {
	field1 string `json:"field1"`
	field2 string `json:"field2"`
	field3 int `json:"field3,string"`
	field4 int `json:"field4,string"`
}

func Benchmark_jsoniter_reflect(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		iter := ParseString(`{"field3": "100"}`)
		struct_ := StructOfTagOne{}
		iter.Read(&struct_)
		//iter := ParseString(`[1,2,3]`)
		//var array []int
		//iter.Read(&array)
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
		struct_ := StructOfTagOne{}
		json.Unmarshal([]byte(`{"field3": "100"}`), &struct_)
		//array := make([]string, 0, 2)
		//json.Unmarshal([]byte(`["hello", "world"]`), &array)
	}
}