package jsoniter

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

//func Test_large_file(t *testing.T) {
//	file, err := os.Open("/tmp/large-file.json")
//	if err != nil {
//		t.Fatal(err)
//	}
//	iter := Parse(file, 4096)
//	count := 0
//	for iter.ReadArray() {
//		iter.Skip()
//		count++
//	}
//	if count != 11351 {
//		t.Fatal(count)
//	}
//}

func Benchmark_jsoniter_large_file(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		file, _ := os.Open("/tmp/large-file.json")
		iter := Parse(ConfigDefault, file, 4096)
		count := 0
		for iter.ReadArray() {
			iter.Skip()
			count++
		}
		file.Close()
	}
}

func Benchmark_json_large_file(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		file, _ := os.Open("/tmp/large-file.json")
		bytes, _ := ioutil.ReadAll(file)
		file.Close()
		result := []struct{}{}
		json.Unmarshal(bytes, &result)
	}
}
