package jsoniter

import (
	"reflect"
	"errors"
	"fmt"
	"unsafe"
	"sync/atomic"
)

type Decoder interface {
	decode(ptr unsafe.Pointer, iter *Iterator)
}

type stringDecoder struct {
}

func (decoder *stringDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*string)(ptr)) = iter.ReadString()
}

type structDecoder struct {
	fields map[string]Decoder
}

func (decoder *structDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
		fieldDecoder := decoder.fields[field]
		if fieldDecoder == nil {
			iter.Skip()
		} else {
			fieldDecoder.decode(ptr, iter)
		}
	}
}

type structFieldDecoder struct {
	offset       uintptr
	fieldDecoder Decoder
}

func (decoder *structFieldDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	fieldPtr := uintptr(ptr) + decoder.offset
	decoder.fieldDecoder.decode(unsafe.Pointer(fieldPtr), iter)
}

var DECODER_STRING *stringDecoder
var DECODERS_STRUCT unsafe.Pointer

func init() {
	DECODER_STRING = &stringDecoder{}
	atomic.StorePointer(&DECODERS_STRUCT, unsafe.Pointer(&map[string]*structDecoder{}))
}

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  *struct{}
	word unsafe.Pointer
}

func (iter *Iterator) Read(obj interface{}) {
	type_ := reflect.TypeOf(obj)
	decoder, err := decoderOfType(type_)
	if err != nil {
		iter.Error = err
		return
	}
	e := (*emptyInterface)(unsafe.Pointer(&obj))
	decoder.decode(e.word, iter)
}

type prefix string

func (p prefix) addTo(decoder Decoder, err error) (Decoder, error) {
	if err != nil {
		return nil, fmt.Errorf("%s: %s", p, err.Error())
	}
	return decoder, err
}

func decoderOfType(type_ reflect.Type) (Decoder, error) {
	switch type_.Kind() {
	case reflect.Ptr:
		return prefix("ptr").addTo(decoderOfPtr(type_.Elem()))
	default:
		return nil, errors.New("expect ptr")
	}
}

func decoderOfPtr(type_ reflect.Type) (Decoder, error) {
	switch type_.Kind() {
	case reflect.String:
		return DECODER_STRING, nil
	case reflect.Struct:
		return decoderOfStruct(type_)
	default:
		return nil, errors.New("expect string")
	}
}

func decoderOfStruct(type_ reflect.Type) (Decoder, error) {
	cacheKey := type_.String()
	cachedDecoder := getStructDecoderFromCache(cacheKey)
	if cachedDecoder == nil {
		fields := map[string]Decoder{}
		for i := 0; i < type_.NumField(); i++ {
			field := type_.Field(i)
			decoder, err := decoderOfPtr(field.Type)
			if err != nil {
				return prefix(fmt.Sprintf("[%s]", field.Name)).addTo(decoder, err)
			}
			fields[field.Name] = &structFieldDecoder{field.Offset, decoder}
		}
		cachedDecoder = &structDecoder{fields}
		addStructDecoderToCache(cacheKey, cachedDecoder)
	}
	return cachedDecoder, nil
}

func addStructDecoderToCache(cacheKey string, decoder *structDecoder) {
	retry := true
	for retry {
		ptr := atomic.LoadPointer(&DECODERS_STRUCT)
		cache := *(*map[string]*structDecoder)(ptr)
		copy := map[string]*structDecoder{}
		for k, v := range cache {
			copy[k] = v
		}
		copy[cacheKey] = decoder
		retry = !atomic.CompareAndSwapPointer(&DECODERS_STRUCT, ptr, unsafe.Pointer(&copy))
	}
}

func getStructDecoderFromCache(cacheKey string) *structDecoder {
	ptr := atomic.LoadPointer(&DECODERS_STRUCT)
	cache := *(*map[string]*structDecoder)(ptr)
	return cache[cacheKey]
}

