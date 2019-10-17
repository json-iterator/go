package test

import (
	"bytes"
	"encoding/json"
	"testing"
	"unicode/utf8"

	"github.com/json-iterator/go"
)

func init() {
	marshalCases = append(marshalCases,
		`>`,
		`"数字山谷"`,
		"he\u2029\u2028he",
	)
	for i := 0; i < utf8.RuneSelf; i++ {
		marshalCases = append(marshalCases, string([]byte{byte(i)}))
	}
}

func Test_read_string(t *testing.T) {
	badInputs := []string{
		``,
		`"`,
		`"\"`,
		`"\\\"`,
		"\"\n\"",
		`"\U0001f64f"`,
		`"\uD83D\u00"`,
	}
	for i := 0; i < 32; i++ {
		// control characters are invalid
		badInputs = append(badInputs, string([]byte{'"', byte(i), '"'}))
	}

	for _, input := range badInputs {
		testReadString(t, input, "", nil, true, "json.Unmarshal", json.Unmarshal, nil)
		testReadString(t, input, "", nil, true, "jsoniter.Unmarshal", jsoniter.Unmarshal, nil)
		testReadString(t, input, "", nil, true, "jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal", jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal, nil)
	}

	goodInputs := []struct {
		input           string
		expectValue     string
		expectRemarshal []byte
	}{
		{input: `""`, expectValue: ""},
		{input: `"a"`, expectValue: "a"},
		{input: `null`, expectValue: "", expectRemarshal: []byte(`""`)},
		{input: `"Iñtërnâtiônàlizætiøn,💝🐹🌇⛔"`, expectValue: "Iñtërnâtiônàlizætiøn,💝🐹🌇⛔"},
		{input: `"\uD83D"`, expectValue: string([]byte{239, 191, 189}), expectRemarshal: []byte{'"', 239, 191, 189, '"'}},
		{input: `"\uD83D\\"`, expectValue: string([]byte{239, 191, 189, '\\'}), expectRemarshal: []byte{'"', 239, 191, 189, '\\', '\\', '"'}},
		{input: `"\uD83D\ub000"`, expectValue: string([]byte{239, 191, 189, 235, 128, 128}), expectRemarshal: []byte{'"', 239, 191, 189, 235, 128, 128, '"'}},
		{input: `"\uD83D\ude04"`, expectValue: "😄", expectRemarshal: []byte(`"😄"`)},
		{input: `"\uDEADBEEF"`, expectValue: string([]byte{239, 191, 189, 66, 69, 69, 70}), expectRemarshal: []byte{'"', 239, 191, 189, 66, 69, 69, 70, '"'}},
		{input: `"hel\"lo"`, expectValue: `hel"lo`, expectRemarshal: []byte(`"hel\"lo"`)},
		{input: `"hel\\\/lo"`, expectValue: `hel\/lo`, expectRemarshal: []byte(`"hel\\/lo"`)},
		{input: `"hel\\blo"`, expectValue: `hel\blo`},
		{input: `"hel\\\blo"`, expectValue: "hel\\\blo", expectRemarshal: []byte(`"hel\\\u0008lo"`)},
		{input: `"hel\\nlo"`, expectValue: `hel\nlo`},
		{input: `"hel\\\nlo"`, expectValue: "hel\\\nlo"},
		{input: `"hel\\tlo"`, expectValue: `hel\tlo`},
		{input: `"hel\\flo"`, expectValue: `hel\flo`},
		{input: `"hel\\\flo"`, expectValue: "hel\\\flo", expectRemarshal: []byte(`"hel\\\u000clo"`)},
		{input: `"hel\\\rlo"`, expectValue: "hel\\\rlo", expectRemarshal: []byte(``)},
		{input: `"hel\\\tlo"`, expectValue: "hel\\\tlo", expectRemarshal: []byte(``)},
		{input: `"\u4e2d\u6587"`, expectValue: "中文", expectRemarshal: []byte(`"中文"`)},
		{input: `"\ud83d\udc4a"`, expectValue: "\xf0\x9f\x91\x8a", expectRemarshal: []byte("\"\xf0\x9f\x91\x8a\"")},
		// single-byte invalid utf8 encoding:
		{input: `"` + string([]byte{147}) + `"`, expectValue: "\ufffd", expectRemarshal: []byte("\"\ufffd\"")},
		// single-byte invalid utf8 encoding followed by valid extended character
		{input: `"` + string([]byte{147}) + `⛔"`, expectValue: "\ufffd⛔", expectRemarshal: []byte("\"\ufffd⛔\"")},
		// multi-byte invalid utf8 encoding
		{input: `"` + string([]byte{226, 128}) + `"`, expectValue: "\ufffd\ufffd", expectRemarshal: []byte("\"\ufffd\ufffd\"")},
		// multi-byte invalid utf8 encoding followed by valid extended character
		{input: `"` + string([]byte{226, 128}) + `⛔"`, expectValue: "\ufffd\ufffd⛔", expectRemarshal: []byte("\"\ufffd\ufffd⛔\"")},
		// valid multi-byte followed by invalid single-byte
		{input: `"` + string([]byte{226, 128, 168, 138}) + `"`, expectValue: string([]byte{226, 128, 168}) + "\ufffd", expectRemarshal: []byte("\"\\u2028\ufffd\"")},
	}

	for _, tc := range goodInputs {
		testReadString(t, tc.input, tc.expectValue, tc.expectRemarshal, false, "json.Unmarshal", json.Unmarshal, json.Marshal)
		testReadString(t, tc.input, tc.expectValue, tc.expectRemarshal, false, "jsoniter.Unmarshal", jsoniter.Unmarshal, jsoniter.Marshal)
		testReadString(t, tc.input, tc.expectValue, tc.expectRemarshal, false, "jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal", jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal, jsoniter.ConfigCompatibleWithStandardLibrary.Marshal)
	}
}

func testReadString(t *testing.T, input string, expectValue string, expectRemarshal []byte, expectError bool, marshalerName string, unmarshaler func([]byte, interface{}) error, marshaler func(interface{}) ([]byte, error)) {
	var value string
	err := unmarshaler([]byte(input), &value)
	if expectError != (err != nil) {
		t.Errorf("%q: %s: expected error %v, got %v", input, marshalerName, expectError, err)
		return
	}
	if value != expectValue {
		t.Errorf("%q: %s: expected %q (%v), got %q (%v)", input, marshalerName, expectValue, []byte(expectValue), value, []byte(value))
		return
	}

	if expectError {
		return
	}

	// Test re-marshal
	remarshal, err := marshaler(expectValue)
	if err != nil {
		t.Errorf("%q: %s: unexpected error remarshaling: %v", input, marshalerName, err)
		return
	}
	if len(expectRemarshal) == 0 {
		expectRemarshal = []byte(input)
	}
	if bytes.Compare(remarshal, expectRemarshal) != 0 {
		t.Errorf("%q: %s: expected %q, got %q remarshaling", input, marshalerName, string(expectRemarshal), string(remarshal))
		return
	}

	// Test round-trip is a no-op
	var value2 string
	err = unmarshaler(remarshal, &value2)
	if expectError != (err != nil) {
		t.Errorf("%q: %s: expected error %v, got %v", input, marshalerName, expectError, err)
		return
	}
	if value2 != expectValue {
		t.Errorf("%q: %s: expected %q (%v), got %q (%v)", input, marshalerName, expectValue, []byte(expectValue), value2, []byte(value2))
		return
	}
}
