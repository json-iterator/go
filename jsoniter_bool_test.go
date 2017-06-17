package jsoniter

import (
	"bytes"
	"encoding/json"
	"github.com/json-iterator/go/require"
	"testing"
)

func Test_true(t *testing.T) {
	should := require.New(t)
	iter := ParseString(ConfigDefault, `true`)
	should.True(iter.ReadBool())
	iter = ParseString(ConfigDefault, `true`)
	should.Equal(true, iter.Read())
}

func Test_false(t *testing.T) {
	should := require.New(t)
	iter := ParseString(ConfigDefault, `false`)
	should.False(iter.ReadBool())
}

func Test_read_bool_as_any(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString("true")
	should.Nil(err)
	should.True(any.ToBool())
}

func Test_write_true_false(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(ConfigDefault, buf, 4096)
	stream.WriteTrue()
	stream.WriteFalse()
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("truefalse", buf.String())
}

func Test_write_val_bool(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(ConfigDefault, buf, 4096)
	stream.WriteVal(true)
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("true", buf.String())
}

func Test_encode_string_bool(t *testing.T) {
	type TestObject struct {
		Field bool `json:",omitempty,string"`
	}
	should := require.New(t)
	output, err := json.Marshal(TestObject{true})
	should.Nil(err)
	should.Equal(`{"Field":"true"}`, string(output))
	output, err = Marshal(TestObject{true})
	should.Nil(err)
	should.Equal(`{"Field":"true"}`, string(output))
}

func Test_decode_string_bool(t *testing.T) {
	type TestObject struct {
		Field bool `json:",omitempty,string"`
	}
	should := require.New(t)
	obj := TestObject{}
	err := json.Unmarshal([]byte(`{"Field":"true"}`), &obj)
	should.Nil(err)
	should.True(obj.Field)

	obj = TestObject{}
	err = json.Unmarshal([]byte(`{"Field":true}`), &obj)
	should.NotNil(err)

	obj = TestObject{}
	err = Unmarshal([]byte(`{"Field":"true"}`), &obj)
	should.Nil(err)
	should.True(obj.Field)

	obj = TestObject{}
	err = Unmarshal([]byte(`{"Field":true}`), &obj)
	should.NotNil(err)
}
