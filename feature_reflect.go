package jsoniter

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"sync/atomic"
	"unsafe"
)

/*
Reflection on type to create decoders, which is then cached
Reflection on value is avoided as we can, as the reflect.Value itself will allocate, with following exceptions
1. create instance of new value, for example *int will need a int to be allocated
2. append to slice, if the existing cap is not enough, allocate will be done using Reflect.New
3. assignment to map, both key and value will be reflect.Value
For a simple struct binding, it will be reflect.Value free and allocation free
*/

type Decoder interface {
	decode(ptr unsafe.Pointer, iter *Iterator)
}

type Encoder interface {
	isEmpty(ptr unsafe.Pointer) bool
	encode(ptr unsafe.Pointer, stream *Stream)
	encodeInterface(val interface{}, stream *Stream)
}

func writeToStream(val interface{}, stream *Stream, encoder Encoder) {
	e := (*emptyInterface)(unsafe.Pointer(&val))
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		encoder.encode(unsafe.Pointer(&e.word), stream)
	} else {
		encoder.encode(e.word, stream)
	}
}

type DecoderFunc func(ptr unsafe.Pointer, iter *Iterator)
type EncoderFunc func(ptr unsafe.Pointer, stream *Stream)
type ExtensionFunc func(typ reflect.Type, field *reflect.StructField) ([]string, EncoderFunc, DecoderFunc)

type funcDecoder struct {
	fun DecoderFunc
}

func (decoder *funcDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	decoder.fun(ptr, iter)
}

type funcEncoder struct {
	fun EncoderFunc
}

func (encoder *funcEncoder) encode(ptr unsafe.Pointer, stream *Stream) {
	encoder.fun(ptr, stream)
}

func (encoder *funcEncoder) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (encoder *funcEncoder) isEmpty(ptr unsafe.Pointer) bool {
	return false
}

var DECODERS unsafe.Pointer
var ENCODERS unsafe.Pointer

var typeDecoders map[string]Decoder
var fieldDecoders map[string]Decoder
var typeEncoders map[string]Encoder
var fieldEncoders map[string]Encoder
var extensions []ExtensionFunc
var jsonNumberType reflect.Type
var jsonRawMessageType reflect.Type
var anyType reflect.Type
var marshalerType reflect.Type
var unmarshalerType reflect.Type
var textUnmarshalerType reflect.Type

func init() {
	typeDecoders = map[string]Decoder{}
	RegisterTypeDecoder("time.Time", TimeDecoder)
	fieldDecoders = map[string]Decoder{}
	typeEncoders = map[string]Encoder{}
	fieldEncoders = map[string]Encoder{}
	extensions = []ExtensionFunc{}
	atomic.StorePointer(&DECODERS, unsafe.Pointer(&map[string]Decoder{}))
	atomic.StorePointer(&ENCODERS, unsafe.Pointer(&map[string]Encoder{}))
	jsonNumberType = reflect.TypeOf((*json.Number)(nil)).Elem()
	jsonRawMessageType = reflect.TypeOf((*json.RawMessage)(nil)).Elem()
	anyType = reflect.TypeOf((*Any)(nil)).Elem()
	marshalerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	unmarshalerType = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
}

func addDecoderToCache(cacheKey reflect.Type, decoder Decoder) {
	done := false
	for !done {
		ptr := atomic.LoadPointer(&DECODERS)
		cache := *(*map[reflect.Type]Decoder)(ptr)
		copied := map[reflect.Type]Decoder{}
		for k, v := range cache {
			copied[k] = v
		}
		copied[cacheKey] = decoder
		done = atomic.CompareAndSwapPointer(&DECODERS, ptr, unsafe.Pointer(&copied))
	}
}

func addEncoderToCache(cacheKey reflect.Type, encoder Encoder) {
	done := false
	for !done {
		ptr := atomic.LoadPointer(&ENCODERS)
		cache := *(*map[reflect.Type]Encoder)(ptr)
		copied := map[reflect.Type]Encoder{}
		for k, v := range cache {
			copied[k] = v
		}
		copied[cacheKey] = encoder
		done = atomic.CompareAndSwapPointer(&ENCODERS, ptr, unsafe.Pointer(&copied))
	}
}

func getDecoderFromCache(cacheKey reflect.Type) Decoder {
	ptr := atomic.LoadPointer(&DECODERS)
	cache := *(*map[reflect.Type]Decoder)(ptr)
	return cache[cacheKey]
}

func getEncoderFromCache(cacheKey reflect.Type) Encoder {
	ptr := atomic.LoadPointer(&ENCODERS)
	cache := *(*map[reflect.Type]Encoder)(ptr)
	return cache[cacheKey]
}

// RegisterTypeDecoder can register a type for json object
func RegisterTypeDecoder(typ string, fun DecoderFunc) {
	typeDecoders[typ] = &funcDecoder{fun}
}

// RegisterFieldDecoder can register a type for json field
func RegisterFieldDecoder(typ string, field string, fun DecoderFunc) {
	fieldDecoders[fmt.Sprintf("%s/%s", typ, field)] = &funcDecoder{fun}
}

func RegisterTypeEncoder(typ string, fun EncoderFunc) {
	typeEncoders[typ] = &funcEncoder{fun}
}

func RegisterFieldEncoder(typ string, field string, fun EncoderFunc) {
	fieldEncoders[fmt.Sprintf("%s/%s", typ, field)] = &funcEncoder{fun}
}

// RegisterExtension can register a custom extension
func RegisterExtension(extension ExtensionFunc) {
	extensions = append(extensions, extension)
}

// CleanDecoders cleans decoders registered or cached
func CleanDecoders() {
	typeDecoders = map[string]Decoder{}
	fieldDecoders = map[string]Decoder{}
	atomic.StorePointer(&DECODERS, unsafe.Pointer(&map[string]Decoder{}))
}

// CleanEncoders cleans decoders registered or cached
func CleanEncoders() {
	typeEncoders = map[string]Encoder{}
	fieldEncoders = map[string]Encoder{}
	atomic.StorePointer(&ENCODERS, unsafe.Pointer(&map[string]Encoder{}))
}

type optionalDecoder struct {
	valueType    reflect.Type
	valueDecoder Decoder
}

func (decoder *optionalDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	if iter.ReadNil() {
		*((*unsafe.Pointer)(ptr)) = nil
	} else {
		if *((*unsafe.Pointer)(ptr)) == nil {
			// pointer to null, we have to allocate memory to hold the value
			value := reflect.New(decoder.valueType)
			decoder.valueDecoder.decode(unsafe.Pointer(value.Pointer()), iter)
			*((*uintptr)(ptr)) = value.Pointer()
		} else {
			// reuse existing instance
			decoder.valueDecoder.decode(*((*unsafe.Pointer)(ptr)), iter)
		}
	}
}

type optionalEncoder struct {
	valueEncoder Encoder
}

func (encoder *optionalEncoder) encode(ptr unsafe.Pointer, stream *Stream) {
	if *((*unsafe.Pointer)(ptr)) == nil {
		stream.WriteNil()
	} else {
		encoder.valueEncoder.encode(*((*unsafe.Pointer)(ptr)), stream)
	}
}

func (encoder *optionalEncoder) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (encoder *optionalEncoder) isEmpty(ptr unsafe.Pointer) bool {
	if *((*unsafe.Pointer)(ptr)) == nil {
		return true
	} else {
		return encoder.valueEncoder.isEmpty(*((*unsafe.Pointer)(ptr)))
	}
}

type placeholderEncoder struct {
	valueEncoder Encoder
}

func (encoder *placeholderEncoder) encode(ptr unsafe.Pointer, stream *Stream) {
	encoder.valueEncoder.encode(ptr, stream)
}

func (encoder *placeholderEncoder) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (encoder *placeholderEncoder) isEmpty(ptr unsafe.Pointer) bool {
	return encoder.valueEncoder.isEmpty(ptr)
}

type placeholderDecoder struct {
	valueDecoder Decoder
}

func (decoder *placeholderDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	decoder.valueDecoder.decode(ptr, iter)
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

// Read converts an Iterator instance into go interface, same as json.Unmarshal
func (iter *Iterator) ReadVal(obj interface{}) {
	typ := reflect.TypeOf(obj)
	cacheKey := typ.Elem()
	cachedDecoder := getDecoderFromCache(cacheKey)
	if cachedDecoder == nil {
		decoder, err := decoderOfType(cacheKey)
		if err != nil {
			iter.Error = err
			return
		}
		cachedDecoder = decoder
		addDecoderToCache(cacheKey, decoder)
	}
	e := (*emptyInterface)(unsafe.Pointer(&obj))
	cachedDecoder.decode(e.word, iter)
}

func (stream *Stream) WriteVal(val interface{}) {
	if nil == val {
		stream.WriteNil()
		return
	}
	typ := reflect.TypeOf(val)
	cacheKey := typ
	cachedEncoder := getEncoderFromCache(cacheKey)
	if cachedEncoder == nil {
		encoder, err := encoderOfType(cacheKey)
		if err != nil {
			stream.Error = err
			return
		}
		cachedEncoder = encoder
		addEncoderToCache(cacheKey, encoder)
	}
	cachedEncoder.encodeInterface(val, stream)
}

type prefix string

func (p prefix) addToDecoder(decoder Decoder, err error) (Decoder, error) {
	if err != nil {
		return nil, fmt.Errorf("%s: %s", p, err.Error())
	}
	return decoder, err
}

func (p prefix) addToEncoder(encoder Encoder, err error) (Encoder, error) {
	if err != nil {
		return nil, fmt.Errorf("%s: %s", p, err.Error())
	}
	return encoder, err
}

func decoderOfType(typ reflect.Type) (Decoder, error) {
	typeName := typ.String()
	typeDecoder := typeDecoders[typeName]
	if typeDecoder != nil {
		return typeDecoder, nil
	}
	if typ.Kind() == reflect.Ptr {
		typeDecoder := typeDecoders[typ.Elem().String()]
		if typeDecoder != nil {
			return &optionalDecoder{typ.Elem(), typeDecoder}, nil
		}
	}
	cacheKey := typ
	cachedDecoder := getDecoderFromCache(cacheKey)
	if cachedDecoder != nil {
		return cachedDecoder, nil
	}
	placeholder := &placeholderDecoder{}
	addDecoderToCache(cacheKey, placeholder)
	newDecoder, err := createDecoderOfType(typ)
	placeholder.valueDecoder = newDecoder
	addDecoderToCache(cacheKey, newDecoder)
	return newDecoder, err
}

func createDecoderOfType(typ reflect.Type) (Decoder, error) {
	if typ.String() == "[]uint8" {
		return &base64Codec{}, nil
	}
	if typ.AssignableTo(jsonRawMessageType) {
		return &jsonRawMessageCodec{}, nil
	}
	if typ.AssignableTo(jsonNumberType) {
		return &jsonNumberCodec{}, nil
	}
	if typ.ConvertibleTo(unmarshalerType) {
		templateInterface := reflect.New(typ).Elem().Interface()
		return &optionalDecoder{typ, &unmarshalerDecoder{extractInterface(templateInterface)}}, nil
	}
	if typ.ConvertibleTo(anyType) {
		return &anyCodec{}, nil
	}
	switch typ.Kind() {
	case reflect.String:
		return &stringCodec{}, nil
	case reflect.Int:
		return &intCodec{}, nil
	case reflect.Int8:
		return &int8Codec{}, nil
	case reflect.Int16:
		return &int16Codec{}, nil
	case reflect.Int32:
		return &int32Codec{}, nil
	case reflect.Int64:
		return &int64Codec{}, nil
	case reflect.Uint:
		return &uintCodec{}, nil
	case reflect.Uint8:
		return &uint8Codec{}, nil
	case reflect.Uint16:
		return &uint16Codec{}, nil
	case reflect.Uint32:
		return &uint32Codec{}, nil
	case reflect.Uint64:
		return &uint64Codec{}, nil
	case reflect.Float32:
		return &float32Codec{}, nil
	case reflect.Float64:
		return &float64Codec{}, nil
	case reflect.Bool:
		return &boolCodec{}, nil
	case reflect.Interface:
		if typ.NumMethod() == 0 {
			return &emptyInterfaceCodec{}, nil
		} else {
			return &nonEmptyInterfaceCodec{}, nil
		}
	case reflect.Struct:
		return prefix(fmt.Sprintf("[%s]", typ.String())).addToDecoder(decoderOfStruct(typ))
	case reflect.Slice:
		return prefix("[slice]").addToDecoder(decoderOfSlice(typ))
	case reflect.Map:
		return prefix("[map]").addToDecoder(decoderOfMap(typ))
	case reflect.Ptr:
		return prefix("[optional]").addToDecoder(decoderOfOptional(typ))
	default:
		return nil, fmt.Errorf("unsupported type: %v", typ)
	}
}

func encoderOfType(typ reflect.Type) (Encoder, error) {
	typeName := typ.String()
	typeEncoder := typeEncoders[typeName]
	if typeEncoder != nil {
		return typeEncoder, nil
	}
	if typ.Kind() == reflect.Ptr {
		typeEncoder := typeEncoders[typ.Elem().String()]
		if typeEncoder != nil {
			return &optionalEncoder{typeEncoder}, nil
		}
	}
	cacheKey := typ
	cachedEncoder := getEncoderFromCache(cacheKey)
	if cachedEncoder != nil {
		return cachedEncoder, nil
	}
	placeholder := &placeholderEncoder{}
	addEncoderToCache(cacheKey, placeholder)
	newEncoder, err := createEncoderOfType(typ)
	placeholder.valueEncoder = newEncoder
	addEncoderToCache(cacheKey, newEncoder)
	return newEncoder, err
}

func createEncoderOfType(typ reflect.Type) (Encoder, error) {
	if typ.String() == "[]uint8" {
		return &base64Codec{}, nil
	}
	if typ.AssignableTo(jsonRawMessageType) {
		return &jsonRawMessageCodec{}, nil
	}
	if typ.AssignableTo(jsonNumberType) {
		return &jsonNumberCodec{}, nil
	}
	if typ.ConvertibleTo(marshalerType) {
		templateInterface := reflect.New(typ).Elem().Interface()
		return &marshalerEncoder{extractInterface(templateInterface)}, nil
	}
	if typ.ConvertibleTo(anyType) {
		return &anyCodec{}, nil
	}
	kind := typ.Kind()
	switch kind {
	case reflect.String:
		return &stringCodec{}, nil
	case reflect.Int:
		return &intCodec{}, nil
	case reflect.Int8:
		return &int8Codec{}, nil
	case reflect.Int16:
		return &int16Codec{}, nil
	case reflect.Int32:
		return &int32Codec{}, nil
	case reflect.Int64:
		return &int64Codec{}, nil
	case reflect.Uint:
		return &uintCodec{}, nil
	case reflect.Uint8:
		return &uint8Codec{}, nil
	case reflect.Uint16:
		return &uint16Codec{}, nil
	case reflect.Uint32:
		return &uint32Codec{}, nil
	case reflect.Uint64:
		return &uint64Codec{}, nil
	case reflect.Float32:
		return &float32Codec{}, nil
	case reflect.Float64:
		return &float64Codec{}, nil
	case reflect.Bool:
		return &boolCodec{}, nil
	case reflect.Interface:
		if typ.NumMethod() == 0 {
			return &emptyInterfaceCodec{}, nil
		} else {
			return &nonEmptyInterfaceCodec{}, nil
		}
	case reflect.Struct:
		return prefix(fmt.Sprintf("[%s]", typ.String())).addToEncoder(encoderOfStruct(typ))
	case reflect.Slice:
		return prefix("[slice]").addToEncoder(encoderOfSlice(typ))
	case reflect.Map:
		return prefix("[map]").addToEncoder(encoderOfMap(typ))
	case reflect.Ptr:
		return prefix("[optional]").addToEncoder(encoderOfOptional(typ))
	default:
		return nil, fmt.Errorf("unsupported type: %v", typ)
	}
}

func decoderOfOptional(typ reflect.Type) (Decoder, error) {
	elemType := typ.Elem()
	decoder, err := decoderOfType(elemType)
	if err != nil {
		return nil, err
	}
	return &optionalDecoder{elemType, decoder}, nil
}

func encoderOfOptional(typ reflect.Type) (Encoder, error) {
	elemType := typ.Elem()
	decoder, err := encoderOfType(elemType)
	if err != nil {
		return nil, err
	}
	return &optionalEncoder{decoder}, nil
}

func decoderOfMap(typ reflect.Type) (Decoder, error) {
	decoder, err := decoderOfType(typ.Elem())
	if err != nil {
		return nil, err
	}
	mapInterface := reflect.New(typ).Interface()
	return &mapDecoder{typ, typ.Key(), typ.Elem(), decoder, extractInterface(mapInterface)}, nil
}

func extractInterface(val interface{}) emptyInterface {
	return *((*emptyInterface)(unsafe.Pointer(&val)))
}

func encoderOfMap(typ reflect.Type) (Encoder, error) {
	elemType := typ.Elem()
	encoder, err := encoderOfType(elemType)
	if err != nil {
		return nil, err
	}
	mapInterface := reflect.New(typ).Elem().Interface()
	return &mapEncoder{typ, elemType, encoder, *((*emptyInterface)(unsafe.Pointer(&mapInterface)))}, nil
}
