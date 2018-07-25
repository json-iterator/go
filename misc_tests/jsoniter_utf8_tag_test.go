package misc_tests

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

func TestUTF8StructTag(t *testing.T) {
	jsonBlob := []byte(`{"姓名":"燕子"}`)

	type User struct {
		Name string `json:"姓名"`
	}

	should := require.New(t)
	user := User{}
	err := jsoniter.Unmarshal(jsonBlob, &user)
	should.NoError(err)
	should.Equal("燕子", user.Name)
}
