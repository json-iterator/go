package jsoniter

import (
	"testing"

	"github.com/json-iterator/go/require"
)

func Test_invalid(t *testing.T) {
	should := require.New(t)
	any := Get([]byte("[]"))
	should.Equal(Invalid, any.Get(0.3).ValueType())
	// is nil correct ?
	should.Equal(nil, any.Get(0.3).GetInterface())

	any = any.Get(0.3)
	should.Equal(false, any.ToBool())
	should.Equal(int(0), any.ToInt())
	should.Equal(int32(0), any.ToInt32())
	should.Equal(int64(0), any.ToInt64())
	should.Equal(uint(0), any.ToUint())
	should.Equal(uint32(0), any.ToUint32())
	should.Equal(uint64(0), any.ToUint64())
	should.Equal(float32(0), any.ToFloat32())
	should.Equal(float64(0), any.ToFloat64())
	should.Equal("", any.ToString())

	should.Equal(Invalid, any.Get(0.1).Get(1).ValueType())

}
