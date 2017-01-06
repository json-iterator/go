package jsoniter

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
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

// Decoder works like a father class for sub-type decoders
type Decoder interface {
	decode(ptr unsafe.Pointer, iter *Iterator)
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

type generalStructDecoder struct {
	typ    reflect.Type
	fields map[string]*structFieldDecoder
}

func (decoder *generalStructDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	field := iter.readObjectField()
	fieldDecoder := decoder.fields[field]
	if fieldDecoder == nil {
		iter.Skip()
	} else {
		fieldDecoder.decode(ptr, iter)
	}
	for iter.nextToken() == ',' {
		field = iter.readObjectField()
		fieldDecoder = decoder.fields[field]
		if fieldDecoder == nil {
			iter.Skip()
		} else {
			fieldDecoder.decode(ptr, iter)
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type skipDecoder struct {
	typ reflect.Type
}

func (decoder *skipDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	iter.Skip()
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type oneFieldStructDecoder struct {
	typ          reflect.Type
	fieldName    string
	fieldDecoder *structFieldDecoder
}

func (decoder *oneFieldStructDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	field := iter.readObjectField()
	if field == decoder.fieldName {
		decoder.fieldDecoder.decode(ptr, iter)
	} else {
		iter.Skip()
	}
	for iter.nextToken() == ',' {
		field = iter.readObjectField()
		if field == decoder.fieldName {
			decoder.fieldDecoder.decode(ptr, iter)
		} else {
			iter.Skip()
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type twoFieldsStructDecoder struct {
	typ           reflect.Type
	fieldName1    string
	fieldDecoder1 *structFieldDecoder
	fieldName2    string
	fieldDecoder2 *structFieldDecoder
}

func (decoder *twoFieldsStructDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	field := iter.readObjectField()
	switch field {
	case decoder.fieldName1:
		decoder.fieldDecoder1.decode(ptr, iter)
	case decoder.fieldName2:
		decoder.fieldDecoder2.decode(ptr, iter)
	default:
		iter.Skip()
	}
	for iter.nextToken() == ',' {
		field = iter.readObjectField()
		switch field {
		case decoder.fieldName1:
			decoder.fieldDecoder1.decode(ptr, iter)
		case decoder.fieldName2:
			decoder.fieldDecoder2.decode(ptr, iter)
		default:
			iter.Skip()
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type threeFieldsStructDecoder struct {
	typ           reflect.Type
	fieldName1    string
	fieldDecoder1 *structFieldDecoder
	fieldName2    string
	fieldDecoder2 *structFieldDecoder
	fieldName3    string
	fieldDecoder3 *structFieldDecoder
}

func (decoder *threeFieldsStructDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	field := iter.readObjectField()
	switch field {
	case decoder.fieldName1:
		decoder.fieldDecoder1.decode(ptr, iter)
	case decoder.fieldName2:
		decoder.fieldDecoder2.decode(ptr, iter)
	case decoder.fieldName3:
		decoder.fieldDecoder3.decode(ptr, iter)
	default:
		iter.Skip()
	}
	for iter.nextToken() == ',' {
		field = iter.readObjectField()
		switch field {
		case decoder.fieldName1:
			decoder.fieldDecoder1.decode(ptr, iter)
		case decoder.fieldName2:
			decoder.fieldDecoder2.decode(ptr, iter)
		case decoder.fieldName3:
			decoder.fieldDecoder3.decode(ptr, iter)
		default:
			iter.Skip()
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type fourFieldsStructDecoder struct {
	typ           reflect.Type
	fieldName1    string
	fieldDecoder1 *structFieldDecoder
	fieldName2    string
	fieldDecoder2 *structFieldDecoder
	fieldName3    string
	fieldDecoder3 *structFieldDecoder
	fieldName4    string
	fieldDecoder4 *structFieldDecoder
}

func (decoder *fourFieldsStructDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	field := iter.readObjectField()
	switch field {
	case decoder.fieldName1:
		decoder.fieldDecoder1.decode(ptr, iter)
	case decoder.fieldName2:
		decoder.fieldDecoder2.decode(ptr, iter)
	case decoder.fieldName3:
		decoder.fieldDecoder3.decode(ptr, iter)
	case decoder.fieldName4:
		decoder.fieldDecoder4.decode(ptr, iter)
	default:
		iter.Skip()
	}
	for iter.nextToken() == ',' {
		field = iter.readObjectField()
		switch field {
		case decoder.fieldName1:
			decoder.fieldDecoder1.decode(ptr, iter)
		case decoder.fieldName2:
			decoder.fieldDecoder2.decode(ptr, iter)
		case decoder.fieldName3:
			decoder.fieldDecoder3.decode(ptr, iter)
		case decoder.fieldName4:
			decoder.fieldDecoder4.decode(ptr, iter)
		default:
			iter.Skip()
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type structFieldDecoder struct {
	field        *reflect.StructField
	fieldDecoder Decoder
}

func (decoder *structFieldDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	fieldPtr := uintptr(ptr) + decoder.field.Offset
	decoder.fieldDecoder.decode(unsafe.Pointer(fieldPtr), iter)
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%s: %s", decoder.field.Name, iter.Error.Error())
	}
}

type mapDecoder struct {
	mapType      reflect.Type
	elemType     reflect.Type
	elemDecoder  Decoder
	mapInterface emptyInterface
}

func (decoder *mapDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	// dark magic to cast unsafe.Pointer back to interface{} using reflect.Type
	mapInterface := decoder.mapInterface
	mapInterface.word = ptr
	realInterface := (*interface{})(unsafe.Pointer(&mapInterface))
	realVal := reflect.ValueOf(*realInterface).Elem()

	for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
		elem := reflect.New(decoder.elemType)
		decoder.elemDecoder.decode(unsafe.Pointer(elem.Pointer()), iter)
		// to put into map, we have to use reflection
		realVal.SetMapIndex(reflect.ValueOf(string([]byte(field))), elem.Elem())
	}
}

type sliceDecoder struct {
	sliceType   reflect.Type
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
	decoder.doDecode(ptr, iter)
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.sliceType, iter.Error.Error())
	}
}

func (decoder *sliceDecoder) doDecode(ptr unsafe.Pointer, iter *Iterator) {
	slice := (*sliceHeader)(ptr)
	reuseSlice(slice, decoder.sliceType, 4)
	if !iter.ReadArray() {
		return
	}
	offset := uintptr(0)
	decoder.elemDecoder.decode(unsafe.Pointer(uintptr(slice.Data) + offset), iter)
	if !iter.ReadArray() {
		slice.Len = 1
		return
	}
	offset += decoder.elemType.Size()
	decoder.elemDecoder.decode(unsafe.Pointer(uintptr(slice.Data) + offset), iter)
	if !iter.ReadArray() {
		slice.Len = 2
		return
	}
	offset += decoder.elemType.Size()
	decoder.elemDecoder.decode(unsafe.Pointer(uintptr(slice.Data) + offset), iter)
	if !iter.ReadArray() {
		slice.Len = 3
		return
	}
	offset += decoder.elemType.Size()
	decoder.elemDecoder.decode(unsafe.Pointer(uintptr(slice.Data) + offset), iter)
	slice.Len = 4
	for iter.ReadArray() {
		growOne(slice, decoder.sliceType, decoder.elemType)
		offset += decoder.elemType.Size()
		decoder.elemDecoder.decode(unsafe.Pointer(uintptr(slice.Data) + offset), iter)
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
	// copy old array into new array
	originalBytesCount := uintptr(slice.Len) * elementType.Size()
	srcPtr := (*[1 << 30]byte)(slice.Data)
	dstPtr := (*[1 << 30]byte)(dst)
	for i := uintptr(0); i < originalBytesCount; i++ {
		dstPtr[i] = srcPtr[i]
	}
	slice.Len = newLen
	slice.Cap = newCap
	slice.Data = dst
}

func reuseSlice(slice *sliceHeader, sliceType reflect.Type, expectedCap int) {
	if expectedCap <= slice.Cap {
		return
	}
	dst := unsafe.Pointer(reflect.MakeSlice(sliceType, 0, expectedCap).Pointer())
	slice.Cap = expectedCap
	slice.Data = dst
}

var DECODERS unsafe.Pointer

func addDecoderToCache(cacheKey reflect.Type, decoder Decoder) {
	retry := true
	for retry {
		ptr := atomic.LoadPointer(&DECODERS)
		cache := *(*map[reflect.Type]Decoder)(ptr)
		copy := map[reflect.Type]Decoder{}
		for k, v := range cache {
			copy[k] = v
		}
		copy[cacheKey] = decoder
		retry = !atomic.CompareAndSwapPointer(&DECODERS, ptr, unsafe.Pointer(&copy))
	}
}

func getDecoderFromCache(cacheKey reflect.Type) Decoder {
	ptr := atomic.LoadPointer(&DECODERS)
	cache := *(*map[reflect.Type]Decoder)(ptr)
	return cache[cacheKey]
}

var typeDecoders map[string]Decoder
var fieldDecoders map[string]Decoder
var extensions []ExtensionFunc

func init() {
	typeDecoders = map[string]Decoder{}
	fieldDecoders = map[string]Decoder{}
	extensions = []ExtensionFunc{}
	atomic.StorePointer(&DECODERS, unsafe.Pointer(&map[string]Decoder{}))
}

type DecoderFunc func(ptr unsafe.Pointer, iter *Iterator)
type ExtensionFunc func(typ reflect.Type, field *reflect.StructField) ([]string, DecoderFunc)

type funcDecoder struct {
	fun DecoderFunc
}

func (decoder *funcDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	decoder.fun(ptr, iter)
}

// RegisterTypeDecoder can register a type for json object
func RegisterTypeDecoder(typ string, fun DecoderFunc) {
	typeDecoders[typ] = &funcDecoder{fun}
}

// RegisterFieldDecoder can register a type for json field
func RegisterFieldDecoder(typ string, field string, fun DecoderFunc) {
	fieldDecoders[fmt.Sprintf("%s/%s", typ, field)] = &funcDecoder{fun}
}

// RegisterExtension can register a custom extension
func RegisterExtension(extension ExtensionFunc) {
	extensions = append(extensions, extension)
}

// CleanDecoders cleans decoders registered
func CleanDecoders() {
	typeDecoders = map[string]Decoder{}
	fieldDecoders = map[string]Decoder{}
}

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  *struct{}
	word unsafe.Pointer
}

// ReadAny converts a json object in a Iterator instance to Any
func (iter *Iterator) ReadAny() (ret *Any) {
	valueType := iter.WhatIsNext()
	switch valueType {
	case String:
		return MakeAny(iter.ReadString())
	case Number:
		return iter.readNumber()
	case Null:
		return MakeAny(nil)
	case Bool:
		return MakeAny(iter.ReadBool())
	case Array:
		val := []interface{}{}
		for iter.ReadArray() {
			element := iter.ReadAny()
			if iter.Error != nil {
				return
			}
			val = append(val, element.val)
		}
		return MakeAny(val)
	case Object:
		val := map[string]interface{}{}
		for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
			element := iter.ReadAny()
			if iter.Error != nil {
				return
			}
			val[string([]byte(field))] = element.val
		}
		return MakeAny(val)
	default:
		iter.reportError("ReadAny", fmt.Sprintf("unexpected value type: %v", valueType))
		return MakeAny(nil)
	}
}

func (iter *Iterator) readNumber() (ret *Any) {
	strBuf := [8]byte{}
	str := strBuf[0:0]
	hasMore := true
	foundFloat := false
	foundNegative := false
	for hasMore {
		for i := iter.head; i < iter.tail; i++ {
			c := iter.buf[i]
			switch c {
			case '-':
				foundNegative = true
				str = append(str, c)
				continue
			case '+', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				str = append(str, c)
				continue
			case '.', 'e', 'E':
				foundFloat = true
				str = append(str, c)
				continue
			default:
				iter.head = i
				hasMore = false
				break
			}
			if !hasMore {
				break
			}
		}
		if hasMore {
			if !iter.loadMore() {
				break
			}
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		return
	}
	number := *(*string)(unsafe.Pointer(&str))
	if foundFloat {
		val, err := strconv.ParseFloat(number, 64)
		if err != nil {
			iter.Error = err
			return
		}
		return MakeAny(val)
	}
	if foundNegative {
		val, err := strconv.ParseInt(number, 10, 64)
		if err != nil {
			iter.Error = err
			return
		}
		return MakeAny(val)
	}
	val, err := strconv.ParseUint(number, 10, 64)
	if err != nil {
		iter.Error = err
		return
	}
	return MakeAny(val)
}

// Read converts an Iterator instance into go interface, same as json.Unmarshal
func (iter *Iterator) Read(obj interface{}) {
	typ := reflect.TypeOf(obj)
	cacheKey := typ.Elem()
	cachedDecoder := getDecoderFromCache(cacheKey)
	if cachedDecoder == nil {
		decoder, err := decoderOfType(typ)
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

func decoderOfType(typ reflect.Type) (Decoder, error) {
	switch typ.Kind() {
	case reflect.Ptr:
		return prefix("ptr").addTo(decoderOfPtr(typ.Elem()))
	default:
		return nil, errors.New("expect ptr")
	}
}

func decoderOfPtr(typ reflect.Type) (Decoder, error) {
	typeName := typ.String()
	if typeName == "jsoniter.Any" {
		return &anyDecoder{}, nil
	}
	typeDecoder := typeDecoders[typeName]
	if typeDecoder != nil {
		return typeDecoder, nil
	}
	switch typ.Kind() {
	case reflect.String:
		return &stringDecoder{}, nil
	case reflect.Int:
		return &intDecoder{}, nil
	case reflect.Int8:
		return &int8Decoder{}, nil
	case reflect.Int16:
		return &int16Decoder{}, nil
	case reflect.Int32:
		return &int32Decoder{}, nil
	case reflect.Int64:
		return &int64Decoder{}, nil
	case reflect.Uint:
		return &uintDecoder{}, nil
	case reflect.Uint8:
		return &uint8Decoder{}, nil
	case reflect.Uint16:
		return &uint16Decoder{}, nil
	case reflect.Uint32:
		return &uint32Decoder{}, nil
	case reflect.Uint64:
		return &uint64Decoder{}, nil
	case reflect.Float32:
		return &float32Decoder{}, nil
	case reflect.Float64:
		return &float64Decoder{}, nil
	case reflect.Bool:
		return &boolDecoder{}, nil
	case reflect.Interface:
		return &interfaceDecoder{}, nil
	case reflect.Struct:
		return decoderOfStruct(typ)
	case reflect.Slice:
		return prefix("[slice]").addTo(decoderOfSlice(typ))
	case reflect.Map:
		return prefix("[map]").addTo(decoderOfMap(typ))
	case reflect.Ptr:
		return prefix("[optional]").addTo(decoderOfOptional(typ.Elem()))
	default:
		return nil, fmt.Errorf("unsupported type: %v", typ)
	}
}

func decoderOfOptional(typ reflect.Type) (Decoder, error) {
	decoder, err := decoderOfPtr(typ)
	if err != nil {
		return nil, err
	}
	return &optionalDecoder{typ, decoder}, nil
}

func decoderOfStruct(typ reflect.Type) (Decoder, error) {
	fields := map[string]*structFieldDecoder{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldDecoderKey := fmt.Sprintf("%s/%s", typ.String(), field.Name)
		var fieldNames []string
		for _, extension := range extensions {
			alternativeFieldNames, fun := extension(typ, &field)
			if alternativeFieldNames != nil {
				fieldNames = alternativeFieldNames
			}
			if fun != nil {
				fieldDecoders[fieldDecoderKey] = &funcDecoder{fun}
			}
		}
		decoder := fieldDecoders[fieldDecoderKey]
		tagParts := strings.Split(field.Tag.Get("json"), ",")
		// if fieldNames set by extension, use theirs, otherwise try tags
		if fieldNames == nil {
			/// tagParts[0] always present, even if no tags
			switch tagParts[0] {
			case "":
				fieldNames = []string{field.Name}
			case "-":
				fieldNames = []string{}
			default:
				fieldNames = []string{tagParts[0]}
			}
		}
		if decoder == nil {
			var err error
			decoder, err = decoderOfPtr(field.Type)
			if err != nil {
				return prefix(fmt.Sprintf("{%s}", field.Name)).addTo(decoder, err)
			}
		}
		if len(tagParts) > 1 && tagParts[1] == "string" {
			decoder = &stringNumberDecoder{decoder}
		}
		for _, fieldName := range fieldNames {
			fields[fieldName] = &structFieldDecoder{&field, decoder}
		}
	}
	switch len(fields) {
	case 0:
		return &skipDecoder{typ}, nil
	case 1:
		for fieldName, fieldDecoder := range fields {
			return &oneFieldStructDecoder{typ, fieldName, fieldDecoder}, nil
		}
	case 2:
		var fieldName1 string
		var fieldName2 string
		var fieldDecoder1 *structFieldDecoder
		var fieldDecoder2 *structFieldDecoder
		for fieldName, fieldDecoder := range fields {
			if fieldName1 == "" {
				fieldName1 = fieldName
				fieldDecoder1 = fieldDecoder
			} else {
				fieldName2 = fieldName
				fieldDecoder2 = fieldDecoder
			}
		}
		return &twoFieldsStructDecoder{typ, fieldName1, fieldDecoder1, fieldName2, fieldDecoder2}, nil
	case 3:
		var fieldName1 string
		var fieldName2 string
		var fieldName3 string
		var fieldDecoder1 *structFieldDecoder
		var fieldDecoder2 *structFieldDecoder
		var fieldDecoder3 *structFieldDecoder
		for fieldName, fieldDecoder := range fields {
			if fieldName1 == "" {
				fieldName1 = fieldName
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == "" {
				fieldName2 = fieldName
				fieldDecoder2 = fieldDecoder
			} else {
				fieldName3 = fieldName
				fieldDecoder3 = fieldDecoder
			}
		}
		return &threeFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1, fieldName2, fieldDecoder2, fieldName3, fieldDecoder3}, nil
	case 4:
		var fieldName1 string
		var fieldName2 string
		var fieldName3 string
		var fieldName4 string
		var fieldDecoder1 *structFieldDecoder
		var fieldDecoder2 *structFieldDecoder
		var fieldDecoder3 *structFieldDecoder
		var fieldDecoder4 *structFieldDecoder
		for fieldName, fieldDecoder := range fields {
			if fieldName1 == "" {
				fieldName1 = fieldName
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == "" {
				fieldName2 = fieldName
				fieldDecoder2 = fieldDecoder
			} else if fieldName3 == "" {
				fieldName3 = fieldName
				fieldDecoder3 = fieldDecoder
			} else {
				fieldName4 = fieldName
				fieldDecoder4 = fieldDecoder
			}
		}
		return &fourFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1, fieldName2, fieldDecoder2, fieldName3, fieldDecoder3,
			fieldName4, fieldDecoder4}, nil
	}
	return &generalStructDecoder{typ, fields}, nil
}

func decoderOfSlice(typ reflect.Type) (Decoder, error) {
	decoder, err := decoderOfPtr(typ.Elem())
	if err != nil {
		return nil, err
	}
	return &sliceDecoder{typ, typ.Elem(), decoder}, nil
}

func decoderOfMap(typ reflect.Type) (Decoder, error) {
	decoder, err := decoderOfPtr(typ.Elem())
	if err != nil {
		return nil, err
	}
	mapInterface := reflect.New(typ).Interface()
	return &mapDecoder{typ, typ.Elem(), decoder, *((*emptyInterface)(unsafe.Pointer(&mapInterface)))}, nil
}
