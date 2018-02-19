package jsoniter

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"unsafe"
	"github.com/v2pro/plz/reflect2"
)

type stringCodec struct {
}

func (codec *stringCodec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*string)(ptr)) = iter.ReadString()
}

func (codec *stringCodec) Encode(ptr unsafe.Pointer, stream *Stream) {
	str := *((*string)(ptr))
	stream.WriteString(str)
}

func (codec *stringCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*string)(ptr)) == ""
}

type intCodec struct {
}

func (codec *intCodec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.ReadNil() {
		*((*int)(ptr)) = iter.ReadInt()
	}
}

func (codec *intCodec) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteInt(*((*int)(ptr)))
}

func (codec *intCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*int)(ptr)) == 0
}

type uintptrCodec struct {
}

func (codec *uintptrCodec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.ReadNil() {
		*((*uintptr)(ptr)) = uintptr(iter.ReadUint64())
	}
}

func (codec *uintptrCodec) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteUint64(uint64(*((*uintptr)(ptr))))
}

func (codec *uintptrCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*uintptr)(ptr)) == 0
}

type int8Codec struct {
}

func (codec *int8Codec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.ReadNil() {
		*((*int8)(ptr)) = iter.ReadInt8()
	}
}

func (codec *int8Codec) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteInt8(*((*int8)(ptr)))
}

func (codec *int8Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*int8)(ptr)) == 0
}

type int16Codec struct {
}

func (codec *int16Codec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.ReadNil() {
		*((*int16)(ptr)) = iter.ReadInt16()
	}
}

func (codec *int16Codec) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteInt16(*((*int16)(ptr)))
}

func (codec *int16Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*int16)(ptr)) == 0
}

type int32Codec struct {
}

func (codec *int32Codec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.ReadNil() {
		*((*int32)(ptr)) = iter.ReadInt32()
	}
}

func (codec *int32Codec) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteInt32(*((*int32)(ptr)))
}

func (codec *int32Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*int32)(ptr)) == 0
}

type int64Codec struct {
}

func (codec *int64Codec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.ReadNil() {
		*((*int64)(ptr)) = iter.ReadInt64()
	}
}

func (codec *int64Codec) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteInt64(*((*int64)(ptr)))
}

func (codec *int64Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*int64)(ptr)) == 0
}

type uintCodec struct {
}

func (codec *uintCodec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.ReadNil() {
		*((*uint)(ptr)) = iter.ReadUint()
		return
	}
}

func (codec *uintCodec) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteUint(*((*uint)(ptr)))
}

func (codec *uintCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*uint)(ptr)) == 0
}

type uint8Codec struct {
}

func (codec *uint8Codec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.ReadNil() {
		*((*uint8)(ptr)) = iter.ReadUint8()
	}
}

func (codec *uint8Codec) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteUint8(*((*uint8)(ptr)))
}

func (codec *uint8Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*uint8)(ptr)) == 0
}

type uint16Codec struct {
}

func (codec *uint16Codec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.ReadNil() {
		*((*uint16)(ptr)) = iter.ReadUint16()
	}
}

func (codec *uint16Codec) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteUint16(*((*uint16)(ptr)))
}

func (codec *uint16Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*uint16)(ptr)) == 0
}

type uint32Codec struct {
}

func (codec *uint32Codec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.ReadNil() {
		*((*uint32)(ptr)) = iter.ReadUint32()
	}
}

func (codec *uint32Codec) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteUint32(*((*uint32)(ptr)))
}

func (codec *uint32Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*uint32)(ptr)) == 0
}

type uint64Codec struct {
}

func (codec *uint64Codec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.ReadNil() {
		*((*uint64)(ptr)) = iter.ReadUint64()
	}
}

func (codec *uint64Codec) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteUint64(*((*uint64)(ptr)))
}

func (codec *uint64Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*uint64)(ptr)) == 0
}

type float32Codec struct {
}

func (codec *float32Codec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.ReadNil() {
		*((*float32)(ptr)) = iter.ReadFloat32()
	}
}

func (codec *float32Codec) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteFloat32(*((*float32)(ptr)))
}

func (codec *float32Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*float32)(ptr)) == 0
}

type float64Codec struct {
}

func (codec *float64Codec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.ReadNil() {
		*((*float64)(ptr)) = iter.ReadFloat64()
	}
}

func (codec *float64Codec) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteFloat64(*((*float64)(ptr)))
}

func (codec *float64Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*float64)(ptr)) == 0
}

type boolCodec struct {
}

func (codec *boolCodec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.ReadNil() {
		*((*bool)(ptr)) = iter.ReadBool()
	}
}

func (codec *boolCodec) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteBool(*((*bool)(ptr)))
}

func (codec *boolCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return !(*((*bool)(ptr)))
}

type emptyInterfaceCodec struct {
}

func (codec *emptyInterfaceCodec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	existing := *((*interface{})(ptr))

	// Checking for both typed and untyped nil pointers.
	if existing != nil &&
		reflect.TypeOf(existing).Kind() == reflect.Ptr &&
		!reflect.ValueOf(existing).IsNil() {

		var ptrToExisting interface{}
		for {
			elem := reflect.ValueOf(existing).Elem()
			if elem.Kind() != reflect.Ptr || elem.IsNil() {
				break
			}
			ptrToExisting = existing
			existing = elem.Interface()
		}

		if iter.ReadNil() {
			if ptrToExisting != nil {
				nilPtr := reflect.Zero(reflect.TypeOf(ptrToExisting).Elem())
				reflect.ValueOf(ptrToExisting).Elem().Set(nilPtr)
			} else {
				*((*interface{})(ptr)) = nil
			}
		} else {
			iter.ReadVal(existing)
		}

		return
	}

	if iter.ReadNil() {
		*((*interface{})(ptr)) = nil
	} else {
		*((*interface{})(ptr)) = iter.Read()
	}
}

func (codec *emptyInterfaceCodec) Encode(ptr unsafe.Pointer, stream *Stream) {
	obj := *((*interface{})(ptr))
	stream.WriteVal(obj)
}

func (codec *emptyInterfaceCodec) IsEmpty(ptr unsafe.Pointer) bool {
	emptyInterface := (*emptyInterface)(ptr)
	return emptyInterface.typ == nil
}

type nonEmptyInterfaceCodec struct {
}

func (codec *nonEmptyInterfaceCodec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if iter.WhatIsNext() == NilValue {
		iter.skipFourBytes('n', 'u', 'l', 'l')
		*((*interface{})(ptr)) = nil
		return
	}
	nonEmptyInterface := (*nonEmptyInterface)(ptr)
	if nonEmptyInterface.itab == nil {
		iter.ReportError("read non-empty interface", "do not know which concrete type to decode to")
		return
	}
	var i interface{}
	e := (*emptyInterface)(unsafe.Pointer(&i))
	e.typ = nonEmptyInterface.itab.typ
	e.word = nonEmptyInterface.word
	iter.ReadVal(&i)
	if e.word == nil {
		nonEmptyInterface.itab = nil
	}
	nonEmptyInterface.word = e.word
}

func (codec *nonEmptyInterfaceCodec) Encode(ptr unsafe.Pointer, stream *Stream) {
	nonEmptyInterface := (*nonEmptyInterface)(ptr)
	var i interface{}
	if nonEmptyInterface.itab != nil {
		e := (*emptyInterface)(unsafe.Pointer(&i))
		e.typ = nonEmptyInterface.itab.typ
		e.word = nonEmptyInterface.word
	}
	stream.WriteVal(i)
}

func (codec *nonEmptyInterfaceCodec) IsEmpty(ptr unsafe.Pointer) bool {
	nonEmptyInterface := (*nonEmptyInterface)(ptr)
	return nonEmptyInterface.word == nil
}

type dynamicEncoder struct {
	valType reflect2.Type
}

func (encoder *dynamicEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	obj := encoder.valType.UnsafeIndirect(ptr)
	stream.WriteVal(obj)
}

func (encoder *dynamicEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.valType.UnsafeIndirect(ptr) == nil
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

type jsoniterNumberCodec struct {
}

func (codec *jsoniterNumberCodec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	switch iter.WhatIsNext() {
	case StringValue:
		*((*Number)(ptr)) = Number(iter.ReadString())
	case NilValue:
		iter.skipFourBytes('n', 'u', 'l', 'l')
		*((*Number)(ptr)) = ""
	default:
		*((*Number)(ptr)) = Number([]byte(iter.readNumberAsString()))
	}
}

func (codec *jsoniterNumberCodec) Encode(ptr unsafe.Pointer, stream *Stream) {
	number := *((*Number)(ptr))
	if len(number) == 0 {
		stream.writeByte('0')
	} else {
		stream.WriteRaw(string(number))
	}
}

func (codec *jsoniterNumberCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return len(*((*Number)(ptr))) == 0
}

type jsonRawMessageCodec struct {
}

func (codec *jsonRawMessageCodec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*json.RawMessage)(ptr)) = json.RawMessage(iter.SkipAndReturnBytes())
}

func (codec *jsonRawMessageCodec) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteRaw(string(*((*json.RawMessage)(ptr))))
}

func (codec *jsonRawMessageCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return len(*((*json.RawMessage)(ptr))) == 0
}

type jsoniterRawMessageCodec struct {
}

func (codec *jsoniterRawMessageCodec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*RawMessage)(ptr)) = RawMessage(iter.SkipAndReturnBytes())
}

func (codec *jsoniterRawMessageCodec) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteRaw(string(*((*RawMessage)(ptr))))
}

func (codec *jsoniterRawMessageCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return len(*((*RawMessage)(ptr))) == 0
}

type base64Codec struct {
	sliceType *reflect2.UnsafeSliceType
	sliceDecoder ValDecoder
}

func (codec *base64Codec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if iter.ReadNil() {
		codec.sliceType.UnsafeSetNil(ptr)
		return
	}
	switch iter.WhatIsNext() {
	case StringValue:
		encoding := base64.StdEncoding
		src := iter.SkipAndReturnBytes()
		src = src[1: len(src)-1]
		decodedLen := encoding.DecodedLen(len(src))
		dst := make([]byte, decodedLen)
		len, err := encoding.Decode(dst, src)
		if err != nil {
			iter.ReportError("decode base64", err.Error())
		} else {
			dst = dst[:len]
			codec.sliceType.UnsafeSet(ptr, unsafe.Pointer(&dst))
		}
	case ArrayValue:
		codec.sliceDecoder.Decode(ptr, iter)
	default:
		iter.ReportError("base64Codec", "invalid input")
	}
}

func (codec *base64Codec) Encode(ptr unsafe.Pointer, stream *Stream) {
	src := *((*[]byte)(ptr))
	if len(src) == 0 {
		stream.WriteNil()
		return
	}
	encoding := base64.StdEncoding
	stream.writeByte('"')
	size := encoding.EncodedLen(len(src))
	buf := make([]byte, size)
	encoding.Encode(buf, src)
	stream.buf = append(stream.buf, buf...)
	stream.writeByte('"')
}

func (codec *base64Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return len(*((*[]byte)(ptr))) == 0
}

type stringModeNumberDecoder struct {
	elemDecoder ValDecoder
}

func (decoder *stringModeNumberDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	c := iter.nextToken()
	if c != '"' {
		iter.ReportError("stringModeNumberDecoder", `expect ", but found `+string([]byte{c}))
		return
	}
	decoder.elemDecoder.Decode(ptr, iter)
	if iter.Error != nil {
		return
	}
	c = iter.readByte()
	if c != '"' {
		iter.ReportError("stringModeNumberDecoder", `expect ", but found `+string([]byte{c}))
		return
	}
}

type stringModeStringDecoder struct {
	elemDecoder ValDecoder
	cfg         *frozenConfig
}

func (decoder *stringModeStringDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	decoder.elemDecoder.Decode(ptr, iter)
	str := *((*string)(ptr))
	tempIter := decoder.cfg.BorrowIterator([]byte(str))
	defer decoder.cfg.ReturnIterator(tempIter)
	*((*string)(ptr)) = tempIter.ReadString()
}

type stringModeNumberEncoder struct {
	elemEncoder ValEncoder
}

func (encoder *stringModeNumberEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.writeByte('"')
	encoder.elemEncoder.Encode(ptr, stream)
	stream.writeByte('"')
}

func (encoder *stringModeNumberEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.elemEncoder.IsEmpty(ptr)
}

type stringModeStringEncoder struct {
	elemEncoder ValEncoder
	cfg         *frozenConfig
}

func (encoder *stringModeStringEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	tempStream := encoder.cfg.BorrowStream(nil)
	defer encoder.cfg.ReturnStream(tempStream)
	encoder.elemEncoder.Encode(ptr, tempStream)
	stream.WriteString(string(tempStream.Buffer()))
}

func (encoder *stringModeStringEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.elemEncoder.IsEmpty(ptr)
}
