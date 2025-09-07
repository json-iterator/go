package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

type Foo struct {
	Bar interface{}
}

func (f Foo) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(f.Bar)
	return buf.Bytes(), err
}

// Standard Encoder has trailing newline.
func TestEncodeMarshalJSON(t *testing.T) {

	foo := Foo{
		Bar: 123,
	}
	should := require.New(t)
	var buf, stdbuf bytes.Buffer
	enc := jsoniter.ConfigCompatibleWithStandardLibrary.NewEncoder(&buf)
	enc.Encode(foo)
	stdenc := json.NewEncoder(&stdbuf)
	stdenc.Encode(foo)
	should.Equal(stdbuf.Bytes(), buf.Bytes())
}

func TestMarshalObjectWithCycle(t *testing.T) {
	type A struct {
		A *A
	}
	a := A{}
	a.A = &a

	api := jsoniter.ConfigCompatibleWithStandardLibrary

	if _, err := jsoniter.Marshal(a); !errors.Is(err, jsoniter.ErrCycleEncountered) {
		t.Fatal(err)
	}

	if err := api.NewEncoder(nil).Encode(a); !errors.Is(err, jsoniter.ErrCycleEncountered) {
		t.Fatal(err)
	}
}
