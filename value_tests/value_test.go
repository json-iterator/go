package test

import (
	"testing"
	"reflect"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"github.com/json-iterator/go"
)

type unmarshalCase struct {
	ptr interface{}
	input string
}

var unmarshalCases []unmarshalCase

var marshalCases = []interface{}{
	nil,
}

func Test_unmarshal(t *testing.T) {
	should := require.New(t)
	for _, testCase := range unmarshalCases {
		valType := reflect.TypeOf(testCase.ptr).Elem()
		ptr1Val := reflect.New(valType)
		err1 := json.Unmarshal([]byte(testCase.input), ptr1Val.Interface())
		should.NoError(err1)
		ptr2Val := reflect.New(valType)
		err2 := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(testCase.input), ptr2Val.Interface())
		should.NoError(err2)
		should.Equal(ptr1Val.Interface(), ptr2Val.Interface())
	}
}

func Test_marshal(t *testing.T) {
	should := require.New(t)
	for _, testCase := range marshalCases {
		output1, err1 := json.Marshal(testCase)
		should.NoError(err1)
		output2, err2 := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(testCase)
		should.NoError(err2)
		should.Equal(string(output1), string(output2))
	}
}