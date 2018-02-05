package test

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
	"github.com/json-iterator/go"
)

func Test_disallowUnknownFields(t *testing.T) {
	should := require.New(t)
	type TestObject struct{}
	var obj TestObject
	decoder := jsoniter.NewDecoder(bytes.NewBufferString(`{"field1":100}`))
	decoder.DisallowUnknownFields()
	should.Error(decoder.Decode(&obj))
}
