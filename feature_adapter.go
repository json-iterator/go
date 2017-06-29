// Package jsoniter implements encoding and decoding of JSON as defined in
// RFC 4627 and provides interfaces with identical syntax of standard lib encoding/json.
// Converting from encoding/json to jsoniter is no more than replacing the package with jsoniter
// and variable type declarations (if any).
// jsoniter interfaces gives 100% compatibility with code using standard lib.
//
// "JSON and Go"
// (https://golang.org/doc/articles/json_and_go.html)
// gives a description of how Marshal/Unmarshal operate
// between arbitrary or predefined json objects and bytes,
// and it applies to jsoniter.Marshal/Unmarshal as well.
package jsoniter

import (
	"bytes"
	"io"
)

type RawMessage []byte

// Unmarshal adapts to json/encoding Unmarshal API
//
// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
// Refer to https://godoc.org/encoding/json#Unmarshal for more information
func Unmarshal(data []byte, v interface{}) error {
	return ConfigDefault.Unmarshal(data, v)
}

func lastNotSpacePos(data []byte) int {
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] != ' ' && data[i] != '\t' && data[i] != '\r' && data[i] != '\n' {
			return i + 1
		}
	}
	return 0
}

func UnmarshalFromString(str string, v interface{}) error {
	return ConfigDefault.UnmarshalFromString(str, v)
}

func Get(data []byte, path ...interface{}) Any {
	return ConfigDefault.Get(data, path...)
}

// Marshal adapts to json/encoding Marshal API
//
// Marshal returns the JSON encoding of v, adapts to json/encoding Marshal API
// Refer to https://godoc.org/encoding/json#Marshal for more information
func Marshal(v interface{}) ([]byte, error) {
	return ConfigDefault.Marshal(v)
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return ConfigDefault.MarshalIndent(v, prefix, indent)
}

func MarshalToString(v interface{}) (string, error) {
	return ConfigDefault.MarshalToString(v)
}

// NewDecoder adapts to json/stream NewDecoder API.
//
// NewDecoder returns a new decoder that reads from r.
//
// Instead of a json/encoding Decoder, an Decoder is returned
// Refer to https://godoc.org/encoding/json#NewDecoder for more information
func NewDecoder(reader io.Reader) *Decoder {
	return ConfigDefault.NewDecoder(reader)
}

// Decoder reads and decodes JSON values from an input stream.
// Decoder provides identical APIs with json/stream Decoder (Token() and UseNumber() are in progress)
type Decoder struct {
	iter *Iterator
}

func (adapter *Decoder) Decode(obj interface{}) error {
	adapter.iter.ReadVal(obj)
	err := adapter.iter.Error
	if err == io.EOF {
		return nil
	}
	return adapter.iter.Error
}

func (adapter *Decoder) More() bool {
	return adapter.iter.head != adapter.iter.tail
}

func (adapter *Decoder) Buffered() io.Reader {
	remaining := adapter.iter.buf[adapter.iter.head:adapter.iter.tail]
	return bytes.NewReader(remaining)
}

func (decoder *Decoder) UseNumber() {
	origCfg := decoder.iter.cfg.configBeforeFrozen
	origCfg.UseNumber = true
	decoder.iter.cfg = origCfg.Froze()
}

func NewEncoder(writer io.Writer) *Encoder {
	return ConfigDefault.NewEncoder(writer)
}

type Encoder struct {
	stream *Stream
}

func (adapter *Encoder) Encode(val interface{}) error {
	adapter.stream.WriteVal(val)
	adapter.stream.Flush()
	return adapter.stream.Error
}

func (adapter *Encoder) SetIndent(prefix, indent string) {
	adapter.stream.cfg.indentionStep = len(indent)
}

func (adapter *Encoder) SetEscapeHTML(escapeHtml bool) {
	config := adapter.stream.cfg.configBeforeFrozen
	config.EscapeHtml = escapeHtml
	adapter.stream.cfg = config.Froze()
}
