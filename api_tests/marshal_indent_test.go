package test

import (
	"encoding/json"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

func Test_marshal_indent(t *testing.T) {
	testcases := []struct {
		name   string
		obj    interface{}
		expect string
	}{
		{
			name: "no writable fields",
			obj: struct {
				Writable   bool `json:"writable,omitempty"`
				unexported int
			}{Writable: false, unexported: 1},
			expect: "{}",
		},
		{
			name: "flattened fields",
			obj: struct {
				F1 int
				F2 []int
			}{F1: 1, F2: []int{2, 3, 4}},
			expect: "{\n  \"F1\": 1,\n  \"F2\": [\n    2,\n    3,\n    4\n  ]\n}",
		},
		{
			name: "nested fields",
			obj: struct {
				F1 map[int]int
				F2 struct{ V int }
			}{F1: map[int]int{1: 1}, F2: struct{ V int }{V: 2}},
			expect: "{\n  \"F1\": {\n    \"1\": 1\n  },\n  \"F2\": {\n    \"V\": 2\n  }\n}",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			should := require.New(t)
			obj := tc.obj

			output, err := json.MarshalIndent(obj, "", "  ")
			should.Nil(err)
			should.Equal(tc.expect, string(output))

			output, err = jsoniter.MarshalIndent(obj, "", "  ")
			should.Nil(err)
			should.Equal(tc.expect, string(output))

			output, err = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalIndent(obj, "", "  ")
			should.Nil(err)
			should.Equal(tc.expect, string(output))
		})
	}
}

func Test_marshal_indent_map(t *testing.T) {
	testcases := []struct {
		name   string
		obj    interface{}
		expect string
	}{
		{
			name:   "empty map",
			obj:    map[int]int{},
			expect: "{}",
		},
		{
			name:   "map with literal value",
			obj:    map[int]int{1: 2},
			expect: "{\n  \"1\": 2\n}",
		},
		{
			name:   "map with array value",
			obj:    map[int][]int{1: {1, 2}},
			expect: "{\n  \"1\": [\n    1,\n    2\n  ]\n}",
		},
		{
			name:   "map with nested map",
			obj:    map[int]map[int]int{1: {1: 2}},
			expect: "{\n  \"1\": {\n    \"1\": 2\n  }\n}",
		},
		{
			name:   "map with object value",
			obj:    map[int]struct{ F int }{1: {F: 2}},
			expect: "{\n  \"1\": {\n    \"F\": 2\n  }\n}",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			should := require.New(t)
			obj := tc.obj

			output, err := json.MarshalIndent(obj, "", "  ")
			should.Nil(err)
			should.Equal(tc.expect, string(output))

			output, err = jsoniter.MarshalIndent(obj, "", "  ")
			should.Nil(err)
			should.Equal(tc.expect, string(output))

			output, err = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalIndent(obj, "", "  ")
			should.Nil(err)
			should.Equal(tc.expect, string(output))
		})
	}
}
