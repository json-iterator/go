package extra

import (
	"encoding/json"
	"testing"

	"github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

func init() {
	jsoniter.RegisterExtension(&BinaryAsStringExtension{})
}

type TestBytesStruct struct {
	Message []byte `json:"message"`
}

type TestStringStruct struct {
	Message string `json:"message"`
}

func TestBinaryAsStringCodec(t *testing.T) {
	t.Run("test quotes & escape struct", func(t *testing.T) {
		should := require.New(t)
		msg := "\\\"hello\"\\"
		expected, err := json.Marshal(TestStringStruct{Message: msg})
		should.NoError(err)
		output, err := jsoniter.Marshal(TestBytesStruct{Message: []byte(msg)})
		should.NoError(err)
		should.Equal(string(expected), string(output))
		var val *TestBytesStruct
		should.NoError(jsoniter.Unmarshal(output, &val))
		should.Equal(msg, string(val.Message))
	})
	t.Run("safe set", func(t *testing.T) {
		should := require.New(t)
		output, err := jsoniter.Marshal([]byte("hello"))
		should.NoError(err)
		should.Equal(`"hello"`, string(output))
		var val []byte
		should.NoError(jsoniter.Unmarshal(output, &val))
		should.Equal(`hello`, string(val))
	})
	t.Run("non safe set", func(t *testing.T) {
		should := require.New(t)
		msg := `\x01\x02\x03\x0f`
		expected, err := json.Marshal(msg)
		should.NoError(err)
		output, err := jsoniter.Marshal([]byte(msg))
		should.NoError(err)
		should.Equal(string(expected), string(output))
		var val []byte
		should.NoError(jsoniter.Unmarshal(output, &val))
		should.Equal([]byte{1, 2, 3, 15}, val)
	})
}
