package jsoniter

import (
	"testing"
	"github.com/json-iterator/go/require"
	"bytes"
)

func Test_read_null(t *testing.T) {
	should := require.New(t)
	iter := ParseString(`null`)
	should.True(iter.ReadNil())
	iter = ParseString(`null`)
	should.Nil(iter.Read())
}

func Test_write_null(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(buf, 4096)
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
	iter := ParseString(`[null,"a"]`)
	iter.ReadArray()
	if iter.ReadObject() != "" {
		t.FailNow()
	}
	iter.ReadArray()
	if iter.ReadString() != "a" {
		t.FailNow()
	}
}

func Test_decode_null_array(t *testing.T) {
	iter := ParseString(`[null,"a"]`)
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
	iter := ParseString(`[null,"a"]`)
	should.True(iter.ReadArray())
	should.True(iter.ReadNil())
	should.True(iter.ReadArray())
	should.Equal("a", iter.ReadString())
}

func Test_decode_null_skip(t *testing.T) {
	iter := ParseString(`[null,"a"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "a" {
		t.FailNow()
	}
}
