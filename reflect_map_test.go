package jsoniter

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncoderIndentionStep(t *testing.T) {
	should := require.New(t)

	res := `{
    "first-map": {
        "0": {
            "a": "some",
            "b": "test"
        },
        "1": {
            "a": "some",
            "b": "test"
        },
        "2": {
            "a": "some",
            "b": "test"
        }
    },
    "second-map": {
        "first": {
            "a": "321",
            "b": "123"
        },
        "second": {
            "a": "321",
            "b": "123"
        },
        "third": {
            "a": "321",
            "b": "123"
        }
    },
    "custom-struct": {
        "a": "some",
        "b": "struct"
    },
    "number": 42
}
`

	json := Config{
		EscapeHTML:  true,
		SortMapKeys: true,
	}.Froze()

	type Struct struct {
		A string `json:"a"`
		B string `json:"b"`
	}

	ForTest := struct {
		M      map[int]Struct    `json:"first-map"`
		M2     map[string]Struct `json:"second-map"`
		S      Struct            `json:"custom-struct"`
		Number int               `json:"number"`
	}{
		M: map[int]Struct{
			0: Struct{"some", "test"},
			1: Struct{"some", "test"},
			2: Struct{"some", "test"},
		},
		M2: map[string]Struct{
			"first":  Struct{"321", "123"},
			"second": Struct{"321", "123"},
			"third":  Struct{"321", "123"},
		},
		S:      Struct{"some", "struct"},
		Number: 42,
	}

	buff := bytes.NewBuffer(nil)

	enc := json.NewEncoder(buff)
	enc.SetIndent("", "    ")
	enc.Encode(ForTest)

	should.Equal(res, buff.String())
	// io.Copy(os.Stdout, buff)
}
