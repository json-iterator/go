package jsoniter

import (
	"unsafe"
	"fmt"
	"reflect"
)

type objectLazyAny struct {
	baseAny
	buf       []byte
	iter      *Iterator
	err       error
	cache     map[string]Any
	remaining []byte
}

func (any *objectLazyAny) ValueType() ValueType {
	return Object
}

func (any *objectLazyAny) Parse() *Iterator {
	iter := any.iter
	if iter == nil {
		iter = NewIterator()
		any.iter = iter
	}
	iter.ResetBytes(any.remaining)
	return iter
}

func (any *objectLazyAny) fillCacheUntil(target string) Any {
	if any.remaining == nil {
		return any.cache[target]
	}
	if any.cache == nil {
		any.cache = map[string]Any{}
	}
	val := any.cache[target]
	if val != nil {
		return val
	}
	iter := any.Parse()
	if len(any.remaining) == len(any.buf) {
		iter.head++
		c := iter.nextToken()
		if c != '}' {
			iter.unreadByte()
			k := string(iter.readObjectFieldAsBytes())
			v := iter.readAny(iter)
			any.cache[k] = v
			if target == k {
				any.remaining = iter.buf[iter.head:]
				any.err = iter.Error
				return v
			}
		} else {
			any.remaining = nil
			any.err = iter.Error
			return nil
		}
	}
	for iter.nextToken() == ',' {
		k := string(iter.readObjectFieldAsBytes())
		v := iter.readAny(iter)
		any.cache[k] = v
		if target == k {
			any.remaining = iter.buf[iter.head:]
			any.err = iter.Error
			return v
		}
	}
	any.remaining = nil
	any.err = iter.Error
	return nil
}

func (any *objectLazyAny) fillCache() {
	if any.remaining == nil {
		return
	}
	if any.cache == nil {
		any.cache = map[string]Any{}
	}
	iter := any.Parse()
	if len(any.remaining) == len(any.buf) {
		iter.head++
		c := iter.nextToken()
		if c != '}' {
			iter.unreadByte()
			k := string(iter.readObjectFieldAsBytes())
			v := iter.readAny(iter)
			any.cache[k] = v
		} else {
			any.remaining = nil
			any.err = iter.Error
			return
		}
	}
	for iter.nextToken() == ',' {
		k := string(iter.readObjectFieldAsBytes())
		v := iter.readAny(iter)
		any.cache[k] = v
	}
	any.remaining = nil
	any.err = iter.Error
	return
}

func (any *objectLazyAny) LastError() error {
	return any.err
}

func (any *objectLazyAny) ToBool() bool {
	if any.cache == nil {
		any.IterateObject() // trigger first value read
	}
	return len(any.cache) != 0
}

func (any *objectLazyAny) ToInt() int {
	if any.cache == nil {
		any.IterateObject() // trigger first value read
	}
	if len(any.cache) == 0 {
		return 0
	}
	return 1
}

func (any *objectLazyAny) ToInt32() int32 {
	if any.cache == nil {
		any.IterateObject() // trigger first value read
	}
	if len(any.cache) == 0 {
		return 0
	}
	return 1
}

func (any *objectLazyAny) ToInt64() int64 {
	if any.cache == nil {
		any.IterateObject() // trigger first value read
	}
	if len(any.cache) == 0 {
		return 0
	}
	return 1
}

func (any *objectLazyAny) ToUint() uint {
	if any.cache == nil {
		any.IterateObject() // trigger first value read
	}
	if len(any.cache) == 0 {
		return 0
	}
	return 1
}

func (any *objectLazyAny) ToUint32() uint32 {
	if any.cache == nil {
		any.IterateObject() // trigger first value read
	}
	if len(any.cache) == 0 {
		return 0
	}
	return 1
}

func (any *objectLazyAny) ToUint64() uint64 {
	if any.cache == nil {
		any.IterateObject() // trigger first value read
	}
	if len(any.cache) == 0 {
		return 0
	}
	return 1
}

func (any *objectLazyAny) ToFloat32() float32 {
	if any.cache == nil {
		any.IterateObject() // trigger first value read
	}
	if len(any.cache) == 0 {
		return 0
	}
	return 1
}

func (any *objectLazyAny) ToFloat64() float64 {
	if any.cache == nil {
		any.IterateObject() // trigger first value read
	}
	if len(any.cache) == 0 {
		return 0
	}
	return 1
}

func (any *objectLazyAny) ToString() string {
	if len(any.remaining) == len(any.buf) {
		// nothing has been parsed yet
		return *(*string)(unsafe.Pointer(&any.buf))
	} else {
		any.fillCache()
		str, err := MarshalToString(any.cache)
		any.err = err
		return str
	}
}

func (any *objectLazyAny) Get(path ...interface{}) Any {
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

func (any *objectLazyAny) Keys() []string {
	any.fillCache()
	keys := make([]string, 0, len(any.cache))
	for key := range any.cache {
		keys = append(keys, key)
	}
	return keys
}

func (any *objectLazyAny) Size() int {
	any.fillCache()
	return len(any.cache)
}

func (any *objectLazyAny) IterateObject() (func() (string, Any, bool), bool) {
	if any.cache == nil {
		any.cache = map[string]Any{}
	}
	remaining := any.remaining
	if len(remaining) == len(any.buf) {
		iter := any.Parse()
		iter.head++
		c := iter.nextToken()
		if c != '}' {
			iter.unreadByte()
			k := string(iter.readObjectFieldAsBytes())
			v := iter.readAny(iter)
			any.cache[k] = v
			remaining = iter.buf[iter.head:]
			any.remaining = remaining
		} else {
			remaining = nil
			any.remaining = nil
			any.err = iter.Error
			return nil, false
		}
	}
	if len(any.cache) == 0 {
		return nil, false
	}
	keys := make([]string, 0, len(any.cache))
	values := make([]Any, 0, len(any.cache))
	for key, value := range any.cache {
		keys = append(keys, key)
		values = append(values, value)
	}
	nextKey := keys[0]
	nextValue := values[0]
	i := 1
	return func() (string, Any, bool) {
		key := nextKey
		value := nextValue
		if i < len(keys) {
			// read from cache
			nextKey = keys[i]
			nextValue = values[i]
			i++
			return key, value, true
		} else {
			// read from buffer
			iter := any.iter
			if iter == nil {
				iter = NewIterator()
				any.iter = iter
			}
			iter.ResetBytes(remaining)
			c := iter.nextToken()
			if c == ',' {
				nextKey = string(iter.readObjectFieldAsBytes())
				nextValue = iter.readAny(iter)
				any.cache[nextKey] = nextValue
				remaining = iter.buf[iter.head:]
				any.remaining = remaining
				any.err = iter.Error
				return key, value, true
			} else {
				remaining = nil
				any.remaining = nil
				any.err = iter.Error
				return key, value, false
			}
		}
	}, true
}

func (any *objectLazyAny) GetObject() map[string]Any {
	any.fillCache()
	return any.cache
}

func (any *objectLazyAny) SetObject(val map[string]Any) bool {
	any.fillCache()
	any.cache = val
	return true
}

func (any *objectLazyAny) WriteTo(stream *Stream) {
	if len(any.remaining) == len(any.buf) {
		// nothing has been parsed yet
		stream.Write(any.buf)
	} else {
		any.fillCache()
		stream.WriteVal(any.cache)
	}
}

func (any *objectLazyAny) GetInterface() interface{} {
	any.fillCache()
	return any.cache
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

func (any *objectAny) IterateObject() (func() (string, Any, bool), bool) {
	if any.cache == nil {
		any.cache = map[string]Any{}
	}
	if any.val.NumField() == 0 {
		return nil, false
	}
	cacheKeys := make([]string, len(any.cache))
	i := 0
	for key := range any.cache {
		cacheKeys[i] = key
		i++
	}
	i = 0
	return func() (string, Any, bool) {
		if i == any.val.NumField() {
			return "", nil, false
		}
		var fieldName string
		var fieldValueAsAny Any
		if i == len(cacheKeys) {
			fieldName = any.val.Type().Field(i).Name
			cacheKeys = append(cacheKeys, fieldName)
			fieldValue := any.val.Field(i)
			if fieldValue.CanInterface() {
				fieldValueAsAny = Wrap(fieldValue.Interface())
				any.cache[fieldName] = fieldValueAsAny
			} else {
				fieldValueAsAny = &invalidAny{baseAny{}, fmt.Errorf("%v not found in %v", fieldName, any.cache)}
				any.cache[fieldName] = fieldValueAsAny
			}
		} else {
			fieldName = cacheKeys[i]
			fieldValueAsAny = any.cache[fieldName]
		}
		i++
		return fieldName, fieldValueAsAny, i != any.val.NumField()
	}, true
}

func (any *objectAny) GetObject() map[string]Any {
	any.fillCache()
	return any.cache
}

func (any *objectAny) SetObject(val map[string]Any) bool {
	any.fillCache()
	any.cache = val
	return true
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

func (any *mapAny) IterateObject() (func() (string, Any, bool), bool) {
	any.fillCache()
	if len(any.cache) == 0 {
		return nil, false
	}
	keys := make([]string, len(any.cache))
	values := make([]Any, len(any.cache))
	i := 0
	for k, v := range any.cache {
		keys[i] = k
		values[i] = v
		i++
	}
	i = 0
	return func() (string, Any, bool) {
		if i == len(keys) {
			return "", nil, false
		}
		k := keys[i]
		v := values[i]
		i++
		return k, v, i != len(keys)
	}, true
}

func (any *mapAny) GetObject() map[string]Any {
	any.fillCache()
	return any.cache
}

func (any *mapAny) SetObject(val map[string]Any) bool {
	any.fillCache()
	any.cache = val
	return true
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
