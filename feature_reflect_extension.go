package jsoniter

import (
	"reflect"
	"fmt"
	"unsafe"
)

var typeDecoders map[string]ValDecoder
var fieldDecoders map[string]ValDecoder
var typeEncoders map[string]ValEncoder
var fieldEncoders map[string]ValEncoder
var extensions []ExtensionFunc

type ExtensionFunc func(typ reflect.Type, field *reflect.StructField) ([]string, EncoderFunc, DecoderFunc)

type funcDecoder struct {
	fun DecoderFunc
}

func init() {
	typeDecoders = map[string]ValDecoder{}
	fieldDecoders = map[string]ValDecoder{}
	typeEncoders = map[string]ValEncoder{}
	fieldEncoders = map[string]ValEncoder{}
	extensions = []ExtensionFunc{}
}

func RegisterTypeDecoderFunc(typ string, fun DecoderFunc) {
	typeDecoders[typ] = &funcDecoder{fun}
}

func RegisterTypeDecoder(typ string, decoder ValDecoder) {
	typeDecoders[typ] = decoder
}

func RegisterFieldDecoderFunc(typ string, field string, fun DecoderFunc) {
	RegisterFieldDecoder(typ, field, &funcDecoder{fun})
}

func RegisterFieldDecoder(typ string, field string, decoder ValDecoder) {
	fieldDecoders[fmt.Sprintf("%s/%s", typ, field)] = decoder
}

func RegisterTypeEncoderFunc(typ string, fun EncoderFunc, isEmptyFunc func(unsafe.Pointer) bool) {
	typeEncoders[typ] = &funcEncoder{fun, isEmptyFunc}
}

func RegisterTypeEncoder(typ string, encoder ValEncoder) {
	typeEncoders[typ] = encoder
}

func RegisterFieldEncoderFunc(typ string, field string, fun EncoderFunc, isEmptyFunc func(unsafe.Pointer) bool) {
	RegisterFieldEncoder(typ, field, &funcEncoder{fun, isEmptyFunc})
}

func RegisterFieldEncoder(typ string, field string, encoder ValEncoder) {
	fieldEncoders[fmt.Sprintf("%s/%s", typ, field)] = encoder
}

func RegisterExtension(extension ExtensionFunc) {
	extensions = append(extensions, extension)
}

func getTypeDecoderFromExtension(typ reflect.Type) ValDecoder {
	typeName := typ.String()
	typeDecoder := typeDecoders[typeName]
	if typeDecoder != nil {
		return typeDecoder
	}
	if typ.Kind() == reflect.Ptr {
		typeDecoder := typeDecoders[typ.Elem().String()]
		if typeDecoder != nil {
			return &optionalDecoder{typ.Elem(), typeDecoder}
		}
	}
	return nil
}

func getTypeEncoderFromExtension(typ reflect.Type) ValEncoder {
	typeName := typ.String()
	typeEncoder := typeEncoders[typeName]
	if typeEncoder != nil {
		return typeEncoder
	}
	if typ.Kind() == reflect.Ptr {
		typeEncoder := typeEncoders[typ.Elem().String()]
		if typeEncoder != nil {
			return &optionalEncoder{typeEncoder}
		}
	}
	return nil
}