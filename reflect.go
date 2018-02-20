package jsoniter

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
	"unsafe"
	"github.com/v2pro/plz/reflect2"
)

// ValDecoder is an internal type registered to cache as needed.
// Don't confuse jsoniter.ValDecoder with json.Decoder.
// For json.Decoder's adapter, refer to jsoniter.AdapterDecoder(todo link).
//
// Reflection on type to create decoders, which is then cached
// Reflection on value is avoided as we can, as the reflect.Value itself will allocate, with following exceptions
// 1. create instance of new value, for example *int will need a int to be allocated
// 2. append to slice, if the existing cap is not enough, allocate will be done using Reflect.New
// 3. assignment to map, both key and value will be reflect.Value
// For a simple struct binding, it will be reflect.Value free and allocation free
type ValDecoder interface {
	Decode(ptr unsafe.Pointer, iter *Iterator)
}

// ValEncoder is an internal type registered to cache as needed.
// Don't confuse jsoniter.ValEncoder with json.Encoder.
// For json.Encoder's adapter, refer to jsoniter.AdapterEncoder(todo godoc link).
type ValEncoder interface {
	IsEmpty(ptr unsafe.Pointer) bool
	Encode(ptr unsafe.Pointer, stream *Stream)
}

type checkIsEmpty interface {
	IsEmpty(ptr unsafe.Pointer) bool
}

// ReadVal copy the underlying JSON into go interface, same as json.Unmarshal
func (iter *Iterator) ReadVal(obj interface{}) {
	typ := reflect.TypeOf(obj)
	cacheKey := typ.Elem()
	decoder := decoderOfType(iter.cfg, "", cacheKey)
	e := (*emptyInterface)(unsafe.Pointer(&obj))
	if e.word == nil {
		iter.ReportError("ReadVal", "can not read into nil pointer")
		return
	}
	decoder.Decode(e.word, iter)
}

// WriteVal copy the go interface into underlying JSON, same as json.Marshal
func (stream *Stream) WriteVal(val interface{}) {
	if nil == val {
		stream.WriteNil()
		return
	}
	typ := reflect.TypeOf(val)
	encoder := stream.cfg.EncoderOf(typ)
	encoder.Encode(reflect2.PtrOf(val), stream)
}

func (cfg *frozenConfig) DecoderOf(typ reflect.Type) ValDecoder {
	cacheKey := typ
	decoder := cfg.getDecoderFromCache(cacheKey)
	if decoder != nil {
		return decoder
	}
	decoder = decoderOfType(cfg, "", typ)
	cfg.addDecoderToCache(cacheKey, decoder)
	return decoder
}

func decoderOfType(cfg *frozenConfig, prefix string, typ reflect.Type) ValDecoder {
	decoder := getTypeDecoderFromExtension(cfg, typ)
	if decoder != nil {
		return decoder
	}
	decoder = createDecoderOfType(cfg, prefix, typ)
	for _, extension := range extensions {
		decoder = extension.DecorateDecoder(typ, decoder)
	}
	for _, extension := range cfg.extensions {
		decoder = extension.DecorateDecoder(typ, decoder)
	}
	return decoder
}

func createDecoderOfType(cfg *frozenConfig, prefix string, typ reflect.Type) ValDecoder {
	decoder := createDecoderOfJsonRawMessage(cfg, prefix, typ)
	if decoder != nil {
		return decoder
	}
	decoder = createDecoderOfJsonNumber(cfg, prefix, typ)
	if decoder != nil {
		return decoder
	}
	decoder = createDecoderOfMarshaler(cfg, prefix, typ)
	if decoder != nil {
		return decoder
	}
	decoder = createDecoderOfAny(cfg, prefix, typ)
	if decoder != nil {
		return decoder
	}
	decoder = createDecoderOfNative(cfg, prefix, typ)
	if decoder != nil {
		return decoder
	}
	switch typ.Kind() {
	case reflect.Interface:
		if typ.NumMethod() == 0 {
			return &emptyInterfaceCodec{}
		}
		return &nonEmptyInterfaceCodec{}
	case reflect.Struct:
		return decoderOfStruct(cfg, prefix, typ)
	case reflect.Array:
		return decoderOfArray(cfg, prefix, typ)
	case reflect.Slice:
		return decoderOfSlice(cfg, prefix, typ)
	case reflect.Map:
		return decoderOfMap(cfg, prefix, typ)
	case reflect.Ptr:
		return decoderOfOptional(cfg, prefix, typ)
	default:
		return &lazyErrorDecoder{err: fmt.Errorf("%s%s is unsupported type", prefix, typ.String())}
	}
}

func (cfg *frozenConfig) EncoderOf(typ reflect.Type) ValEncoder {
	cacheKey := typ
	encoder := cfg.getEncoderFromCache(cacheKey)
	if encoder != nil {
		return encoder
	}
	encoder = encoderOfType(cfg, "", typ)
	if shouldFixOnePtr(typ) {
		encoder = &onePtrEncoder{encoder}
	}
	cfg.addEncoderToCache(cacheKey, encoder)
	return encoder
}

type onePtrEncoder struct {
	encoder ValEncoder
}

func (encoder *onePtrEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.encoder.IsEmpty(unsafe.Pointer(&ptr))
}

func (encoder *onePtrEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	encoder.encoder.Encode(unsafe.Pointer(&ptr), stream)
}

func shouldFixOnePtr(typ reflect.Type) bool {
	return reflect2.Type2(typ).LikePtr()
}

func encoderOfType(cfg *frozenConfig, prefix string, typ reflect.Type) ValEncoder {
	encoder := getTypeEncoderFromExtension(cfg, typ)
	if encoder != nil {
		return encoder
	}
	encoder = createEncoderOfType(cfg, prefix, typ)
	for _, extension := range extensions {
		encoder = extension.DecorateEncoder(typ, encoder)
	}
	for _, extension := range cfg.extensions {
		encoder = extension.DecorateEncoder(typ, encoder)
	}
	return encoder
}

func createEncoderOfType(cfg *frozenConfig, prefix string, typ reflect.Type) ValEncoder {
	encoder := createEncoderOfJsonRawMessage(cfg, prefix, typ)
	if encoder != nil {
		return encoder
	}
	encoder = createEncoderOfJsonNumber(cfg, prefix, typ)
	if encoder != nil {
		return encoder
	}
	encoder = createEncoderOfMarshaler(cfg, prefix, typ)
	if encoder != nil {
		return encoder
	}
	encoder = createEncoderOfAny(cfg, prefix, typ)
	if encoder != nil {
		return encoder
	}
	encoder = createEncoderOfNative(cfg, prefix, typ)
	if encoder != nil {
		return encoder
	}
	kind := typ.Kind()
	switch kind {
	case reflect.Interface:
		return &dynamicEncoder{reflect2.Type2(typ)}
	case reflect.Struct:
		return encoderOfStruct(cfg, prefix, typ)
	case reflect.Array:
		return encoderOfArray(cfg, prefix, typ)
	case reflect.Slice:
		return encoderOfSlice(cfg, prefix, typ)
	case reflect.Map:
		return encoderOfMap(cfg, prefix, typ)
	case reflect.Ptr:
		return encoderOfOptional(cfg, prefix, typ)
	default:
		return &lazyErrorEncoder{err: fmt.Errorf("%s%s is unsupported type", prefix, typ.String())}
	}
}

type lazyErrorDecoder struct {
	err error
}

func (decoder *lazyErrorDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if iter.WhatIsNext() != NilValue {
		if iter.Error == nil {
			iter.Error = decoder.err
		}
	} else {
		iter.Skip()
	}
}

type lazyErrorEncoder struct {
	err error
}

func (encoder *lazyErrorEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	if ptr == nil {
		stream.WriteNil()
	} else if stream.Error == nil {
		stream.Error = encoder.err
	}
}

func (encoder *lazyErrorEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func extractInterface(val interface{}) emptyInterface {
	return *((*emptyInterface)(unsafe.Pointer(&val)))
}

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  unsafe.Pointer
	word unsafe.Pointer
}

// emptyInterface is the header for an interface with method (not interface{})
type nonEmptyInterface struct {
	// see ../runtime/iface.go:/Itab
	itab *struct {
		ityp   unsafe.Pointer // static interface type
		typ    unsafe.Pointer // dynamic concrete type
		link   unsafe.Pointer
		bad    int32
		unused int32
		fun    [100000]unsafe.Pointer // method table
	}
	word unsafe.Pointer
}
