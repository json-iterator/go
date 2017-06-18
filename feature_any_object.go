package jsoniter

import (
	"fmt"
	"reflect"
	"unsafe"
)

type objectLazyAny struct {
	baseAny
	cfg *frozenConfig
	buf []byte
	err error
}

func (any *objectLazyAny) ValueType() ValueType {
	return Object
}

func (any *objectLazyAny) LastError() error {
	return any.err
}

func (any *objectLazyAny) ToBool() bool {
	iter := any.cfg.BorrowIterator(any.buf)
	defer any.cfg.ReturnIterator(iter)
	return iter.ReadObject() != ""
}

func (any *objectLazyAny) ToInt() int {
	if any.ToBool() {
		return 1
	} else {
		return 0
	}
}

func (any *objectLazyAny) ToInt32() int32 {
	if any.ToBool() {
		return 1
	} else {
		return 0
	}
}

func (any *objectLazyAny) ToInt64() int64 {
	if any.ToBool() {
		return 1
	} else {
		return 0
	}
}

func (any *objectLazyAny) ToUint() uint {
	if any.ToBool() {
		return 1
	} else {
		return 0
	}
}

func (any *objectLazyAny) ToUint32() uint32 {
	if any.ToBool() {
		return 1
	} else {
		return 0
	}
}

func (any *objectLazyAny) ToUint64() uint64 {
	if any.ToBool() {
		return 1
	} else {
		return 0
	}
}

func (any *objectLazyAny) ToFloat32() float32 {
	if any.ToBool() {
		return 1
	} else {
		return 0
	}
}

func (any *objectLazyAny) ToFloat64() float64 {
	if any.ToBool() {
		return 1
	} else {
		return 0
	}
}

func (any *objectLazyAny) ToString() string {
	return *(*string)(unsafe.Pointer(&any.buf))
}

func (any *objectLazyAny) Get(path ...interface{}) Any {
	if len(path) == 0 {
		return any
	}
	switch firstPath := path[0].(type) {
	case string:
		iter := any.cfg.BorrowIterator(any.buf)
		defer any.cfg.ReturnIterator(iter)
		valueBytes := locateObjectField(iter, firstPath)
		if valueBytes == nil {
			return newInvalidAny(path)
		} else {
			iter.ResetBytes(valueBytes)
			return locatePath(iter, path[1:])
		}
	case int32:
		if '*' == firstPath {
			mappedAll := map[string]Any{}
			iter := any.cfg.BorrowIterator(any.buf)
			defer any.cfg.ReturnIterator(iter)
			iter.ReadObjectCB(func(iter *Iterator, field string) bool {
				mapped := locatePath(iter, path[1:])
				if mapped.ValueType() != Invalid {
					mappedAll[field] = mapped
				}
				return true
			})
			return wrapMap(mappedAll)
		} else {
			return newInvalidAny(path)
		}
	default:
		return newInvalidAny(path)
	}
}

func (any *objectLazyAny) Keys() []string {
	keys := []string{}
	iter := any.cfg.BorrowIterator(any.buf)
	defer any.cfg.ReturnIterator(iter)
	iter.ReadObjectCB(func(iter *Iterator, field string) bool {
		iter.Skip()
		keys = append(keys, field)
		return true
	})
	return keys
}

func (any *objectLazyAny) Size() int {
	size := 0
	iter := any.cfg.BorrowIterator(any.buf)
	defer any.cfg.ReturnIterator(iter)
	iter.ReadObjectCB(func(iter *Iterator, field string) bool {
		iter.Skip()
		size ++
		return true
	})
	return size
}

func (any *objectLazyAny) GetObject() map[string]Any {
	asMap := map[string]Any{}
	iter := any.cfg.BorrowIterator(any.buf)
	defer any.cfg.ReturnIterator(iter)
	iter.ReadObjectCB(func(iter *Iterator, field string) bool {
		asMap[field] = iter.ReadAny()
		return true
	})
	return asMap
}

func (any *objectLazyAny) WriteTo(stream *Stream) {
	stream.Write(any.buf)
}

func (any *objectLazyAny) GetInterface() interface{} {
	iter := any.cfg.BorrowIterator(any.buf)
	defer any.cfg.ReturnIterator(iter)
	return iter.Read()
}

type objectAny struct {
	baseAny
	err   error
	cache map[string]Any
	val   reflect.Value
}

func wrapStruct(val interface{}) *objectAny {
	return &objectAny{baseAny{}, nil, nil, reflect.ValueOf(val)}
}

func (any *objectAny) ValueType() ValueType {
	return Object
}

func (any *objectAny) Parse() *Iterator {
	return nil
}

func (any *objectAny) fillCacheUntil(target string) Any {
	if any.cache == nil {
		any.cache = map[string]Any{}
	}
	element, found := any.cache[target]
	if found {
		return element
	}
	for i := len(any.cache); i < any.val.NumField(); i++ {
		field := any.val.Field(i)
		fieldName := any.val.Type().Field(i).Name
		var element Any
		if field.CanInterface() {
			element = Wrap(field.Interface())
		} else {
			element = &invalidAny{baseAny{}, fmt.Errorf("%v not found in %v", fieldName, any.cache)}
		}
		any.cache[fieldName] = element
		if fieldName == target {
			return element
		}
	}
	return nil
}

func (any *objectAny) fillCache() {
	if any.cache == nil {
		any.cache = map[string]Any{}
	}
	if len(any.cache) == any.val.NumField() {
		return
	}
	for i := 0; i < any.val.NumField(); i++ {
		field := any.val.Field(i)
		fieldName := any.val.Type().Field(i).Name
		var element Any
		if field.CanInterface() {
			element = Wrap(field.Interface())
		} else {
			element = &invalidAny{baseAny{}, fmt.Errorf("%v not found in %v", fieldName, any.cache)}
		}
		any.cache[fieldName] = element
	}
}

func (any *objectAny) LastError() error {
	return any.err
}

func (any *objectAny) ToBool() bool {
	return any.val.NumField() != 0
}

func (any *objectAny) ToInt() int {
	if any.val.NumField() == 0 {
		return 0
	}
	return 1
}

func (any *objectAny) ToInt32() int32 {
	if any.val.NumField() == 0 {
		return 0
	}
	return 1
}

func (any *objectAny) ToInt64() int64 {
	if any.val.NumField() == 0 {
		return 0
	}
	return 1
}

func (any *objectAny) ToUint() uint {
	if any.val.NumField() == 0 {
		return 0
	}
	return 1
}

func (any *objectAny) ToUint32() uint32 {
	if any.val.NumField() == 0 {
		return 0
	}
	return 1
}

func (any *objectAny) ToUint64() uint64 {
	if any.val.NumField() == 0 {
		return 0
	}
	return 1
}

func (any *objectAny) ToFloat32() float32 {
	if any.val.NumField() == 0 {
		return 0
	}
	return 1
}

func (any *objectAny) ToFloat64() float64 {
	if any.val.NumField() == 0 {
		return 0
	}
	return 1
}

func (any *objectAny) ToString() string {
	if len(any.cache) == 0 {
		str, err := MarshalToString(any.val.Interface())
		any.err = err
		return str
	} else {
		any.fillCache()
		str, err := MarshalToString(any.cache)
		any.err = err
		return str
	}
}

func (any *objectAny) Get(path ...interface{}) Any {
	if len(path) == 0 {
		return any
	}
	var element Any
	switch firstPath := path[0].(type) {
	case string:
		element = any.fillCacheUntil(firstPath)
		if element == nil {
			element = &invalidAny{baseAny{}, fmt.Errorf("%v not found in %v", firstPath, any.cache)}
		}
	case int32:
		if '*' == firstPath {
			any.fillCache()
			mappedAll := map[string]Any{}
			for key, value := range any.cache {
				mapped := value.Get(path[1:]...)
				if mapped.ValueType() != Invalid {
					mappedAll[key] = mapped
				}
			}
			return wrapMap(mappedAll)
		} else {
			element = &invalidAny{baseAny{}, fmt.Errorf("%v not found in %v", firstPath, any.cache)}
		}
	default:
		element = &invalidAny{baseAny{}, fmt.Errorf("%v not found in %v", firstPath, any.cache)}
	}
	if len(path) == 1 {
		return element
	} else {
		return element.Get(path[1:]...)
	}
}

func (any *objectAny) Keys() []string {
	any.fillCache()
	keys := make([]string, 0, len(any.cache))
	for key := range any.cache {
		keys = append(keys, key)
	}
	return keys
}

func (any *objectAny) Size() int {
	any.fillCache()
	return len(any.cache)
}

func (any *objectAny) GetObject() map[string]Any {
	any.fillCache()
	return any.cache
}

func (any *objectAny) WriteTo(stream *Stream) {
	if len(any.cache) == 0 {
		// nothing has been parsed yet
		stream.WriteVal(any.val)
	} else {
		any.fillCache()
		stream.WriteVal(any.cache)
	}
}

func (any *objectAny) GetInterface() interface{} {
	any.fillCache()
	return any.cache
}

type mapAny struct {
	baseAny
	err   error
	cache map[string]Any
	val   reflect.Value
}

func wrapMap(val interface{}) *mapAny {
	return &mapAny{baseAny{}, nil, nil, reflect.ValueOf(val)}
}

func (any *mapAny) ValueType() ValueType {
	return Object
}

func (any *mapAny) Parse() *Iterator {
	return nil
}

func (any *mapAny) fillCacheUntil(target string) Any {
	if any.cache == nil {
		any.cache = map[string]Any{}
	}
	element, found := any.cache[target]
	if found {
		return element
	}
	for _, key := range any.val.MapKeys() {
		keyAsStr := key.String()
		_, found := any.cache[keyAsStr]
		if found {
			continue
		}
		element := Wrap(any.val.MapIndex(key).Interface())
		any.cache[keyAsStr] = element
		if keyAsStr == target {
			return element
		}
	}
	return nil
}

func (any *mapAny) fillCache() {
	if any.cache == nil {
		any.cache = map[string]Any{}
	}
	if len(any.cache) == any.val.Len() {
		return
	}
	for _, key := range any.val.MapKeys() {
		keyAsStr := key.String()
		element := Wrap(any.val.MapIndex(key).Interface())
		any.cache[keyAsStr] = element
	}
}

func (any *mapAny) LastError() error {
	return any.err
}

func (any *mapAny) ToBool() bool {
	return any.val.Len() != 0
}

func (any *mapAny) ToInt() int {
	if any.val.Len() == 0 {
		return 0
	}
	return 1
}

func (any *mapAny) ToInt32() int32 {
	if any.val.Len() == 0 {
		return 0
	}
	return 1
}

func (any *mapAny) ToInt64() int64 {
	if any.val.Len() == 0 {
		return 0
	}
	return 1
}

func (any *mapAny) ToUint() uint {
	if any.val.Len() == 0 {
		return 0
	}
	return 1
}

func (any *mapAny) ToUint32() uint32 {
	if any.val.Len() == 0 {
		return 0
	}
	return 1
}

func (any *mapAny) ToUint64() uint64 {
	if any.val.Len() == 0 {
		return 0
	}
	return 1
}

func (any *mapAny) ToFloat32() float32 {
	if any.val.Len() == 0 {
		return 0
	}
	return 1
}

func (any *mapAny) ToFloat64() float64 {
	if any.val.Len() == 0 {
		return 0
	}
	return 1
}

func (any *mapAny) ToString() string {
	if len(any.cache) == 0 {
		str, err := MarshalToString(any.val.Interface())
		any.err = err
		return str
	} else {
		any.fillCache()
		str, err := MarshalToString(any.cache)
		any.err = err
		return str
	}
}

func (any *mapAny) Get(path ...interface{}) Any {
	if len(path) == 0 {
		return any
	}
	var element Any
	switch firstPath := path[0].(type) {
	case string:
		element = any.fillCacheUntil(firstPath)
		if element == nil {
			element = &invalidAny{baseAny{}, fmt.Errorf("%v not found in %v", firstPath, any.cache)}
		}
	case int32:
		if '*' == firstPath {
			any.fillCache()
			mappedAll := map[string]Any{}
			for key, value := range any.cache {
				mapped := value.Get(path[1:]...)
				if mapped.ValueType() != Invalid {
					mappedAll[key] = mapped
				}
			}
			return wrapMap(mappedAll)
		} else {
			element = &invalidAny{baseAny{}, fmt.Errorf("%v not found in %v", firstPath, any.cache)}
		}
	default:
		element = &invalidAny{baseAny{}, fmt.Errorf("%v not found in %v", firstPath, any.cache)}
	}
	if len(path) == 1 {
		return element
	} else {
		return element.Get(path[1:]...)
	}
}

func (any *mapAny) Keys() []string {
	any.fillCache()
	keys := make([]string, 0, len(any.cache))
	for key := range any.cache {
		keys = append(keys, key)
	}
	return keys
}

func (any *mapAny) Size() int {
	any.fillCache()
	return len(any.cache)
}

func (any *mapAny) GetObject() map[string]Any {
	any.fillCache()
	return any.cache
}

func (any *mapAny) WriteTo(stream *Stream) {
	if len(any.cache) == 0 {
		// nothing has been parsed yet
		stream.WriteVal(any.val)
	} else {
		any.fillCache()
		stream.WriteVal(any.cache)
	}
}

func (any *mapAny) GetInterface() interface{} {
	any.fillCache()
	return any.cache
}
