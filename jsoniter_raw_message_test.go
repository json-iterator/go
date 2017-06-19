package jsoniter

import (
	"testing"
	"encoding/json"
	"github.com/json-iterator/go/require"
)

func Test_json_RawMessage(t *testing.T) {
	should := require.New(t)
	var data json.RawMessage
	should.Nil(Unmarshal([]byte(`[1,2,3]`), &data))
	should.Equal(`[1,2,3]`, string(data))
	str, err := MarshalToString(data)
	should.Nil(err)
	should.Equal(`[1,2,3]`, str)
}

func Test_json_RawMessage_in_struct(t *testing.T) {
	type TestObject struct {
		Field1 string
		Field2 json.RawMessage
	}
	should := require.New(t)
	var data TestObject
	should.Nil(Unmarshal([]byte(`{"field1": "hello", "field2": [1,2,3]}`), &data))
	should.Equal(` [1,2,3]`, string(data.Field2))
	should.Equal(`hello`, data.Field1)
}
