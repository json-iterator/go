package jsoniter

import (
	"testing"
	"github.com/json-iterator/go/require"
)

func Test_read_float_as_any(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString("12.3")
	should.Nil(err)
	should.Equal(float64(12.3), any.ToFloat64())
	should.Equal("12.3", any.ToString())
	should.True(any.ToBool())
}
