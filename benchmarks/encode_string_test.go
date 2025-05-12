package test

import (
	"bytes"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func Benchmark_encode_string_with_SetEscapeHTML(b *testing.B) {
	type V struct {
		S string
		B bool
		I int
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b.ReportAllocs()
	for b.Loop() {
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(true)
		if err := enc.Encode(V{S: "s", B: true, I: 233}); err != nil {
			b.Fatal(err)
		}
	}
}
