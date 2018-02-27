package extra

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/json-iterator/go"
)

func init() {
	jsoniter.RegisterExtension(&BinaryAsStringExtension{})
}

func TestBinaryAsStringCodec(t *testing.T) {
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
		output, err := jsoniter.Marshal([]byte{1, 2, 3, 15})
		should.NoError(err)
		should.Equal(`"\\x01\\x02\\x03\\x0f"`, string(output))
		var val []byte
		should.NoError(jsoniter.Unmarshal(output, &val))
		should.Equal([]byte{1, 2, 3, 15}, val)
	})
}
