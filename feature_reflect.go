package jsoniter

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
	"unsafe"
)

// Decoder is an internal type registered to cache as needed.
// Don't confuse jsoniter.Decoder with json.Decoder.
// For json.Decoder's adapter, refer to jsoniter.AdapterDecoder(todo link).
//
// Reflection on type to create decoders, which is then cached
// Reflection on value is avoided as we can, as the reflect.Value itself will allocate, with following exceptions
// 1. create instance of new value, for example *int will need a int to be allocated
// 2. append to slice, if the existing cap is not enough, allocate will be done using Reflect.New
// 3. assignment to map, both key and value will be reflect.Value
// For a simple struct binding, it will be reflect.Value free and allocation free
type Decoder interface {
	decode(ptr unsafe.Pointer, iter *Iterator)
}

// Encoder is an internal type registered to cache as needed.
// Don't confuse jsoniter.Encoder with json.Encoder.
// For json.Encoder's adapter, refer to jsoniter.AdapterEncoder(todo godoc link).
type Encoder interface {
	isEmpty(ptr unsafe.Pointer) bool
	encode(ptr unsafe.Pointer, stream *Stream)
	encodeInterface(val interface{}, stream *Stream)
}

func writeToStream(val interface{}, stream *Stream, encoder Encoder) {
	e := (*emptyInterface)(unsafe.Pointer(&val))
	if e.word == nil {
		stream.WriteNil()
		return
	}
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

var typeDecoders map[string]Decoder
var fieldDecoders map[string]Decoder
var typeEncoders map[string]Encoder
var fieldEncoders map[string]Encoder
var extensions []ExtensionFunc
var jsonNumberType reflect.Type
var jsonRawMessageType reflect.Type
var jsoniterRawMessageType reflect.Type
var anyType reflect.Type
var marshalerType reflect.Type
var unmarshalerType reflect.Type
var textUnmarshalerType reflect.Type

func init() {
	typeDecoders = map[string]Decoder{}
	fieldDecoders = map[string]Decoder{}
	typeEncoders = map[string]Encoder{}
	fieldEncoders = map[string]Encoder{}
	extensions = []ExtensionFunc{}
	jsonNumberType = reflect.TypeOf((*json.Number)(nil)).Elem()
	jsonRawMessageType = reflect.TypeOf((*json.RawMessage)(nil)).Elem()
	jsoniterRawMessageType = reflect.TypeOf((*RawMessage)(nil)).Elem()
	anyType = reflect.TypeOf((*Any)(nil)).Elem()
	marshalerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	unmarshalerType = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
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
	cfg      *frozenConfig
	cacheKey reflect.Type
}

func (encoder *placeholderEncoder) encode(ptr unsafe.Pointer, stream *Stream) {
	encoder.getRealEncoder().encode(ptr, stream)
}

func (encoder *placeholderEncoder) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (encoder *placeholderEncoder) isEmpty(ptr unsafe.Pointer) bool {
	return encoder.getRealEncoder().isEmpty(ptr)
}

func (encoder *placeholderEncoder) getRealEncoder() Encoder {
	for i := 0; i < 30; i++ {
		realDecoder := encoder.cfg.getEncoderFromCache(encoder.cacheKey)
		_, isPlaceholder := realDecoder.(*placeholderEncoder)
		if isPlaceholder {
			time.Sleep(time.Second)
		} else {
			return realDecoder
		}
	}
	panic(fmt.Sprintf("real encoder not found for cache key: %v", encoder.cacheKey))
}

type placeholderDecoder struct {
	cfg      *frozenConfig
	cacheKey reflect.Type
}

func (decoder *placeholderDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	for i := 0; i < 30; i++ {
		realDecoder := decoder.cfg.getDecoderFromCache(decoder.cacheKey)
		_, isPlaceholder := realDecoder.(*placeholderDecoder)
		if isPlaceholder {
			time.Sleep(time.Second)
		} else {
			realDecoder.decode(ptr, iter)
			return
		}
	}
	panic(fmt.Sprintf("real decoder not found for cache key: %v", decoder.cacheKey))
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
	cachedDecoder := iter.cfg.getDecoderFromCache(cacheKey)
	if cachedDecoder == nil {
		decoder, err := decoderOfType(iter.cfg, cacheKey)
		if err != nil {
			iter.Error = err
			return
		}
		cachedDecoder = decoder
		iter.cfg.addDecoderToCache(cacheKey, decoder)
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
	cachedEncoder := stream.cfg.getEncoderFromCache(cacheKey)
	if cachedEncoder == nil {
		encoder, err := encoderOfType(stream.cfg, cacheKey)
		if err != nil {
			stream.Error = err
			return
		}
		cachedEncoder = encoder
		stream.cfg.addEncoderToCache(cacheKey, encoder)
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

func decoderOfType(cfg *frozenConfig, typ reflect.Type) (Decoder, error) {
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
	cachedDecoder := cfg.getDecoderFromCache(cacheKey)
	if cachedDecoder != nil {
		return cachedDecoder, nil
	}
	placeholder := &placeholderDecoder{cfg: cfg, cacheKey: cacheKey}
	cfg.addDecoderToCache(cacheKey, placeholder)
	newDecoder, err := createDecoderOfType(cfg, typ)
	cfg.addDecoderToCache(cacheKey, newDecoder)
	return newDecoder, err
}

func createDecoderOfType(cfg *frozenConfig, typ reflect.Type) (Decoder, error) {
	if typ.String() == "[]uint8" {
		return &base64Codec{}, nil
	}
	if typ.AssignableTo(jsonRawMessageType) {
		return &jsonRawMessageCodec{}, nil
	}
	if typ.AssignableTo(jsoniterRawMessageType) {
		return &jsoniterRawMessageCodec{}, nil
	}
	if typ.AssignableTo(jsonNumberType) {
		return &jsonNumberCodec{}, nil
	}
	if typ.ConvertibleTo(unmarshalerType) {
		templateInterface := reflect.New(typ).Elem().Interface()
		var decoder Decoder = &unmarshalerDecoder{extractInterface(templateInterface)}
		if typ.Kind() != reflect.Struct {
			decoder = &optionalDecoder{typ, decoder}
		}
		return decoder, nil
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
		return prefix(fmt.Sprintf("[%s]", typ.String())).addToDecoder(decoderOfStruct(cfg, typ))
	case reflect.Array:
		return prefix("[array]").addToDecoder(decoderOfArray(cfg, typ))
	case reflect.Slice:
		return prefix("[slice]").addToDecoder(decoderOfSlice(cfg, typ))
	case reflect.Map:
		return prefix("[map]").addToDecoder(decoderOfMap(cfg, typ))
	case reflect.Ptr:
		return prefix("[optional]").addToDecoder(decoderOfOptional(cfg, typ))
	default:
		return nil, fmt.Errorf("unsupported type: %v", typ)
	}
}

func encoderOfType(cfg *frozenConfig, typ reflect.Type) (Encoder, error) {
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
	cachedEncoder := cfg.getEncoderFromCache(cacheKey)
	if cachedEncoder != nil {
		return cachedEncoder, nil
	}
	placeholder := &placeholderEncoder{cfg: cfg, cacheKey: cacheKey}
	cfg.addEncoderToCache(cacheKey, placeholder)
	newEncoder, err := createEncoderOfType(cfg, typ)
	cfg.addEncoderToCache(cacheKey, newEncoder)
	return newEncoder, err
}

func createEncoderOfType(cfg *frozenConfig, typ reflect.Type) (Encoder, error) {
	if typ.String() == "[]uint8" {
		return &base64Codec{}, nil
	}
	if typ.AssignableTo(jsonRawMessageType) {
		return &jsonRawMessageCodec{}, nil
	}
	if typ.AssignableTo(jsoniterRawMessageType) {
		return &jsoniterRawMessageCodec{}, nil
	}
	if typ.AssignableTo(jsonNumberType) {
		return &jsonNumberCodec{}, nil
	}
	if typ.ConvertibleTo(marshalerType) {
		templateInterface := reflect.New(typ).Elem().Interface()
		var encoder Encoder = &marshalerEncoder{extractInterface(templateInterface)}
		if typ.Kind() != reflect.Struct {
			encoder = &optionalEncoder{encoder}
		}
		return encoder, nil
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
		return prefix(fmt.Sprintf("[%s]", typ.String())).addToEncoder(encoderOfStruct(cfg, typ))
	case reflect.Array:
		return prefix("[array]").addToEncoder(encoderOfArray(cfg, typ))
	case reflect.Slice:
		return prefix("[slice]").addToEncoder(encoderOfSlice(cfg, typ))
	case reflect.Map:
		return prefix("[map]").addToEncoder(encoderOfMap(cfg, typ))
	case reflect.Ptr:
		return prefix("[optional]").addToEncoder(encoderOfOptional(cfg, typ))
	default:
		return nil, fmt.Errorf("unsupported type: %v", typ)
	}
}

func decoderOfOptional(cfg *frozenConfig, typ reflect.Type) (Decoder, error) {
	elemType := typ.Elem()
	decoder, err := decoderOfType(cfg, elemType)
	if err != nil {
		return nil, err
	}
	return &optionalDecoder{elemType, decoder}, nil
}

func encoderOfOptional(cfg *frozenConfig, typ reflect.Type) (Encoder, error) {
	elemType := typ.Elem()
	elemEncoder, err := encoderOfType(cfg, elemType)
	if err != nil {
		return nil, err
	}
	encoder := &optionalEncoder{elemEncoder}
	if elemType.Kind() == reflect.Map {
		encoder = &optionalEncoder{encoder}
	}
	return encoder, nil
}

func decoderOfMap(cfg *frozenConfig, typ reflect.Type) (Decoder, error) {
	decoder, err := decoderOfType(cfg, typ.Elem())
	if err != nil {
		return nil, err
	}
	mapInterface := reflect.New(typ).Interface()
	return &mapDecoder{typ, typ.Key(), typ.Elem(), decoder, extractInterface(mapInterface)}, nil
}

func extractInterface(val interface{}) emptyInterface {
	return *((*emptyInterface)(unsafe.Pointer(&val)))
}

func encoderOfMap(cfg *frozenConfig, typ reflect.Type) (Encoder, error) {
	elemType := typ.Elem()
	encoder, err := encoderOfType(cfg, elemType)
	if err != nil {
		return nil, err
	}
	mapInterface := reflect.New(typ).Elem().Interface()
	if cfg.sortMapKeys {
		return &sortKeysMapEncoder{typ, elemType, encoder, *((*emptyInterface)(unsafe.Pointer(&mapInterface)))}, nil
	} else {
		return &mapEncoder{typ, elemType, encoder, *((*emptyInterface)(unsafe.Pointer(&mapInterface)))}, nil
	}
}
