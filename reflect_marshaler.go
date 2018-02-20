package jsoniter

import (
	"github.com/v2pro/plz/reflect2"
	"unsafe"
	"encoding"
	"encoding/json"
	"reflect"
)

var marshalerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
var unmarshalerType = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
var textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
var textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

func createDecoderOfMarshaler(cfg *frozenConfig, prefix string, typ reflect.Type) ValDecoder {
	ptrType := reflect.PtrTo(typ)
	if ptrType.Implements(unmarshalerType) {
		return &referenceDecoder{
			&unmarshalerDecoder{reflect2.Type2(ptrType)},
		}
	}
	if ptrType.Implements(textUnmarshalerType) {
		return &referenceDecoder{
			&textUnmarshalerDecoder{reflect2.Type2(ptrType)},
		}
	}
	return nil
}

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
	valType reflect2.Type
}

func (decoder *unmarshalerDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	valType := decoder.valType
	obj := valType.UnsafeIndirect(ptr)
	unmarshaler := obj.(json.Unmarshaler)
	iter.nextToken()
	iter.unreadByte() // skip spaces
	bytes := iter.SkipAndReturnBytes()
	err := unmarshaler.UnmarshalJSON(bytes)
	if err != nil {
		iter.ReportError("unmarshalerDecoder", err.Error())
	}
}

type textUnmarshalerDecoder struct {
	valType reflect2.Type
}

func (decoder *textUnmarshalerDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	valType := decoder.valType
	obj := valType.UnsafeIndirect(ptr)
	if reflect2.IsNil(obj) {
		ptrType := valType.(*reflect2.UnsafePtrType)
		elemType := ptrType.Elem()
		elem := elemType.UnsafeNew()
		ptrType.UnsafeSet(ptr, unsafe.Pointer(&elem))
		obj = valType.UnsafeIndirect(ptr)
	}
	unmarshaler := (obj).(encoding.TextUnmarshaler)
	str := iter.ReadString()
	err := unmarshaler.UnmarshalText([]byte(str))
	if err != nil {
		iter.ReportError("textUnmarshalerDecoder", err.Error())
	}
}
