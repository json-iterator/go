package any_tests

import (
	"github.com/stretchr/testify/require"
	"testing"
	"github.com/json-iterator/go"
)

func Test_wrap_map(t *testing.T) {
	should := require.New(t)
	any := jsoniter.Wrap(map[string]string{"Field1": "hello"})
	should.Equal("hello", any.Get("Field1").ToString())
	any = jsoniter.Wrap(map[string]string{"Field1": "hello"})
	should.Equal(1, any.Size())
}
