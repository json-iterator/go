package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/json-iterator/go"
	"encoding/json"
)

func Test_use_number_for_unmarshal(t *testing.T) {
	should := require.New(t)
	api := jsoniter.Config{UseNumber: true}.Froze()
	var obj interface{}
	should.Nil(api.UnmarshalFromString("123", &obj))
	should.Equal(json.Number("123"), obj)
}
