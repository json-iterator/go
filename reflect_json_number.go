package jsoniter

import (
	"encoding/json"
	"unsafe"

	"github.com/modern-go/reflect2"
)

type Number = json.Number

func CastJsonNumber(val interface{}) (string, bool) {
	switch typedVal := val.(type) {
	case json.Number:
		return string(typedVal), true
	}
	return "", false
}

var jsonNumberType = reflect2.TypeOfPtr((*json.Number)(nil)).Elem()

func createDecoderOfJsonNumber(ctx *ctx, typ reflect2.Type) ValDecoder {
	if typ.AssignableTo(jsonNumberType) {
		return &jsonNumberCodec{}
	}
	return nil
}

func createEncoderOfJsonNumber(ctx *ctx, typ reflect2.Type) ValEncoder {
	if typ.AssignableTo(jsonNumberType) {
		return &jsonNumberCodec{}
	}
	return nil
}

type jsonNumberCodec struct {
}

func (codec *jsonNumberCodec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	switch iter.WhatIsNext() {
	case StringValue:
		*((*json.Number)(ptr)) = json.Number(iter.ReadString())
	case NilValue:
		iter.skipFourBytes('n', 'u', 'l', 'l')
		*((*json.Number)(ptr)) = ""
	default:
		*((*json.Number)(ptr)) = json.Number([]byte(iter.readNumberAsString()))
	}
}

func (codec *jsonNumberCodec) Encode(ptr unsafe.Pointer, stream *Stream) {
	number := *((*json.Number)(ptr))
	if len(number) == 0 {
		stream.writeByte('0')
	} else {
		stream.WriteRaw(string(number))
	}
}

func (codec *jsonNumberCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return len(*((*json.Number)(ptr))) == 0
}
