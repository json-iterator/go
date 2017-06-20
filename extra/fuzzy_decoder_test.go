package extra

import (
	"github.com/json-iterator/go"
	"github.com/json-iterator/go/require"
	"testing"
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

func Test_float_to_int(t *testing.T) {
	should := require.New(t)
	var val int
	should.Nil(jsoniter.UnmarshalFromString(`1.23`, &val))
	should.Equal(1, val)
}

func Test_large_float_to_int(t *testing.T) {
	should := require.New(t)
	var val int
	should.NotNil(jsoniter.UnmarshalFromString(`1234512345123451234512345.0`, &val))
}

func Test_string_to_float32(t *testing.T) {
	should := require.New(t)
	var val float32
	should.Nil(jsoniter.UnmarshalFromString(`"100"`, &val))
	should.Equal(float32(100), val)
}

func Test_float_to_float32(t *testing.T) {
	should := require.New(t)
	var val float32
	should.Nil(jsoniter.UnmarshalFromString(`1.23`, &val))
	should.Equal(float32(1.23), val)
}

func Test_string_to_float64(t *testing.T) {
	should := require.New(t)
	var val float64
	should.Nil(jsoniter.UnmarshalFromString(`"100"`, &val))
	should.Equal(float64(100), val)
}

func Test_float_to_float64(t *testing.T) {
	should := require.New(t)
	var val float64
	should.Nil(jsoniter.UnmarshalFromString(`1.23`, &val))
	should.Equal(float64(1.23), val)
}

func Test_empty_array_as_map(t *testing.T) {
	should := require.New(t)
	var val map[string]interface{}
	should.Nil(jsoniter.UnmarshalFromString(`[]`, &val))
	should.Equal(map[string]interface{}{}, val)
}

func Test_empty_array_as_object(t *testing.T) {
	should := require.New(t)
	var val struct{}
	should.Nil(jsoniter.UnmarshalFromString(`[]`, &val))
	should.Equal(struct{}{}, val)
}
