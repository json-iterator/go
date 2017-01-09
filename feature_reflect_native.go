package jsoniter

import "unsafe"

type stringCodec struct {
}

func (codec *stringCodec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*string)(ptr)) = iter.ReadString()
}

func (codec *stringCodec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteString(*((*string)(ptr)))
}

type intCodec struct {
}

func (codec *intCodec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int)(ptr)) = iter.ReadInt()
}

func (codec *intCodec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteInt(*((*int)(ptr)))
}

type int8Codec struct {
}

func (codec *int8Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int8)(ptr)) = iter.ReadInt8()
}

func (codec *int8Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteInt8(*((*int8)(ptr)))
}

type int16Codec struct {
}

func (codec *int16Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int16)(ptr)) = iter.ReadInt16()
}

func (codec *int16Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteInt16(*((*int16)(ptr)))
}

type int32Codec struct {
}

func (codec *int32Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int32)(ptr)) = iter.ReadInt32()
}

func (codec *int32Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteInt32(*((*int32)(ptr)))
}

type int64Codec struct {
}

func (codec *int64Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int64)(ptr)) = iter.ReadInt64()
}

func (codec *int64Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteInt64(*((*int64)(ptr)))
}

type uintCodec struct {
}

func (codec *uintCodec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint)(ptr)) = iter.ReadUint()
}

func (codec *uintCodec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteUint(*((*uint)(ptr)))
}

type uint8Codec struct {
}

func (codec *uint8Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint8)(ptr)) = iter.ReadUint8()
}

func (codec *uint8Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteUint8(*((*uint8)(ptr)))
}

type uint16Codec struct {
}

func (decoder *uint16Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint16)(ptr)) = iter.ReadUint16()
}

func (decoder *uint16Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteUint16(*((*uint16)(ptr)))
}

type uint32Codec struct {
}

func (codec *uint32Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint32)(ptr)) = iter.ReadUint32()
}

func (codec *uint32Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteUint32(*((*uint32)(ptr)))
}

type uint64Codec struct {
}

func (codec *uint64Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint64)(ptr)) = iter.ReadUint64()
}

func (codec *uint64Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteUint64(*((*uint64)(ptr)))
}

type float32Codec struct {
}

func (codec *float32Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*float32)(ptr)) = iter.ReadFloat32()
}

func (codec *float32Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteFloat32(*((*float32)(ptr)))
}

type float64Codec struct {
}

func (codec *float64Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*float64)(ptr)) = iter.ReadFloat64()
}

func (codec *float64Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteFloat64(*((*float64)(ptr)))
}

type boolCodec struct {
}

func (codec *boolCodec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*bool)(ptr)) = iter.ReadBool()
}

func (codec *boolCodec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteBool(*((*bool)(ptr)))
}

type interfaceCodec struct {
}

func (codec *interfaceCodec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*interface{})(ptr)) = iter.ReadAny().Get()
}

func (codec *interfaceCodec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteVal(*((*interface{})(ptr)))
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