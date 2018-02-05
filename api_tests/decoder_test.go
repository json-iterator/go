package test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_disallowUnknownFields(t *testing.T) {
	should := require.New(t)
	type TestObject struct{}
	var obj TestObject
	decoder := json.NewDecoder(bytes.NewBufferString(`{"field1":100}`))
	decoder.DisallowUnknownFields()
	should.Error(decoder.Decode(&obj))
}
