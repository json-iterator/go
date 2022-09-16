package extra

import (
	"github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_lower_case_with_underscores(t *testing.T) {
	should := require.New(t)
	should.Equal("hello_world", LowerCaseWithUnderscores("helloWorld"))
	should.Equal("hello_world", LowerCaseWithUnderscores("HelloWorld"))
	SetNamingStrategy(LowerCaseWithUnderscores)
	output, err := jsoniter.Marshal(struct {
		UserName      string
		FirstLanguage string
	}{
		UserName:      "taowen",
		FirstLanguage: "Chinese",
	})
	should.Nil(err)
	should.Equal(`{"user_name":"taowen","first_language":"Chinese"}`, string(output))
}

func Test_set_naming_strategy_with_overrides(t *testing.T) {
	should := require.New(t)
	SetNamingStrategy(LowerCaseWithUnderscores)
	output, err := jsoniter.Marshal(struct {
		UserName      string `json:"UserName"`
		FirstLanguage string
	}{
		UserName:      "taowen",
		FirstLanguage: "Chinese",
	})
	should.Nil(err)
	should.Equal(`{"UserName":"taowen","first_language":"Chinese"}`, string(output))
}

func Test_set_naming_strategy_with_omitempty(t *testing.T) {
	should := require.New(t)
	SetNamingStrategy(LowerCaseWithUnderscores)
	output, err := jsoniter.Marshal(struct {
		UserName      string
		FirstLanguage string `json:",omitempty"`
	}{
		UserName: "taowen",
	})
	should.Nil(err)
	should.Equal(`{"user_name":"taowen"}`, string(output))
}

func Test_set_naming_strategy_with_private_field(t *testing.T) {
	should := require.New(t)
	SetNamingStrategy(LowerCaseWithUnderscores)
	output, err := jsoniter.Marshal(struct {
		UserName string
		userId   int
		_UserAge int
	}{
		UserName: "allen",
		userId:   100,
		_UserAge: 30,
	})
	should.Nil(err)
	should.Equal(`{"user_name":"allen"}`, string(output))
}

func Test_first_case_to_lower(t *testing.T) {
	should := require.New(t)
	should.Equal("helloWorld", FirstCaseToLower("HelloWorld"))
	should.Equal("hello_World", FirstCaseToLower("Hello_World"))
}

func Test_first_case_to_lower_with_first_case_already_lowercase(t *testing.T) {
	should := require.New(t)
	should.Equal("helloWorld", FirstCaseToLower("helloWorld"))
}

func Test_first_case_to_lower_with_first_case_be_anything(t *testing.T) {
	should := require.New(t)
	should.Equal("_HelloWorld", FirstCaseToLower("_HelloWorld"))
	should.Equal("*HelloWorld", FirstCaseToLower("*HelloWorld"))
	should.Equal("?HelloWorld", FirstCaseToLower("?HelloWorld"))
	should.Equal(".HelloWorld", FirstCaseToLower(".HelloWorld"))
}
