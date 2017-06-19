package jsoniter

import (
	"github.com/json-iterator/go/require"
	"testing"
)

func Test_read_bool_as_any(t *testing.T) {
	should := require.New(t)
	any := Get([]byte("true"))
	should.True(any.ToBool())
}
