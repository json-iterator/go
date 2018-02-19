package jsoniter

import (
	"github.com/v2pro/plz/reflect2"
	"unsafe"
	"encoding"
	"encoding/json"
	"reflect"
)

func createEncoderOfMarshaler(cfg *frozenConfig, prefix string, typ reflect.Type) ValEncoder {
	if typ == marshalerType {
		checkIsEmpty := createCheckIsEmpty(cfg, typ)
		var encoder ValEncoder = &directMarshalerEncoder{
			checkIsEmpty: checkIsEmpty,
		}
		return encoder
	}
	if typ.Implements(marshalerType) {
		checkIsEmpty := createCheckIsEmpty(cfg, typ)
		var encoder ValEncoder = &marshalerEncoder{
			valType:      reflect2.Type2(typ),
			checkIsEmpty: checkIsEmpty,
		}
		return encoder
	}
	ptrType := reflect.PtrTo(typ)
	if prefix != "" && ptrType.Implements(marshalerType) {
		checkIsEmpty := createCheckIsEmpty(cfg, ptrType)
		var encoder ValEncoder = &marshalerEncoder{
			valType:      reflect2.Type2(ptrType),
			checkIsEmpty: checkIsEmpty,
		}
		return &referenceEncoder{encoder}
	}
	if typ == textMarshalerType {
		checkIsEmpty := createCheckIsEmpty(cfg, typ)
		var encoder ValEncoder = &directTextMarshalerEncoder{
			checkIsEmpty:  checkIsEmpty,
			stringEncoder: cfg.EncoderOf(reflect.TypeOf("")),
		}
		return encoder
	}
	if typ.Implements(textMarshalerType) {
		checkIsEmpty := createCheckIsEmpty(cfg, typ)
		var encoder ValEncoder = &textMarshalerEncoder{
			valType:       reflect2.Type2(typ),
			stringEncoder: cfg.EncoderOf(reflect.TypeOf("")),
			checkIsEmpty:  checkIsEmpty,
		}
		return encoder
	}
	// if prefix is empty, the type is the root type
	if prefix != "" && ptrType.Implements(textMarshalerType) {
		checkIsEmpty := createCheckIsEmpty(cfg, ptrType)
		var encoder ValEncoder = &textMarshalerEncoder{
			valType:       reflect2.Type2(ptrType),
			stringEncoder: cfg.EncoderOf(reflect.TypeOf("")),
			checkIsEmpty:  checkIsEmpty,
		}
		return &referenceEncoder{encoder}
	}
	return nil
}

type marshalerEncoder struct {
	checkIsEmpty checkIsEmpty
	valType      reflect2.Type
}

func (encoder *marshalerEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	obj := encoder.valType.UnsafeIndirect(ptr)
	if encoder.valType.IsNullable() && reflect2.IsNil(obj) {
		stream.WriteNil()
		return
	}
	marshaler := obj.(json.Marshaler)
	bytes, err := marshaler.MarshalJSON()
	if err != nil {
		stream.Error = err
	} else {
		stream.Write(bytes)
	}
}

func (encoder *marshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.checkIsEmpty.IsEmpty(ptr)
}

type directMarshalerEncoder struct {
	checkIsEmpty checkIsEmpty
}

func (encoder *directMarshalerEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	marshaler := *(*json.Marshaler)(ptr)
	if marshaler == nil {
		stream.WriteNil()
		return
	}
	bytes, err := marshaler.MarshalJSON()
	if err != nil {
		stream.Error = err
	} else {
		stream.Write(bytes)
	}
}

func (encoder *directMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.checkIsEmpty.IsEmpty(ptr)
}

type textMarshalerEncoder struct {
	valType       reflect2.Type
	stringEncoder ValEncoder
	checkIsEmpty  checkIsEmpty
}

func (encoder *textMarshalerEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	obj := encoder.valType.UnsafeIndirect(ptr)
	if encoder.valType.IsNullable() && reflect2.IsNil(obj) {
		stream.WriteNil()
		return
	}
	marshaler := (obj).(encoding.TextMarshaler)
	bytes, err := marshaler.MarshalText()
	if err != nil {
		stream.Error = err
	} else {
		str := string(bytes)
		encoder.stringEncoder.Encode(unsafe.Pointer(&str), stream)
	}
}

func (encoder *textMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.checkIsEmpty.IsEmpty(ptr)
}

type directTextMarshalerEncoder struct {
	stringEncoder ValEncoder
	checkIsEmpty  checkIsEmpty
}

func (encoder *directTextMarshalerEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	marshaler := *(*encoding.TextMarshaler)(ptr)
	if marshaler == nil {
		stream.WriteNil()
		return
	}
	bytes, err := marshaler.MarshalText()
	if err != nil {
		stream.Error = err
	} else {
		str := string(bytes)
		encoder.stringEncoder.Encode(unsafe.Pointer(&str), stream)
	}
}

func (encoder *directTextMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.checkIsEmpty.IsEmpty(ptr)
}

type unmarshalerDecoder struct {
	templateInterface emptyInterface
}

func (decoder *unmarshalerDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	templateInterface := decoder.templateInterface
	templateInterface.word = ptr
	realInterface := (*interface{})(unsafe.Pointer(&templateInterface))
	unmarshaler := (*realInterface).(json.Unmarshaler)
	iter.nextToken()
	iter.unreadByte() // skip spaces
	bytes := iter.SkipAndReturnBytes()
	err := unmarshaler.UnmarshalJSON(bytes)
	if err != nil {
		iter.ReportError("unmarshalerDecoder", err.Error())
	}
}

type textUnmarshalerDecoder struct {
	templateInterface emptyInterface
}

func (decoder *textUnmarshalerDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	templateInterface := decoder.templateInterface
	templateInterface.word = ptr
	realInterface := (*interface{})(unsafe.Pointer(&templateInterface))
	unmarshaler := (*realInterface).(encoding.TextUnmarshaler)
	str := iter.ReadString()
	err := unmarshaler.UnmarshalText([]byte(str))
	if err != nil {
		iter.ReportError("textUnmarshalerDecoder", err.Error())
	}
}
