package test

import (
	"encoding/json"
	"io"
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/json-iterator/go"
	"errors"
)

func Test_skip(t *testing.T) {
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			should := require.New(t)
			dst := typeForTest(0)
			stdErr := json.Unmarshal([]byte(input), &dst)
			iter := jsoniter.ParseString(jsoniter.ConfigDefault, input)
			iter.Skip()
			iter.ReadNil() // trigger looking forward
			err := iter.Error
			if err == io.EOF {
				err = nil
			} else {
				err = errors.New("remaining bytes")
			}
			if stdErr == nil {
				should.Nil(err)
			} else {
				should.NotNil(err)
			}
		})
	}
}