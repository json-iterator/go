package extra

import (
	"testing"

	"github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
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

func Test_lower_case_with_underscores_omit_empty(t *testing.T) {
	should := require.New(t)
	SetNamingStrategy(LowerCaseWithUnderscores)
	output, err := jsoniter.Marshal(FooOmitEmpty{BarValue: "Test"})
	should.Nil(err)
	should.Equal(`{"bar_value":"Test","additional_value":""}`, string(output))
}

func Test_lower_case_with_underscores_field_customziation(t *testing.T) {
	t.Skip("undo skip after PR is merged https://github.com/json-iterator/go/pull/275")
	should := require.New(t)
	SetNamingStrategy(LowerCaseWithUnderscores)
	output, err := jsoniter.Marshal(FooFieldCustomization{BarValue: "test"})
	should.Nil(err)
	should.Equal(`{"bar":"test"}`, string(output))
}

type Foo struct {
	Bar struct {
		Thing1 string
		Thing2 string
	}
}

type FooOmitEmpty struct {
	BarValue        string `json:"omitempty"`
	AdditionalValue string
}
type FooFieldCustomization struct {
	BarValue string `json:"bar,omitempty"`
}
