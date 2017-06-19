package jsoniter

import (
	"github.com/json-iterator/go/require"
	"testing"
)

func Test_read_float_as_any(t *testing.T) {
	should := require.New(t)
	any := Get([]byte("12.3"))
	should.Equal(float64(12.3), any.ToFloat64())
	should.Equal("12.3", any.ToString())
	should.True(any.ToBool())
}
