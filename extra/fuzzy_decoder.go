package extra

import (
	"github.com/json-iterator/go"
	"unsafe"
	"encoding/json"
	"strings"
)

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func RegisterFuzzyDecoders() {
	jsoniter.RegisterTypeDecoder("string", &FuzzyStringDecoder{})
	jsoniter.RegisterTypeDecoder("int", &FuzzyNumberDecoder{func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator, errorReporter *jsoniter.Iterator) {
		if isFloat {
			val := iter.ReadFloat64()
			if val > float64(MaxInt) || val < float64(MinInt) {
				errorReporter.ReportError("fuzzy decode int", "exceed range")
				return
			}
			*((*int)(ptr)) = int(val)
		} else {
			*((*int)(ptr)) = iter.ReadInt()
		}
	}})
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

type FuzzyNumberDecoder struct {
	fun func(isFloat bool, ptr unsafe.Pointer, iter *jsoniter.Iterator, errorReporter *jsoniter.Iterator)
}

func (decoder *FuzzyNumberDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valueType := iter.WhatIsNext()
	var str string
	switch valueType {
	case jsoniter.Number:
		var number json.Number
		iter.ReadVal(&number)
		str = string(number)
	case jsoniter.String:
		str = iter.ReadString()
	default:
		iter.ReportError("FuzzyNumberDecoder", "not number or string")
	}
	newIter := iter.Config().BorrowIterator([]byte(str))
	defer iter.Config().ReturnIterator(newIter)
	isFloat := strings.IndexByte(str, '.') != -1
	decoder.fun(isFloat, ptr, newIter, iter)
}
