package test

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
	"testing"
)

type Obj struct {
	Bar []int64 `json:",,string"`
}

func TestMarshal(t *testing.T) {
	should := require.New(t)
	obj := Obj{
		Bar: []int64{
			1, 2, 3, 4,
		},
	}

	str, err := jsoniter.Marshal(&obj)
	should.Nil(err)
	should.NotEmpty(str)
	fmt.Println(string(str))
}

func TestUnmarshal(t *testing.T) {
	should := require.New(t)
	str := `{"Bar":["1","2","3","4", "122313213"]}`
	obj := Obj{}
	err := jsoniter.Unmarshal([]byte(str), &obj)
	should.NotNil(obj)
	should.Nil(err)
	fmt.Println(obj)
}

