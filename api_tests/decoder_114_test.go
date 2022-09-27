//+build go1.14

package test

import (
	"encoding/json"
	"strings"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

type decoderIfc interface {
	Decode(obj interface{}) error
	More() bool
	InputOffset() int64
}

func Test_decoder_input_offset(t *testing.T) {
	should := require.New(t)
	input := `{"foo" : "bar"}` + "\n" + `{"qoo" : "baz"}`
	newlinePos := strings.IndexByte(input, '\n')

	runChecks := func(decoder decoderIfc) {
		should.True(decoder.More())
		obj := map[string]interface{}{}
		should.NoError(decoder.Decode(&obj))
		should.Len(obj, 1)
		should.True(decoder.More())
		should.EqualValues(newlinePos+1, decoder.InputOffset())
		should.NoError(decoder.Decode(&obj))
		should.EqualValues(len(input), decoder.InputOffset())
		should.False(decoder.More())
	}

	// try with stdlib json
	runChecks(json.NewDecoder(strings.NewReader(input)))
	// and with jsoniter
	runChecks(jsoniter.NewDecoder(strings.NewReader(input)))
}
