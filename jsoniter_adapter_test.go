package jsoniter

import (
	"testing"
	"github.com/json-iterator/go/require"
	"encoding/json"
	"bytes"
)

func Test_new_decoder(t *testing.T) {
	should := require.New(t)
	decoder1 := json.NewDecoder(bytes.NewBufferString(`[1]`))
	decoder2 := NewDecoder(bytes.NewBufferString(`[1]`))
	arr1 := []int{}
	should.Nil(decoder1.Decode(&arr1))
	should.Equal([]int{1}, arr1)
	arr2 := []int{}
	should.Nil(decoder2.Decode(&arr2))
}