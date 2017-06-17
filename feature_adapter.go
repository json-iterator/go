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
	"encoding/json"
	"io"
	"unsafe"
)

// Unmarshal adapts to json/encoding Unmarshal API
//
// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
// Refer to https://godoc.org/encoding/json#Unmarshal for more information
func Unmarshal(data []byte, v interface{}) error {
	return ConfigDefault.Unmarshal(data, v)
}

// UnmarshalAny adapts to
func UnmarshalAny(data []byte) (Any, error) {
	return ConfigDefault.UnmarshalAny(data)
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

func UnmarshalAnyFromString(str string) (Any, error) {
	return ConfigDefault.UnmarshalAnyFromString(str)
}

// Marshal adapts to json/encoding Marshal API
//
// Marshal returns the JSON encoding of v, adapts to json/encoding Marshal API
// Refer to https://godoc.org/encoding/json#Marshal for more information
func Marshal(v interface{}) ([]byte, error) {
	return ConfigDefault.Marshal(v)
}

func MarshalToString(v interface{}) (string, error) {
	return ConfigDefault.MarshalToString(v)
}

// NewDecoder adapts to json/stream NewDecoder API.
//
// NewDecoder returns a new decoder that reads from r.
//
// Instead of a json/encoding Decoder, an AdaptedDecoder is returned
// Refer to https://godoc.org/encoding/json#NewDecoder for more information
func NewDecoder(reader io.Reader) *AdaptedDecoder {
	return ConfigDefault.NewDecoder(reader)
}

// AdaptedDecoder reads and decodes JSON values from an input stream.
// AdaptedDecoder provides identical APIs with json/stream Decoder (Token() and UseNumber() are in progress)
type AdaptedDecoder struct {
	iter *Iterator
}

func (adapter *AdaptedDecoder) Decode(obj interface{}) error {
	adapter.iter.ReadVal(obj)
	err := adapter.iter.Error
	if err == io.EOF {
		return nil
	}
	return adapter.iter.Error
}

func (adapter *AdaptedDecoder) More() bool {
	return adapter.iter.head != adapter.iter.tail
}

func (adapter *AdaptedDecoder) Buffered() io.Reader {
	remaining := adapter.iter.buf[adapter.iter.head:adapter.iter.tail]
	return bytes.NewReader(remaining)
}

func (decoder *AdaptedDecoder) UseNumber() {
	RegisterTypeDecoder("interface {}", func(ptr unsafe.Pointer, iter *Iterator) {
		if iter.WhatIsNext() == Number {
			*((*interface{})(ptr)) = json.Number(iter.readNumberAsString())
		} else {
			*((*interface{})(ptr)) = iter.Read()
		}
	})
}

func NewEncoder(writer io.Writer) *AdaptedEncoder {
	return ConfigDefault.NewEncoder(writer)
}

type AdaptedEncoder struct {
	stream *Stream
}

func (adapter *AdaptedEncoder) Encode(val interface{}) error {
	adapter.stream.WriteVal(val)
	adapter.stream.Flush()
	return adapter.stream.Error
}

func (adapter *AdaptedEncoder) SetIndent(prefix, indent string) {
	adapter.stream.cfg.indentionStep = len(indent)
}

func (adapter *AdaptedEncoder) SetEscapeHTML(escapeHtml bool) {
	config := adapter.stream.cfg.configBeforeFrozen
	config.EscapeHtml = escapeHtml
	adapter.stream.cfg = config.Froze()
}
