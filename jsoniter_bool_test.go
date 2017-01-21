package jsoniter

import (
	"testing"
	"bytes"
	"github.com/json-iterator/go/require"
)

func Test_true(t *testing.T) {
	should := require.New(t)
	iter := ParseString(`true`)
	should.True(iter.ReadBool())
	iter = ParseString(`true`)
	should.Equal(true, iter.Read())
}

func Test_false(t *testing.T) {
	should := require.New(t)
	iter := ParseString(`false`)
	should.False(iter.ReadBool())
}

func Test_write_true_false(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(buf, 4096)
	stream.WriteTrue()
	stream.WriteFalse()
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("truefalse", buf.String())
}


func Test_write_val_bool(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(buf, 4096)
	stream.WriteVal(true)
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("true", buf.String())
}