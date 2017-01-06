package jsoniter

import "unsafe"

type stringDecoder struct {
}

func (decoder *stringDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*string)(ptr)) = iter.ReadString()
}

type intDecoder struct {
}

func (decoder *intDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int)(ptr)) = iter.ReadInt()
}

type int8Decoder struct {
}

func (decoder *int8Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int8)(ptr)) = iter.ReadInt8()
}

type int16Decoder struct {
}

func (decoder *int16Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int16)(ptr)) = iter.ReadInt16()
}

type int32Decoder struct {
}

func (decoder *int32Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int32)(ptr)) = iter.ReadInt32()
}

type int64Decoder struct {
}

func (decoder *int64Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int64)(ptr)) = iter.ReadInt64()
}

type uintDecoder struct {
}

func (decoder *uintDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint)(ptr)) = iter.ReadUint()
}

type uint8Decoder struct {
}

func (decoder *uint8Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint8)(ptr)) = iter.ReadUint8()
}

type uint16Decoder struct {
}

func (decoder *uint16Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint16)(ptr)) = iter.ReadUint16()
}

type uint32Decoder struct {
}

func (decoder *uint32Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint32)(ptr)) = iter.ReadUint32()
}

type uint64Decoder struct {
}

func (decoder *uint64Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint64)(ptr)) = iter.ReadUint64()
}

type float32Decoder struct {
}

func (decoder *float32Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*float32)(ptr)) = iter.ReadFloat32()
}

type float64Decoder struct {
}

func (decoder *float64Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*float64)(ptr)) = iter.ReadFloat64()
}

type boolDecoder struct {
}

func (decoder *boolDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*bool)(ptr)) = iter.ReadBool()
}

type interfaceDecoder struct {
}

func (decoder *interfaceDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*interface{})(ptr)) = iter.ReadAny().Get()
}

type anyDecoder struct {
}

func (decoder *anyDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*Any)(ptr)) = *iter.ReadAny()
}

type stringNumberDecoder struct {
	elemDecoder Decoder
}

func (decoder *stringNumberDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	c := iter.nextToken()
	if c != '"' {
		iter.reportError("stringNumberDecoder", `expect "`)
		return
	}
	decoder.elemDecoder.decode(ptr, iter)
	if iter.Error != nil {
		return
	}
	c = iter.readByte()
	if c != '"' {
		iter.reportError("stringNumberDecoder", `expect "`)
		return
	}
}