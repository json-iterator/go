package test

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func BenchmarkFloatUnmarshal(b *testing.B) {
	type floaty struct {
		A float32
		B float64
	}

	data, err := jsoniter.Marshal(floaty{A: 1.111111111111111, B: 1.11111111111111})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()
	var out floaty
	for i := 0; i < b.N; i++ {
		if err := jsoniter.Unmarshal(data, &out); err != nil {
			b.Fatal(err)
		}
	}
}
