package jsoniter

import (
	"bytes"
	"encoding/json"
	"github.com/json-iterator/go/require"
	"testing"
)

func Test_read_null(t *testing.T) {
	should := require.New(t)
	iter := ParseString(ConfigOfDefault, `null`)
	should.True(iter.ReadNil())
	iter = ParseString(ConfigOfDefault, `null`)
	should.Nil(iter.Read())
	iter = ParseString(ConfigOfDefault, `null`)
	any, err := UnmarshalAnyFromString(`null`)
	should.Nil(err)
	should.Equal(0, any.ToInt())
	should.Equal(float64(0), any.ToFloat64())
	should.Equal("", any.ToString())
	should.False(any.ToBool())
}

func Test_write_null(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(ConfigOfDefault, buf, 4096)
	stream.WriteNil()
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("null", buf.String())
}

func Test_encode_null(t *testing.T) {
	should := require.New(t)
	str, err := MarshalToString(nil)
	should.Nil(err)
	should.Equal("null", str)
}

func Test_decode_null_object(t *testing.T) {
	should := require.New(t)
	iter := ParseString(ConfigOfDefault, `[null,"a"]`)
	iter.ReadArray()
	if iter.ReadObject() != "" {
		t.FailNow()
	}
	iter.ReadArray()
	if iter.ReadString() != "a" {
		t.FailNow()
	}
	type TestObject struct {
		Field string
	}
	objs := []TestObject{}
	should.Nil(UnmarshalFromString("[null]", &objs))
	should.Len(objs, 1)
}

func Test_decode_null_array(t *testing.T) {
	iter := ParseString(ConfigOfDefault, `[null,"a"]`)
	iter.ReadArray()
	if iter.ReadArray() != false {
		t.FailNow()
	}
	iter.ReadArray()
	if iter.ReadString() != "a" {
		t.FailNow()
	}
}

func Test_decode_null_string(t *testing.T) {
	should := require.New(t)
	iter := ParseString(ConfigOfDefault, `[null,"a"]`)
	should.True(iter.ReadArray())
	should.Equal("", iter.ReadString())
	should.True(iter.ReadArray())
	should.Equal("a", iter.ReadString())
}

func Test_decode_null_skip(t *testing.T) {
	iter := ParseString(ConfigOfDefault, `[null,"a"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "a" {
		t.FailNow()
	}
}

func Test_encode_nil_map(t *testing.T) {
	should := require.New(t)
	type Ttest map[string]string
	var obj1 Ttest
	output, err := json.Marshal(obj1)
	should.Nil(err)
	should.Equal("null", string(output))
	output, err = json.Marshal(&obj1)
	should.Nil(err)
	should.Equal("null", string(output))
	output, err = Marshal(obj1)
	should.Nil(err)
	should.Equal("null", string(output))
	output, err = Marshal(&obj1)
	should.Nil(err)
	should.Equal("null", string(output))
}

func Test_encode_nil_array(t *testing.T) {
	should := require.New(t)
	type Ttest []string
	var obj1 Ttest
	output, err := json.Marshal(obj1)
	should.Nil(err)
	should.Equal("null", string(output))
	output, err = json.Marshal(&obj1)
	should.Nil(err)
	should.Equal("null", string(output))
	output, err = Marshal(obj1)
	should.Nil(err)
	should.Equal("null", string(output))
	output, err = Marshal(&obj1)
	should.Nil(err)
	should.Equal("null", string(output))
}
