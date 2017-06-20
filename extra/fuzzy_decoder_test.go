package extra

import (
	"testing"
	"github.com/json-iterator/go"
	"github.com/json-iterator/go/require"
)

func init() {
	RegisterFuzzyDecoders()
}

func Test_string_to_string(t *testing.T) {
	should := require.New(t)
	var val string
	should.Nil(jsoniter.UnmarshalFromString(`"100"`, &val))
	should.Equal("100", val)
}

func Test_int_to_string(t *testing.T) {
	should := require.New(t)
	var val string
	should.Nil(jsoniter.UnmarshalFromString(`100`, &val))
	should.Equal("100", val)
}

func Test_float_to_string(t *testing.T) {
	should := require.New(t)
	var val string
	should.Nil(jsoniter.UnmarshalFromString(`12.0`, &val))
	should.Equal("12.0", val)
}

func Test_string_to_int(t *testing.T) {
	should := require.New(t)
	var val int
	should.Nil(jsoniter.UnmarshalFromString(`"100"`, &val))
	should.Equal(100, val)
}

func Test_int_to_int(t *testing.T) {
	should := require.New(t)
	var val int
	should.Nil(jsoniter.UnmarshalFromString(`100`, &val))
	should.Equal(100, val)
}
