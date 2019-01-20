package jsoniter

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapIndention(t *testing.T) {
	should := require.New(t)

	type object struct {
		A string `json:"a"`
		B string `json:"b"`
	}

	testStruct := struct {
		M      map[int]object    `json:"first-map"`
		M2     map[string]object `json:"second-map"`
		S      object            `json:"custom-struct"`
		Number int               `json:"number"`
	}{
		M: map[int]object{
			0: object{"some", "test"},
			1: object{"some", "test"},
			2: object{"some", "test"},
		},
		M2: map[string]object{
			"first":  object{"321", "123"},
			"second": object{"321", "123"},
			"third":  object{"321", "123"},
		},
		S:      object{"some", "struct"},
		Number: 42,
	}

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
}`

	json := Config{
		EscapeHTML:  true,
		SortMapKeys: true,
	}.Froze()

	// Test MarshalIndent
	data, err := json.MarshalIndent(testStruct, "", "    ")
	should.Nil(err)
	should.Equal(res, string(data))

	// Test Encoder. Have to add '\n' into res because Encode() adds additional '\n'
	res += "\n"
	buff := bytes.NewBuffer(nil)

	enc := json.NewEncoder(buff)
	enc.SetIndent("", "    ")
	enc.Encode(testStruct)

	should.Equal(res, buff.String())

	// io.Copy(os.Stdout, buff) // for debug
}
