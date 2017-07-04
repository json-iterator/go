package jsoniter

import (
	"testing"

	"github.com/json-iterator/go/require"
)

var floatConvertMap = map[string]float64{
	"null":  0,
	"true":  1,
	"false": 0,

	`"true"`:  0,
	`"false"`: 0,

	"123":       123,
	`"123true"`: 123,

	`"-123true"`: -123,
	"0":          0,
	`"0"`:        0,
	"-1":         -1,

	"1.1":       1.1,
	"0.0":       0,
	"-1.1":      -1.1,
	`"+1.1"`:    1.1,
	`""`:        0,
	"[1,2]":     1,
	"[]":        0,
	"{}":        0,
	`{"abc":1}`: 0,
}

func Test_read_any_to_float(t *testing.T) {
	should := require.New(t)
	for k, v := range floatConvertMap {
		any := Get([]byte(k))
		should.Equal(float64(v), any.ToFloat64(), "the original val is "+k)
	}

	for k, v := range floatConvertMap {
		any := Get([]byte(k))
		should.Equal(float32(v), any.ToFloat32(), "the original val is "+k)
	}
}

func Test_read_float_as_any(t *testing.T) {
	should := require.New(t)
	any := Get([]byte("12.3"))
	should.Equal(float64(12.3), any.ToFloat64())
	should.Equal("12.3", any.ToString())
	should.True(any.ToBool())
}
