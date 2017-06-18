package jsoniter

import (
	"testing"
	"github.com/json-iterator/go/require"
)

func Test_read_bool_as_any(t *testing.T) {
	should := require.New(t)
	any := Get([]byte("true"))
	should.True(any.ToBool())
}
