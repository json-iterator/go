package extra

import (
	"github.com/json-iterator/go"
	"unsafe"
	"encoding/json"
)

func RegisterFuzzyDecoders() {
	jsoniter.RegisterTypeDecoder("string", &FuzzyStringDecoder{})
	jsoniter.RegisterTypeDecoder("int", &FuzzyIntDecoder{})
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

type FuzzyIntDecoder struct {
}

func (decoder *FuzzyIntDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	switch valueType {
	case jsoniter.Number:
		// use current iterator
	case jsoniter.String:
		str := iter.ReadString()
		iter = iter.Config().BorrowIterator([]byte(str))
		defer iter.Config().ReturnIterator(iter)
	default:
		iter.ReportError("FuzzyIntDecoder", "not number or string")
	}
	*((*int)(ptr)) = int(iter.ReadFloat64())
}
