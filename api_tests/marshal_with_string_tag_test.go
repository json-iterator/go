package test

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

type Obj struct {
	Bar []int64  `json:",,string"`
	Foo []*int `json:",string"`
}

func TestMarshal(t *testing.T) {
	should := require.New(t)
	a, b, c := 1, 2, 3
	foo := []*int{&a, &b, &c, nil}
	obj := Obj{
		Bar: []int64{
			1, 2, 3, 4,
		},
		Foo: foo,
	}

	out, err := jsoniter.Marshal(&obj)
	should.Nil(err)
	should.NotEmpty(out)
	should.JSONEq(`{"Bar":["1","2","3","4"], "Foo": ["1", "2", "3", "null"]}`, string(out))
}

func TestUnmarshal(t *testing.T) {
	should := require.New(t)
	str := `{"Bar":["1","2","3","4", "122313213"], "Foo": ["1", "2", "3", "null"]}`
	obj := Obj{}
	err := jsoniter.Unmarshal([]byte(str), &obj)
	should.NotNil(obj)
	should.Nil(err)
	should.NotNil(obj.Bar)
	should.NotNil(obj.Foo)
	should.Equal(obj.Bar[0], int64(1))
	should.Equal(*obj.Foo[0], 1)
	should.Equal(*obj.Foo[1], 2)
	should.Equal(*obj.Foo[2], 3)
	should.Nil(obj.Foo[3])
}
