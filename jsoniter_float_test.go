package jsoniter

import (
	"testing"
	"encoding/json"
)

func Test_float64_0(t *testing.T) {
	iter := ParseString(`0`)
	val := iter.ReadFloat64()
	if val != 0 {
		t.Fatal(val)
	}
}

func Test_float64_1_dot_1(t *testing.T) {
	iter := ParseString(`1.1`)
	val := iter.ReadFloat64()
	if val != 1.1 {
		t.Fatal(val)
	}
}

func Test_float32_1_dot_1_comma(t *testing.T) {
	iter := ParseString(`1.1,`)
	val := iter.ReadFloat32()
	if val != 1.1 {
		t.Fatal(val)
	}
}

func Benchmark_jsoniter_float(b *testing.B) {
	for n := 0; n < b.N; n++ {
		iter := ParseString(`1.1`)
		iter.ReadFloat64()
	}
}

func Benchmark_json_float(b *testing.B) {
	for n := 0; n < b.N; n++ {
		result := float64(0)
		json.Unmarshal([]byte(`1.1`), &result)
	}
}