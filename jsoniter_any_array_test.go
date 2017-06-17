package jsoniter

import (
	"testing"
	"github.com/json-iterator/go/require"
	"io"
)

func Test_read_empty_array_as_any(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString("[]")
	should.Nil(err)
	should.Equal(0, any.Size())
}

func Test_read_one_element_array_as_any(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString("[1]")
	should.Nil(err)
	should.Equal(1, any.Size())
}

func Test_read_two_element_array_as_any(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString("[1,2]")
	should.Nil(err)
	should.Equal(1, any.Get(0).ToInt())
	should.Equal(2, any.Size())
	should.True(any.ToBool())
	should.Equal(1, any.ToInt())
}

func Test_read_array_with_any_iterator(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString("[1,2]")
	should.Nil(err)
	var element Any
	var elements []int
	for next, hasNext := any.IterateArray(); hasNext; {
		element, hasNext = next()
		elements = append(elements, element.ToInt())
	}
	should.Equal([]int{1, 2}, elements)
}

func Test_wrap_array(t *testing.T) {
	should := require.New(t)
	any := Wrap([]int{1, 2, 3})
	should.Equal("[1,2,3]", any.ToString())
	var element Any
	var elements []int
	for next, hasNext := any.IterateArray(); hasNext; {
		element, hasNext = next()
		elements = append(elements, element.ToInt())
	}
	should.Equal([]int{1, 2, 3}, elements)
	any = Wrap([]int{1, 2, 3})
	should.Equal(3, any.Size())
	any = Wrap([]int{1, 2, 3})
	should.Equal(2, any.Get(1).ToInt())
}

func Test_array_lazy_any_get(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString("[1,[2,3],4]")
	should.Nil(err)
	should.Equal(3, any.Get(1, 1).ToInt())
	should.Equal("[1,[2,3],4]", any.ToString())
}

func Test_array_lazy_any_get_all(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString("[[1],[2],[3,4]]")
	should.Nil(err)
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
	any, err := UnmarshalAnyFromString("[]")
	should.Nil(err)
	should.Equal(Invalid, any.Get(1, 1).ValueType())
	should.NotNil(any.Get(1, 1).LastError())
	should.Equal(Invalid, any.Get("1").ValueType())
	should.NotNil(any.Get("1").LastError())
}

func Test_array_lazy_any_set(t *testing.T) {
	should := require.New(t)
	any, err := UnmarshalAnyFromString("[1,[2,3],4]")
	should.Nil(err)
	any.GetArray()[0] = WrapInt64(2)
	str, err := MarshalToString(any)
	should.Nil(err)
	should.Equal("[2,[2,3],4]", str)
}

func Test_invalid_array(t *testing.T) {
	_, err := UnmarshalAnyFromString("[")
	if err == nil || err == io.EOF {
		t.FailNow()
	}
}