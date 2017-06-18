package jsoniter

import (
	"testing"
	"github.com/json-iterator/go/require"
)

func Test_read_empty_array_as_any(t *testing.T) {
	should := require.New(t)
	any := Get([]byte("[]"))
	should.Equal(Array, any.Get().ValueType())
	should.Equal(Invalid, any.Get(0.3).ValueType())
	should.Equal(0, any.Size())
	should.Equal(Array, any.ValueType())
	should.Nil(any.LastError())
	should.Equal(0, any.ToInt())
	should.Equal(int32(0), any.ToInt32())
	should.Equal(int64(0), any.ToInt64())
	should.Equal(uint(0), any.ToUint())
	should.Equal(uint32(0), any.ToUint32())
	should.Equal(uint64(0), any.ToUint64())
	should.Equal(float32(0), any.ToFloat32())
	should.Equal(float64(0), any.ToFloat64())
}

func Test_read_one_element_array_as_any(t *testing.T) {
	should := require.New(t)
	any := Get([]byte("[1]"))
	should.Equal(1, any.Size())
}

func Test_read_two_element_array_as_any(t *testing.T) {
	should := require.New(t)
	any := Get([]byte("[1,2]"))
	should.Equal(1, any.Get(0).ToInt())
	should.Equal(2, any.Size())
	should.True(any.ToBool())
	should.Equal(1, any.ToInt())
	should.Equal(1, any.GetArray()[0].ToInt())
	should.Equal([]interface{}{float64(1), float64(2)}, any.GetInterface())
	stream := NewStream(ConfigDefault, nil, 32)
	any.WriteTo(stream)
	should.Equal("[1,2]", string(stream.Buffer()))
}

func Test_wrap_array(t *testing.T) {
	should := require.New(t)
	any := Wrap([]int{1, 2, 3})
	should.Equal("[1,2,3]", any.ToString())
}

func Test_array_lazy_any_get(t *testing.T) {
	should := require.New(t)
	any := Get([]byte("[1,[2,3],4]"))
	should.Equal(3, any.Get(1, 1).ToInt())
	should.Equal("[1,[2,3],4]", any.ToString())
}

func Test_array_lazy_any_get_all(t *testing.T) {
	should := require.New(t)
	any := Get([]byte("[[1],[2],[3,4]]"))
	should.Equal("[1,2,3]", any.Get('*', 0).ToString())
}

func Test_array_wrapper_any_get_all(t *testing.T) {
	should := require.New(t)
	any := wrapArray([][]int{
		{1, 2},
		{3, 4},
		{5, 6},
	})
	should.Equal("[1,3,5]", any.Get('*', 0).ToString())
}

func Test_array_lazy_any_get_invalid(t *testing.T) {
	should := require.New(t)
	any := Get([]byte("[]"))
	should.Equal(Invalid, any.Get(1, 1).ValueType())
	should.NotNil(any.Get(1, 1).LastError())
	should.Equal(Invalid, any.Get("1").ValueType())
	should.NotNil(any.Get("1").LastError())
}

func Test_invalid_array(t *testing.T) {
	should := require.New(t)
	any := Get([]byte("["), 0)
	should.Equal(Invalid, any.ValueType())
}