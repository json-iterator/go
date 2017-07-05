package jsoniter

import (
	"testing"

	"github.com/json-iterator/go/require"
)

func Test_wrap_and_valuetype_everything(t *testing.T) {
	should := require.New(t)
	any := Wrap(int8(10))
	should.Equal(any.ValueType(), Number)
	any = Wrap(int16(10))
	should.Equal(any.ValueType(), Number)
	any = Wrap(int32(10))
	should.Equal(any.ValueType(), Number)
	any = Wrap(int64(10))
	should.Equal(any.ValueType(), Number)

	any = Wrap(uint(10))
	should.Equal(any.ValueType(), Number)
	any = Wrap(uint8(10))
	should.Equal(any.ValueType(), Number)
	any = Wrap(uint16(10))
	should.Equal(any.ValueType(), Number)
	any = Wrap(uint32(10))
	should.Equal(any.ValueType(), Number)
	any = Wrap(uint64(10))
	should.Equal(any.ValueType(), Number)

	any = Wrap(float32(10))
	should.Equal(any.ValueType(), Number)
	any = Wrap(float64(10))
	should.Equal(any.ValueType(), Number)

	any = Wrap(true)
	should.Equal(any.ValueType(), Bool)
	any = Wrap(false)
	should.Equal(any.ValueType(), Bool)

}
