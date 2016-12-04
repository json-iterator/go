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

type optionalDecoder struct {
	valueType reflect.Type
	valueDecoder Decoder
}

func (decoder *optionalDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	if iter.ReadNull() {
		*((*unsafe.Pointer)(ptr)) = nil
	} else {
		value := reflect.New(decoder.valueType)
		decoder.valueDecoder.decode(unsafe.Pointer(value.Pointer()), iter)
		*((*uintptr)(ptr)) = value.Pointer()
	}
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

type sliceDecoder struct {
	sliceType    reflect.Type
	elemType    reflect.Type
	elemDecoder Decoder
}

// sliceHeader is a safe version of SliceHeader used within this package.
type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

func (decoder *sliceDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	slice := (*sliceHeader)(ptr)
	slice.Len = 0
	for iter.ReadArray() {
		offset := uintptr(slice.Len) * decoder.elemType.Size()
		growOne(slice, decoder.sliceType, decoder.elemType)
		dataPtr := uintptr(slice.Data) + offset
		decoder.elemDecoder.decode(unsafe.Pointer(dataPtr), iter)
	}
}

// grow grows the slice s so that it can hold extra more values, allocating
// more capacity if needed. It also returns the old and new slice lengths.
func growOne(slice *sliceHeader, sliceType reflect.Type, elementType reflect.Type) {
	newLen := slice.Len + 1
	if newLen <= slice.Cap {
		slice.Len = newLen
		return
	}
	newCap := slice.Cap
	if newCap == 0 {
		newCap = 1
	} else {
		for newCap < newLen {
			if slice.Len < 1024 {
				newCap += newCap
			} else {
				newCap += newCap / 4
			}
		}
	}
	dst := unsafe.Pointer(reflect.MakeSlice(sliceType, newLen, newCap).Pointer())
	originalBytesCount := uintptr(slice.Len) * elementType.Size()
	srcPtr := (*[1<<30]byte)(slice.Data)
	dstPtr := (*[1<<30]byte)(dst)
	for i := uintptr(0); i < originalBytesCount; i++ {
		dstPtr[i] = srcPtr[i]
	}
	slice.Len = newLen
	slice.Cap = newCap
	slice.Data = dst
}

var DECODER_STRING *stringDecoder
var DECODERS unsafe.Pointer

func addDecoderToCache(cacheKey string, decoder Decoder) {
	retry := true
	for retry {
		ptr := atomic.LoadPointer(&DECODERS)
		cache := *(*map[string]Decoder)(ptr)
		copy := map[string]Decoder{}
		for k, v := range cache {
			copy[k] = v
		}
		copy[cacheKey] = decoder
		retry = !atomic.CompareAndSwapPointer(&DECODERS, ptr, unsafe.Pointer(&copy))
	}
}

func getDecoderFromCache(cacheKey string) Decoder {
	ptr := atomic.LoadPointer(&DECODERS)
	cache := *(*map[string]Decoder)(ptr)
	return cache[cacheKey]
}

func init() {
	DECODER_STRING = &stringDecoder{}
	atomic.StorePointer(&DECODERS, unsafe.Pointer(&map[string]Decoder{}))
}

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  *struct{}
	word unsafe.Pointer
}

func (iter *Iterator) Read(obj interface{}) {
	type_ := reflect.TypeOf(obj)
	cacheKey := type_.String()
	cachedDecoder := getDecoderFromCache(cacheKey)
	if cachedDecoder == nil {
		decoder, err := decoderOfType(type_)
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
	case reflect.Slice:
		return decoderOfSlice(type_)
	case reflect.Ptr:
		return prefix("optional").addTo(decoderOfOptional(type_.Elem()))
	default:
		return nil, errors.New("expect string, struct, slice")
	}
}

func decoderOfOptional(type_ reflect.Type) (Decoder, error) {
	switch type_.Kind() {
	case reflect.String:
		return &optionalDecoder{type_, DECODER_STRING}, nil
	default:
		return nil, errors.New("expect string")
	}
}


func decoderOfStruct(type_ reflect.Type) (Decoder, error) {
	fields := map[string]Decoder{}
	for i := 0; i < type_.NumField(); i++ {
		field := type_.Field(i)
		decoder, err := decoderOfPtr(field.Type)
		if err != nil {
			return prefix(fmt.Sprintf("{%s}", field.Name)).addTo(decoder, err)
		}
		fields[field.Name] = &structFieldDecoder{field.Offset, decoder}
	}
	return &structDecoder{fields}, nil
}

func decoderOfSlice(type_ reflect.Type) (Decoder, error) {
	decoder, err := decoderOfPtr(type_.Elem())
	if err != nil {
		return prefix("[elem]").addTo(decoder, err)
	}
	return &sliceDecoder{type_, type_.Elem(), decoder}, nil
}
