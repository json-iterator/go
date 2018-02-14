package test

import (
	"time"
	"encoding/json"
)

func init() {
	unmarshalCases = append(unmarshalCases, unmarshalCase{
		ptr: (*struct {
			Field interface{}
		})(nil),
		input: `{"Field": "hello"}`,
	}, unmarshalCase{
		ptr: (*struct {
			Field int `json:"field"`
		})(nil),
		input: `{"field": null}`,
	})
	marshalCases = append(marshalCases,
		struct {
			Field map[string]interface{}
		}{
			map[string]interface{}{"hello": "world"},
		},
		struct {
			Field  map[string]interface{}
			Field2 string
		}{
			map[string]interface{}{"hello": "world"}, "",
		},
		struct {
			Field interface{}
		}{
			1024,
		},
		struct {
			Field MyInterface
		}{
			MyString("hello"),
		},
		struct {
			F *float64
		}{},
		// TODO: fix this
		//struct {
		//	*time.Time
		//}{},
		struct {
			*time.Time
		}{&time.Time{}},
		struct {
			*StructVarious
		}{&StructVarious{}},
		struct {
			*StructVarious
		}{},
		struct {
			Field1 int
			Field2 [1]*float64
		}{},
		struct {
			Field interface{} `json:"field,omitempty"`
		}{},
		struct {
			Field MyInterface `json:"field,omitempty"`
		}{},
		struct {
			Field MyInterface `json:"field,omitempty"`
		}{MyString("hello")},
		struct {
			Field json.Marshaler `json:"field"`
		}{},
		struct {
			Field MyInterface `json:"field"`
		}{},
		struct {
			Field MyInterface `json:"field"`
		}{MyString("hello")},
	)
}

type StructVarious struct {
	Field0 string
	Field1 []string
	Field2 map[string]interface{}
}