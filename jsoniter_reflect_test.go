package jsoniter

import (
	"encoding/json"
	"fmt"
	"github.com/json-iterator/go/require"
	"testing"
	"unsafe"
)

func Test_decode_slice(t *testing.T) {
	should := require.New(t)
	slice := make([]string, 0, 5)
	UnmarshalFromString(`["hello", "world"]`, &slice)
	should.Equal([]string{"hello", "world"}, slice)
}

func Test_decode_large_slice(t *testing.T) {
	should := require.New(t)
	slice := make([]int, 0, 1)
	UnmarshalFromString(`[1,2,3,4,5,6,7,8,9]`, &slice)
	should.Equal([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, slice)
}

func Test_decode_nested(t *testing.T) {
	type StructOfString struct {
		Field1 string
		Field2 string
	}
	iter := ParseString(ConfigDefault, `[{"field1": "hello"}, null, {"field2": "world"}]`)
	slice := []*StructOfString{}
	iter.ReadVal(&slice)
	if len(slice) != 3 {
		fmt.Println(iter.Error)
		t.Fatal(len(slice))
	}
	if slice[0].Field1 != "hello" {
		fmt.Println(iter.Error)
		t.Fatal(slice[0])
	}
	if slice[1] != nil {
		fmt.Println(iter.Error)
		t.Fatal(slice[1])
	}
	if slice[2].Field2 != "world" {
		fmt.Println(iter.Error)
		t.Fatal(slice[2])
	}
}

func Test_decode_base64(t *testing.T) {
	iter := ParseString(ConfigDefault, `"YWJj"`)
	val := []byte{}
	RegisterTypeDecoder("[]uint8", func(ptr unsafe.Pointer, iter *Iterator) {
		*((*[]byte)(ptr)) = iter.ReadBase64()
	})
	defer ConfigDefault.CleanDecoders()
	iter.ReadVal(&val)
	if "abc" != string(val) {
		t.Fatal(string(val))
	}
}

type StructOfTagOne struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
	Field3 int    `json:"field3,string"`
	Field4 int    `json:"field4,string"`
}

func Benchmark_jsoniter_reflect(b *testing.B) {
	b.ReportAllocs()
	iter := NewIterator(ConfigDefault)
	Struct := &StructOfTagOne{}
	//var Struct *StructOfTagOne
	input := []byte(`{"field3": "100", "field4": "100"}`)
	//input := []byte(`null`)
	for n := 0; n < b.N; n++ {
		iter.ResetBytes(input)
		iter.ReadVal(&Struct)
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
		//		struct_.Field1 = iter.ReadString()
		//	case "field2":
		//		struct_.Field2 = iter.ReadString()
		//	default:
		//		iter.Skip()
		//	}
		//}
		iter := ParseString(ConfigDefault, `["hello", "world"]`)
		array := make([]string, 0, 2)
		for iter.ReadArray() {
			array = append(array, iter.ReadString())
		}
	}
}

func Benchmark_json_reflect(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		Struct := StructOfTagOne{}
		json.Unmarshal([]byte(`{"field3": "100"}`), &Struct)
		//array := make([]string, 0, 2)
		//json.Unmarshal([]byte(`["hello", "world"]`), &array)
	}
}
