package jsoniter

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"sync/atomic"
	"unsafe"
)

type Config struct {
	IndentionStep                 int
	MarshalFloatWith6Digits       bool
	SupportUnexportedStructFields bool
	EscapeHtml                    bool
	SortMapKeys                   bool
	UseNumber                     bool
}

type frozenConfig struct {
	configBeforeFrozen Config
	sortMapKeys        bool
	indentionStep      int
	decoderCache       unsafe.Pointer
	encoderCache       unsafe.Pointer
	extensions         []ExtensionFunc
	streamPool         chan *Stream
	iteratorPool       chan *Iterator
}

type Api interface {
	MarshalToString(v interface{}) (string, error)
	Marshal(v interface{}) ([]byte, error)
	UnmarshalFromString(str string, v interface{}) error
	Unmarshal(data []byte, v interface{}) error
	Get(data []byte, path ...interface{}) Any
	NewEncoder(writer io.Writer) *AdaptedEncoder
	NewDecoder(reader io.Reader) *AdaptedDecoder
}

var ConfigDefault = Config{
	EscapeHtml: true,
}.Froze()

// Trying to be 100% compatible with standard library behavior
var ConfigCompatibleWithStandardLibrary = Config{
	EscapeHtml:  true,
	SortMapKeys: true,
}.Froze()

var ConfigFastest = Config{
	MarshalFloatWith6Digits: true,
}.Froze()

func (cfg Config) Froze() *frozenConfig {
	frozenConfig := &frozenConfig{
		sortMapKeys:   cfg.SortMapKeys,
		indentionStep: cfg.IndentionStep,
		streamPool:    make(chan *Stream, 16),
		iteratorPool:  make(chan *Iterator, 16),
	}
	atomic.StorePointer(&frozenConfig.decoderCache, unsafe.Pointer(&map[string]Decoder{}))
	atomic.StorePointer(&frozenConfig.encoderCache, unsafe.Pointer(&map[string]Encoder{}))
	if cfg.MarshalFloatWith6Digits {
		frozenConfig.marshalFloatWith6Digits()
	}
	if cfg.SupportUnexportedStructFields {
		frozenConfig.supportUnexportedStructFields()
	}
	if cfg.EscapeHtml {
		frozenConfig.escapeHtml()
	}
	if cfg.UseNumber {
		frozenConfig.useNumber()
	}
	frozenConfig.configBeforeFrozen = cfg
	return frozenConfig
}

func (cfg *frozenConfig) useNumber() {
	cfg.addDecoderToCache(reflect.TypeOf((*interface{})(nil)).Elem(), &funcDecoder{func(ptr unsafe.Pointer, iter *Iterator) {
		if iter.WhatIsNext() == Number {
			*((*interface{})(ptr)) = json.Number(iter.readNumberAsString())
		} else {
			*((*interface{})(ptr)) = iter.Read()
		}
	}})
}

// RegisterExtension can register a custom extension
func (cfg *frozenConfig) registerExtension(extension ExtensionFunc) {
	cfg.extensions = append(cfg.extensions, extension)
}

func (cfg *frozenConfig) supportUnexportedStructFields() {
	cfg.registerExtension(func(type_ reflect.Type, field *reflect.StructField) ([]string, EncoderFunc, DecoderFunc) {
		return []string{field.Name}, nil, nil
	})
}

// EnableLossyFloatMarshalling keeps 10**(-6) precision
// for float variables for better performance.
func (cfg *frozenConfig) marshalFloatWith6Digits() {
	// for better performance
	cfg.addEncoderToCache(reflect.TypeOf((*float32)(nil)).Elem(), &funcEncoder{func(ptr unsafe.Pointer, stream *Stream) {
		val := *((*float32)(ptr))
		stream.WriteFloat32Lossy(val)
	}})
	cfg.addEncoderToCache(reflect.TypeOf((*float64)(nil)).Elem(), &funcEncoder{func(ptr unsafe.Pointer, stream *Stream) {
		val := *((*float64)(ptr))
		stream.WriteFloat64Lossy(val)
	}})
}

type htmlEscapedStringEncoder struct {
}

func (encoder *htmlEscapedStringEncoder) encode(ptr unsafe.Pointer, stream *Stream) {
	str := *((*string)(ptr))
	stream.WriteStringWithHtmlEscaped(str)
}

func (encoder *htmlEscapedStringEncoder) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (encoder *htmlEscapedStringEncoder) isEmpty(ptr unsafe.Pointer) bool {
	return *((*string)(ptr)) == ""
}

func (cfg *frozenConfig) escapeHtml() {
	// for better performance
	cfg.addEncoderToCache(reflect.TypeOf((*string)(nil)).Elem(), &htmlEscapedStringEncoder{})
}

func (cfg *frozenConfig) addDecoderToCache(cacheKey reflect.Type, decoder Decoder) {
	done := false
	for !done {
		ptr := atomic.LoadPointer(&cfg.decoderCache)
		cache := *(*map[reflect.Type]Decoder)(ptr)
		copied := map[reflect.Type]Decoder{}
		for k, v := range cache {
			copied[k] = v
		}
		copied[cacheKey] = decoder
		done = atomic.CompareAndSwapPointer(&cfg.decoderCache, ptr, unsafe.Pointer(&copied))
	}
}

func (cfg *frozenConfig) addEncoderToCache(cacheKey reflect.Type, encoder Encoder) {
	done := false
	for !done {
		ptr := atomic.LoadPointer(&cfg.encoderCache)
		cache := *(*map[reflect.Type]Encoder)(ptr)
		copied := map[reflect.Type]Encoder{}
		for k, v := range cache {
			copied[k] = v
		}
		copied[cacheKey] = encoder
		done = atomic.CompareAndSwapPointer(&cfg.encoderCache, ptr, unsafe.Pointer(&copied))
	}
}

func (cfg *frozenConfig) getDecoderFromCache(cacheKey reflect.Type) Decoder {
	ptr := atomic.LoadPointer(&cfg.decoderCache)
	cache := *(*map[reflect.Type]Decoder)(ptr)
	return cache[cacheKey]
}

func (cfg *frozenConfig) getEncoderFromCache(cacheKey reflect.Type) Encoder {
	ptr := atomic.LoadPointer(&cfg.encoderCache)
	cache := *(*map[reflect.Type]Encoder)(ptr)
	return cache[cacheKey]
}

// cleanDecoders cleans decoders registered or cached
func (cfg *frozenConfig) cleanDecoders() {
	typeDecoders = map[string]Decoder{}
	fieldDecoders = map[string]Decoder{}
	atomic.StorePointer(&cfg.decoderCache, unsafe.Pointer(&map[string]Decoder{}))
}

// cleanEncoders cleans encoders registered or cached
func (cfg *frozenConfig) cleanEncoders() {
	typeEncoders = map[string]Encoder{}
	fieldEncoders = map[string]Encoder{}
	atomic.StorePointer(&cfg.encoderCache, unsafe.Pointer(&map[string]Encoder{}))
}

func (cfg *frozenConfig) MarshalToString(v interface{}) (string, error) {
	stream := cfg.BorrowStream(nil)
	defer cfg.ReturnStream(stream)
	stream.WriteVal(v)
	if stream.Error != nil {
		return "", stream.Error
	}
	return string(stream.Buffer()), nil
}

func (cfg *frozenConfig) Marshal(v interface{}) ([]byte, error) {
	stream := cfg.BorrowStream(nil)
	defer cfg.ReturnStream(stream)
	stream.WriteVal(v)
	if stream.Error != nil {
		return nil, stream.Error
	}
	result := stream.Buffer()
	copied := make([]byte, len(result))
	copy(copied, result)
	return copied, nil
}

func (cfg *frozenConfig) UnmarshalFromString(str string, v interface{}) error {
	data := []byte(str)
	data = data[:lastNotSpacePos(data)]
	iter := cfg.BorrowIterator(data)
	defer cfg.ReturnIterator(iter)
	iter.ReadVal(v)
	if iter.head == iter.tail {
		iter.loadMore()
	}
	if iter.Error == io.EOF {
		return nil
	}
	if iter.Error == nil {
		iter.reportError("UnmarshalFromString", "there are bytes left after unmarshal")
	}
	return iter.Error
}

func (cfg *frozenConfig) Get(data []byte, path ...interface{}) Any {
	iter := cfg.BorrowIterator(data)
	defer cfg.ReturnIterator(iter)
	return locatePath(iter, path)
}

func (cfg *frozenConfig) Unmarshal(data []byte, v interface{}) error {
	data = data[:lastNotSpacePos(data)]
	iter := cfg.BorrowIterator(data)
	defer cfg.ReturnIterator(iter)
	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Ptr {
		// return non-pointer error
		return errors.New("the second param must be ptr type")
	}
	iter.ReadVal(v)
	if iter.head == iter.tail {
		iter.loadMore()
	}
	if iter.Error == io.EOF {
		return nil
	}
	if iter.Error == nil {
		iter.reportError("Unmarshal", "there are bytes left after unmarshal")
	}
	return iter.Error
}

func (cfg *frozenConfig) NewEncoder(writer io.Writer) *AdaptedEncoder {
	stream := NewStream(cfg, writer, 512)
	return &AdaptedEncoder{stream}
}

func (cfg *frozenConfig) NewDecoder(reader io.Reader) *AdaptedDecoder {
	iter := Parse(cfg, reader, 512)
	return &AdaptedDecoder{iter}
}
