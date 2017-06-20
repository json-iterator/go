package extra

import (
	"testing"
	"github.com/json-iterator/go"
	"github.com/json-iterator/go/require"
)

func Test_lower_case_with_underscores(t *testing.T) {
	should := require.New(t)
	should.Equal("hello_world", LowerCaseWithUnderscores("helloWorld"))
	should.Equal("hello_world", LowerCaseWithUnderscores("HelloWorld"))
	SetNamingStrategy(LowerCaseWithUnderscores)
	output, err := jsoniter.MarshalToString(struct {
		HelloWorld string
	}{
		HelloWorld: "hi",
	})
	should.Nil(err)
	should.Equal(`{"hello_world":"hi"}`, output)
}


