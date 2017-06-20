package extra

import (
	"github.com/json-iterator/go"
	"unsafe"
	"encoding/json"
)

func RegisterFuzzyDecoders() {
	jsoniter.RegisterTypeDecoder("string", &FuzzyStringDecoder{})
}

type FuzzyStringDecoder struct {
}

func (decoder *FuzzyStringDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	switch valueType {
	case jsoniter.Number:
		var number json.Number
		iter.ReadVal(&number)
		*((*string)(ptr)) = string(number)
	case jsoniter.String:
		*((*string)(ptr)) = iter.ReadString()
	default:
		iter.ReportError("FuzzyStringDecoder", "not number or string")
	}
}
