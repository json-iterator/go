package jsoniter

import (
	"encoding/base64"
	"encoding/json"
	"unsafe"
)

type stringCodec struct {
}

func (codec *stringCodec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*string)(ptr)) = iter.ReadString()
}

func (codec *stringCodec) encode(ptr unsafe.Pointer, stream *Stream) {
	str := *((*string)(ptr))
	stream.WriteString(str)
}

func (codec *stringCodec) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, codec)
}

func (codec *stringCodec) isEmpty(ptr unsafe.Pointer) bool {
	return *((*string)(ptr)) == ""
}

type intCodec struct {
}

func (codec *intCodec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int)(ptr)) = iter.ReadInt()
}

func (codec *intCodec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteInt(*((*int)(ptr)))
}

func (encoder *intCodec) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (codec *intCodec) isEmpty(ptr unsafe.Pointer) bool {
	return *((*int)(ptr)) == 0
}

type int8Codec struct {
}

func (codec *int8Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int8)(ptr)) = iter.ReadInt8()
}

func (codec *int8Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteInt8(*((*int8)(ptr)))
}

func (encoder *int8Codec) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (codec *int8Codec) isEmpty(ptr unsafe.Pointer) bool {
	return *((*int8)(ptr)) == 0
}

type int16Codec struct {
}

func (codec *int16Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int16)(ptr)) = iter.ReadInt16()
}

func (codec *int16Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteInt16(*((*int16)(ptr)))
}

func (encoder *int16Codec) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (codec *int16Codec) isEmpty(ptr unsafe.Pointer) bool {
	return *((*int16)(ptr)) == 0
}

type int32Codec struct {
}

func (codec *int32Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int32)(ptr)) = iter.ReadInt32()
}

func (codec *int32Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteInt32(*((*int32)(ptr)))
}

func (encoder *int32Codec) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (codec *int32Codec) isEmpty(ptr unsafe.Pointer) bool {
	return *((*int32)(ptr)) == 0
}

type int64Codec struct {
}

func (codec *int64Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int64)(ptr)) = iter.ReadInt64()
}

func (codec *int64Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteInt64(*((*int64)(ptr)))
}

func (encoder *int64Codec) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (codec *int64Codec) isEmpty(ptr unsafe.Pointer) bool {
	return *((*int64)(ptr)) == 0
}

type uintCodec struct {
}

func (codec *uintCodec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint)(ptr)) = iter.ReadUint()
}

func (codec *uintCodec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteUint(*((*uint)(ptr)))
}

func (encoder *uintCodec) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (codec *uintCodec) isEmpty(ptr unsafe.Pointer) bool {
	return *((*uint)(ptr)) == 0
}

type uint8Codec struct {
}

func (codec *uint8Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint8)(ptr)) = iter.ReadUint8()
}

func (codec *uint8Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteUint8(*((*uint8)(ptr)))
}

func (encoder *uint8Codec) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (codec *uint8Codec) isEmpty(ptr unsafe.Pointer) bool {
	return *((*uint8)(ptr)) == 0
}

type uint16Codec struct {
}

func (decoder *uint16Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint16)(ptr)) = iter.ReadUint16()
}

func (codec *uint16Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteUint16(*((*uint16)(ptr)))
}

func (encoder *uint16Codec) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (codec *uint16Codec) isEmpty(ptr unsafe.Pointer) bool {
	return *((*uint16)(ptr)) == 0
}

type uint32Codec struct {
}

func (codec *uint32Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint32)(ptr)) = iter.ReadUint32()
}

func (codec *uint32Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteUint32(*((*uint32)(ptr)))
}

func (encoder *uint32Codec) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (codec *uint32Codec) isEmpty(ptr unsafe.Pointer) bool {
	return *((*uint32)(ptr)) == 0
}

type uint64Codec struct {
}

func (codec *uint64Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint64)(ptr)) = iter.ReadUint64()
}

func (codec *uint64Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteUint64(*((*uint64)(ptr)))
}

func (encoder *uint64Codec) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (codec *uint64Codec) isEmpty(ptr unsafe.Pointer) bool {
	return *((*uint64)(ptr)) == 0
}

type float32Codec struct {
}

func (codec *float32Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*float32)(ptr)) = iter.ReadFloat32()
}

func (codec *float32Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteFloat32(*((*float32)(ptr)))
}

func (encoder *float32Codec) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (codec *float32Codec) isEmpty(ptr unsafe.Pointer) bool {
	return *((*float32)(ptr)) == 0
}

type float64Codec struct {
}

func (codec *float64Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*float64)(ptr)) = iter.ReadFloat64()
}

func (codec *float64Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteFloat64(*((*float64)(ptr)))
}

func (encoder *float64Codec) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (codec *float64Codec) isEmpty(ptr unsafe.Pointer) bool {
	return *((*float64)(ptr)) == 0
}

type boolCodec struct {
}

func (codec *boolCodec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*bool)(ptr)) = iter.ReadBool()
}

func (codec *boolCodec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteBool(*((*bool)(ptr)))
}

func (encoder *boolCodec) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (codec *boolCodec) isEmpty(ptr unsafe.Pointer) bool {
	return !(*((*bool)(ptr)))
}

type emptyInterfaceCodec struct {
}

func (codec *emptyInterfaceCodec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*interface{})(ptr)) = iter.Read()
}

func (codec *emptyInterfaceCodec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteVal(*((*interface{})(ptr)))
}

func (encoder *emptyInterfaceCodec) encodeInterface(val interface{}, stream *Stream) {
	stream.WriteVal(val)
}

func (codec *emptyInterfaceCodec) isEmpty(ptr unsafe.Pointer) bool {
	return ptr == nil
}

type nonEmptyInterfaceCodec struct {
}

func (codec *nonEmptyInterfaceCodec) decode(ptr unsafe.Pointer, iter *Iterator) {
	nonEmptyInterface := (*nonEmptyInterface)(ptr)
	if nonEmptyInterface.itab == nil {
		iter.reportError("read non-empty interface", "do not know which concrete type to decode to")
		return
	}
	var i interface{}
	e := (*emptyInterface)(unsafe.Pointer(&i))
	e.typ = nonEmptyInterface.itab.typ
	e.word = nonEmptyInterface.word
	iter.ReadVal(&i)
	nonEmptyInterface.word = e.word
}

func (codec *nonEmptyInterfaceCodec) encode(ptr unsafe.Pointer, stream *Stream) {
	nonEmptyInterface := (*nonEmptyInterface)(ptr)
	var i interface{}
	e := (*emptyInterface)(unsafe.Pointer(&i))
	e.typ = nonEmptyInterface.itab.typ
	e.word = nonEmptyInterface.word
	stream.WriteVal(i)
}

func (encoder *nonEmptyInterfaceCodec) encodeInterface(val interface{}, stream *Stream) {
	stream.WriteVal(val)
}

func (codec *nonEmptyInterfaceCodec) isEmpty(ptr unsafe.Pointer) bool {
	nonEmptyInterface := (*nonEmptyInterface)(ptr)
	return nonEmptyInterface.word == nil
}

type anyCodec struct {
}

func (codec *anyCodec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*Any)(ptr)) = iter.ReadAny()
}

func (codec *anyCodec) encode(ptr unsafe.Pointer, stream *Stream) {
	(*((*Any)(ptr))).WriteTo(stream)
}

func (encoder *anyCodec) encodeInterface(val interface{}, stream *Stream) {
	(val.(Any)).WriteTo(stream)
}

func (encoder *anyCodec) isEmpty(ptr unsafe.Pointer) bool {
	return (*((*Any)(ptr))).Size() == 0
}

type jsonNumberCodec struct {
}

func (codec *jsonNumberCodec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*json.Number)(ptr)) = json.Number([]byte(iter.readNumberAsString()))
}

func (codec *jsonNumberCodec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteRaw(string(*((*json.Number)(ptr))))
}

func (encoder *jsonNumberCodec) encodeInterface(val interface{}, stream *Stream) {
	stream.WriteRaw(string(val.(json.Number)))
}

func (encoder *jsonNumberCodec) isEmpty(ptr unsafe.Pointer) bool {
	return len(*((*json.Number)(ptr))) == 0
}

type jsonRawMessageCodec struct {
}

func (codec *jsonRawMessageCodec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*json.RawMessage)(ptr)) = json.RawMessage(iter.SkipAndReturnBytes())
}

func (codec *jsonRawMessageCodec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteRaw(string(*((*json.RawMessage)(ptr))))
}

func (encoder *jsonRawMessageCodec) encodeInterface(val interface{}, stream *Stream) {
	stream.WriteRaw(string(val.(json.RawMessage)))
}

func (encoder *jsonRawMessageCodec) isEmpty(ptr unsafe.Pointer) bool {
	return len(*((*json.RawMessage)(ptr))) == 0
}

type jsoniterRawMessageCodec struct {
}

func (codec *jsoniterRawMessageCodec) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*RawMessage)(ptr)) = RawMessage(iter.SkipAndReturnBytes())
}

func (codec *jsoniterRawMessageCodec) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteRaw(string(*((*RawMessage)(ptr))))
}

func (encoder *jsoniterRawMessageCodec) encodeInterface(val interface{}, stream *Stream) {
	stream.WriteRaw(string(val.(RawMessage)))
}

func (encoder *jsoniterRawMessageCodec) isEmpty(ptr unsafe.Pointer) bool {
	return len(*((*RawMessage)(ptr))) == 0
}

type base64Codec struct {
}

func (codec *base64Codec) decode(ptr unsafe.Pointer, iter *Iterator) {
	encoding := base64.StdEncoding
	src := iter.SkipAndReturnBytes()
	src = src[1 : len(src)-1]
	decodedLen := encoding.DecodedLen(len(src))
	dst := make([]byte, decodedLen)
	_, err := encoding.Decode(dst, src)
	if err != nil {
		iter.reportError("decode base64", err.Error())
	} else {
		*((*[]byte)(ptr)) = dst
	}
}

func (codec *base64Codec) encode(ptr unsafe.Pointer, stream *Stream) {
	encoding := base64.StdEncoding
	stream.writeByte('"')
	src := *((*[]byte)(ptr))
	toGrow := encoding.EncodedLen(len(src))
	stream.ensure(toGrow)
	encoding.Encode(stream.buf[stream.n:], src)
	stream.n += toGrow
	stream.writeByte('"')
}

func (encoder *base64Codec) encodeInterface(val interface{}, stream *Stream) {
	encoding := base64.StdEncoding
	stream.writeByte('"')
	src := val.([]byte)
	toGrow := encoding.EncodedLen(len(src))
	stream.ensure(toGrow)
	encoding.Encode(stream.buf[stream.n:], src)
	stream.n += toGrow
	stream.writeByte('"')
}

func (encoder *base64Codec) isEmpty(ptr unsafe.Pointer) bool {
	return len(*((*[]byte)(ptr))) == 0
}

type stringModeDecoder struct {
	elemDecoder Decoder
}

func (decoder *stringModeDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	c := iter.nextToken()
	if c != '"' {
		iter.reportError("stringModeDecoder", `expect "`)
		return
	}
	decoder.elemDecoder.decode(ptr, iter)
	if iter.Error != nil {
		return
	}
	c = iter.readByte()
	if c != '"' {
		iter.reportError("stringModeDecoder", `expect "`)
		return
	}
}

type stringModeEncoder struct {
	elemEncoder Encoder
}

func (encoder *stringModeEncoder) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.writeByte('"')
	encoder.elemEncoder.encode(ptr, stream)
	stream.writeByte('"')
}

func (encoder *stringModeEncoder) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (encoder *stringModeEncoder) isEmpty(ptr unsafe.Pointer) bool {
	return encoder.elemEncoder.isEmpty(ptr)
}

type marshalerEncoder struct {
	templateInterface emptyInterface
}

func (encoder *marshalerEncoder) encode(ptr unsafe.Pointer, stream *Stream) {
	templateInterface := encoder.templateInterface
	templateInterface.word = ptr
	realInterface := (*interface{})(unsafe.Pointer(&templateInterface))
	marshaler := (*realInterface).(json.Marshaler)
	bytes, err := marshaler.MarshalJSON()
	if err != nil {
		stream.Error = err
	} else {
		stream.Write(bytes)
	}
}
func (encoder *marshalerEncoder) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (encoder *marshalerEncoder) isEmpty(ptr unsafe.Pointer) bool {
	templateInterface := encoder.templateInterface
	templateInterface.word = ptr
	realInterface := (*interface{})(unsafe.Pointer(&templateInterface))
	marshaler := (*realInterface).(json.Marshaler)
	bytes, err := marshaler.MarshalJSON()
	if err != nil {
		return true
	} else {
		return len(bytes) > 0
	}
}

type unmarshalerDecoder struct {
	templateInterface emptyInterface
}

func (decoder *unmarshalerDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	templateInterface := decoder.templateInterface
	templateInterface.word = ptr
	realInterface := (*interface{})(unsafe.Pointer(&templateInterface))
	unmarshaler := (*realInterface).(json.Unmarshaler)
	bytes := iter.SkipAndReturnBytes()
	err := unmarshaler.UnmarshalJSON(bytes)
	if err != nil {
		iter.reportError("unmarshaler", err.Error())
	}
}
