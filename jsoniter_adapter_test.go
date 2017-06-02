package jsoniter

import (
	"testing"
	"github.com/json-iterator/go/require"
	"encoding/json"
	"bytes"
	"io/ioutil"
)

func Test_new_decoder(t *testing.T) {
	should := require.New(t)
	decoder1 := json.NewDecoder(bytes.NewBufferString(`[1][2]`))
	decoder2 := NewDecoder(bytes.NewBufferString(`[1][2]`))
	arr1 := []int{}
	should.Nil(decoder1.Decode(&arr1))
	should.Equal([]int{1}, arr1)
	arr2 := []int{}
	should.True(decoder1.More())
	buffered, _ := ioutil.ReadAll(decoder1.Buffered())
	should.Equal("[2]", string(buffered))
	should.Nil(decoder2.Decode(&arr2))
	should.Equal([]int{1}, arr2)
	should.True(decoder2.More())
	buffered, _ = ioutil.ReadAll(decoder2.Buffered())
	should.Equal("[2]", string(buffered))

	should.Nil(decoder1.Decode(&arr1))
	should.Equal([]int{2}, arr1)
	should.False(decoder1.More())
	should.Nil(decoder2.Decode(&arr2))
	should.Equal([]int{2}, arr2)
	should.False(decoder2.More())
}

func Test_new_encoder(t *testing.T) {
	should := require.New(t)
	buf1 := &bytes.Buffer{}
	encoder1 := json.NewEncoder(buf1)
	encoder1.Encode([]int{1})
	should.Equal("[1]\n", buf1.String())
	buf2 := &bytes.Buffer{}
	encoder2 := NewEncoder(buf2)
	encoder2.Encode([]int{1})
	should.Equal("[1]", buf2.String())
}