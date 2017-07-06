package jsoniter

import (
	"testing"
	"github.com/json-iterator/go/require"
)

func Test_missing_object_end(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Metric string                 `json:"metric"`
		Tags   map[string]interface{} `json:"tags"`
	}
	obj := TestObject{}
	should.NotNil(UnmarshalFromString(`{"metric": "sys.777","tags": {"a":"123"}`, &obj))
}

func Test_missing_array_end(t *testing.T) {
	should := require.New(t)
	should.NotNil(UnmarshalFromString(`[1,2,3`, &[]int{}))
}
