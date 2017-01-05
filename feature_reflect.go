package jsoniter

import (
	"reflect"
	"errors"
	"fmt"
	"unsafe"
	"sync/atomic"
	"strings"
	"io"
	"strconv"
)

type Decoder interface {
	decode(ptr unsafe.Pointer, iter *Iterator)
}

type stringDecoder struct {
}

func (decoder *stringDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*string)(ptr)) = iter.ReadString()
}

type intDecoder struct {
}

func (decoder *intDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int)(ptr)) = iter.ReadInt()
}

type int8Decoder struct {
}

func (decoder *int8Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int8)(ptr)) = iter.ReadInt8()
}

type int16Decoder struct {
}

func (decoder *int16Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int16)(ptr)) = iter.ReadInt16()
}

type int32Decoder struct {
}

func (decoder *int32Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int32)(ptr)) = iter.ReadInt32()
}

type int64Decoder struct {
}

func (decoder *int64Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*int64)(ptr)) = iter.ReadInt64()
}

type uintDecoder struct {
}

func (decoder *uintDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint)(ptr)) = iter.ReadUint()
}

type uint8Decoder struct {
}

func (decoder *uint8Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint8)(ptr)) = iter.ReadUint8()
}

type uint16Decoder struct {
}

func (decoder *uint16Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint16)(ptr)) = iter.ReadUint16()
}

type uint32Decoder struct {
}

func (decoder *uint32Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint32)(ptr)) = iter.ReadUint32()
}

type uint64Decoder struct {
}

func (decoder *uint64Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*uint64)(ptr)) = iter.ReadUint64()
}

type float32Decoder struct {
}

func (decoder *float32Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*float32)(ptr)) = iter.ReadFloat32()
}

type float64Decoder struct {
}

func (decoder *float64Decoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*float64)(ptr)) = iter.ReadFloat64()
}

type boolDecoder struct {
}

func (decoder *boolDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*bool)(ptr)) = iter.ReadBool()
}

type interfaceDecoder struct {
}

func (decoder *interfaceDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*interface{})(ptr)) = iter.ReadAny().Get()
}


type anyDecoder struct {
}

func (decoder *anyDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	*((*Any)(ptr)) = *iter.ReadAny()
}

type stringNumberDecoder struct {
	elemDecoder Decoder
}

func (decoder *stringNumberDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	c := iter.readByte()
	if c != '"' {
		iter.ReportError("stringNumberDecoder", `expect "`)
		return
	}
	decoder.elemDecoder.decode(ptr, iter)
	if iter.Error != nil {
		return
	}
	c = iter.readByte()
	if c != '"' {
		iter.ReportError("stringNumberDecoder", `expect "`)
		return
	}
}

type optionalDecoder struct {
	valueType    reflect.Type
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

type generalStructDecoder struct {
	type_  reflect.Type
	fields map[string]*structFieldDecoder
}

func (decoder *generalStructDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
		fieldDecoder := decoder.fields[field]
		if fieldDecoder == nil {
			iter.Skip()
		} else {
			fieldDecoder.decode(ptr, iter)
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.type_, iter.Error.Error())
	}
}

type skipDecoder struct {
	type_ reflect.Type
}

func (decoder *skipDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	iter.Skip()
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.type_, iter.Error.Error())
	}
}

type oneFieldStructDecoder struct {
	type_        reflect.Type
	fieldName    string
	fieldDecoder *structFieldDecoder
}

func (decoder *oneFieldStructDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
		if field == decoder.fieldName {
			decoder.fieldDecoder.decode(ptr, iter)
		} else {
			iter.Skip()
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.type_, iter.Error.Error())
	}
}

type twoFieldsStructDecoder struct {
	type_         reflect.Type
	fieldName1    string
	fieldDecoder1 *structFieldDecoder
	fieldName2    string
	fieldDecoder2 *structFieldDecoder
}

func (decoder *twoFieldsStructDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
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
		iter.Error = fmt.Errorf("%v: %s", decoder.type_, iter.Error.Error())
	}
}

type threeFieldsStructDecoder struct {
	type_         reflect.Type
	fieldName1    string
	fieldDecoder1 *structFieldDecoder
	fieldName2    string
	fieldDecoder2 *structFieldDecoder
	fieldName3    string
	fieldDecoder3 *structFieldDecoder
}

func (decoder *threeFieldsStructDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
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
		iter.Error = fmt.Errorf("%v: %s", decoder.type_, iter.Error.Error())
	}
}

type fourFieldsStructDecoder struct {
	type_         reflect.Type
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
	for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
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
		iter.Error = fmt.Errorf("%v: %s", decoder.type_, iter.Error.Error())
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
	mapType  reflect.Type
	elemType  reflect.Type
	elemDecoder Decoder
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
type ExtensionFunc func(type_ reflect.Type, field *reflect.StructField) ([]string, DecoderFunc)

type funcDecoder struct {
	func_ DecoderFunc
}

func (decoder *funcDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	decoder.func_(ptr, iter)
}

func RegisterTypeDecoder(type_ string, func_ DecoderFunc) {
	typeDecoders[type_] = &funcDecoder{func_}
}

func RegisterFieldDecoder(type_ string, field string, func_ DecoderFunc) {
	fieldDecoders[fmt.Sprintf("%s/%s", type_, field)] = &funcDecoder{func_}
}

func RegisterExtension(extension ExtensionFunc) {
	extensions = append(extensions, extension)
}

func ClearDecoders() {
	typeDecoders = map[string]Decoder{}
	fieldDecoders = map[string]Decoder{}
}

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  *struct{}
	word unsafe.Pointer
}

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
		for (iter.ReadArray()) {
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
		iter.ReportError("ReadAny", fmt.Sprintf("unexpected value type: %v", valueType))
		return MakeAny(nil)
	}
}



func (iter *Iterator) readNumber() (ret *Any) {
	strBuf := [8]byte{}
	str := strBuf[0:0]
	hasMore := true
	foundFloat := false
	foundNegative := false
	for(hasMore) {
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
				hasMore = false
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

func (iter *Iterator) Read(obj interface{}) {
	type_ := reflect.TypeOf(obj)
	cacheKey := type_.Elem()
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
	typeName := type_.String()
	if typeName == "jsoniter.Any" {
		return &anyDecoder{}, nil
	}

	for _, extension := range extensions {
		alternativeFieldNames, func_ := extension(type_, nil)
		if alternativeFieldNames != nil {
			return nil, fmt.Errorf("%v should not return alternative field names when only type is being passed", extension)
		}
		typeDecoders[typeName] = &funcDecoder{func_}
	}
	typeDecoder := typeDecoders[typeName]
	if typeDecoder != nil {
		return typeDecoder, nil
	}
	switch type_.Kind() {
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
		return decoderOfStruct(type_)
	case reflect.Slice:
		return prefix("[slice]").addTo(decoderOfSlice(type_))
	case reflect.Map:
		return prefix("[map]").addTo(decoderOfMap(type_))
	case reflect.Ptr:
		return prefix("[optional]").addTo(decoderOfOptional(type_.Elem()))
	default:
		return nil, fmt.Errorf("unsupported type: %v", type_)
	}
}

func decoderOfOptional(type_ reflect.Type) (Decoder, error) {
	decoder, err := decoderOfPtr(type_)
	if err != nil {
		return nil, err
	}
	return &optionalDecoder{type_, decoder}, nil
}

func decoderOfStruct(type_ reflect.Type) (Decoder, error) {
	fields := map[string]*structFieldDecoder{}
	for i := 0; i < type_.NumField(); i++ {
		field := type_.Field(i)
		fieldDecoderKey := fmt.Sprintf("%s/%s", type_.String(), field.Name)
		var fieldNames []string
		for _, extension := range extensions {
			alternativeFieldNames, func_ := extension(type_, &field)
			if alternativeFieldNames != nil {
				fieldNames = alternativeFieldNames
			}
			if func_ != nil {
				fieldDecoders[fieldDecoderKey] = &funcDecoder{func_}
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
		return &skipDecoder{type_}, nil
	case 1:
		for fieldName, fieldDecoder := range fields {
			return &oneFieldStructDecoder{type_, fieldName, fieldDecoder}, nil
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
		return &twoFieldsStructDecoder{type_, fieldName1, fieldDecoder1, fieldName2, fieldDecoder2}, nil
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
		return &threeFieldsStructDecoder{type_,
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
		return &fourFieldsStructDecoder{type_,
			fieldName1, fieldDecoder1, fieldName2, fieldDecoder2, fieldName3, fieldDecoder3,
			fieldName4, fieldDecoder4}, nil
	}
	return &generalStructDecoder{type_, fields}, nil
}

func decoderOfSlice(type_ reflect.Type) (Decoder, error) {
	decoder, err := decoderOfPtr(type_.Elem())
	if err != nil {
		return nil, err
	}
	return &sliceDecoder{type_, type_.Elem(), decoder}, nil
}

func decoderOfMap(type_ reflect.Type) (Decoder, error) {
	decoder, err := decoderOfPtr(type_.Elem())
	if err != nil {
		return nil, err
	}
	mapInterface := reflect.New(type_).Interface()
	return &mapDecoder{type_, type_.Elem(), decoder, *((*emptyInterface)(unsafe.Pointer(&mapInterface)))}, nil
}
