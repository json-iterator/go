package jsoniter

import (
	"fmt"
	"testing"

	"github.com/json-iterator/go/require"
)

var boolConvertMap = map[string]bool{
	"null":  false,
	"true":  true,
	"false": false,

	`"true"`:  true,
	`"false"`: true,

	"123":   true,
	`"123"`: true,
	"0":     false,
	`"0"`:   false,
	"-1":    true,
	`"-1"`:  true,

	"1.1":       true,
	"0.0":       false,
	"-1.1":      true,
	`""`:        false,
	"[1,2]":     true,
	"[]":        false,
	"{}":        true,
	`{"abc":1}`: true,
}

func Test_read_bool_as_any(t *testing.T) {
	should := require.New(t)

	var any Any
	for k, v := range boolConvertMap {
		any = Get([]byte(k))
		if v {
			should.True(any.ToBool(), fmt.Sprintf("origin val is %v", k))
		} else {
			should.False(any.ToBool(), fmt.Sprintf("origin val is %v", k))
		}
	}
}
