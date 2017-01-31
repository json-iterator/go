package jsoniter

import (
	"unsafe"
	"fmt"
	"reflect"
)

type arrayLazyAny struct {
	baseAny
	buf       []byte
	iter      *Iterator
	err       error
	cache     []Any
	remaining []byte
}

func (any *arrayLazyAny) ValueType() ValueType {
	return Array
}

func (any *arrayLazyAny) Parse() *Iterator {
	iter := any.iter
	if iter == nil {
		iter = NewIterator()
		any.iter = iter
	}
	iter.ResetBytes(any.remaining)
	return iter
}

func (any *arrayLazyAny) fillCacheUntil(target int) Any {
	if any.remaining == nil {
		if target >= len(any.cache) {
			return nil
		}
		return any.cache[target]
	}
	if any.cache == nil {
		any.cache = make([]Any, 0, 8)
	}
	i := len(any.cache)
	if target < i {
		return any.cache[target]
	}
	iter := any.Parse()
	if (len(any.remaining) == len(any.buf)) {
		iter.head++
		c := iter.nextToken()
		if c != ']' {
			iter.unreadByte()
			element := iter.readAny(iter)
			any.cache = append(any.cache, element)
			if target == 0 {
				any.remaining = iter.buf[iter.head:]
				any.err = iter.Error
				return element
			}
			i = 1
		} else {
			any.remaining = nil
			any.err = iter.Error
			return nil
		}
	}
	for iter.nextToken() == ',' {
		element := iter.readAny(iter)
		any.cache = append(any.cache, element)
		if i == target {
			any.remaining = iter.buf[iter.head:]
			any.err = iter.Error
			return element
		}
		i++
	}
	any.remaining = nil
	any.err = iter.Error
	return nil
}

func (any *arrayLazyAny) fillCache() {
	if any.remaining == nil {
		return
	}
	if any.cache == nil {
		any.cache = make([]Any, 0, 8)
	}
	iter := any.Parse()
	if len(any.remaining) == len(any.buf) {
		iter.head++
		c := iter.nextToken()
		if c != ']' {
			iter.unreadByte()
			any.cache = append(any.cache, iter.readAny(iter))
		} else {
			any.remaining = nil
			any.err = iter.Error
			return
		}
	}
	for iter.nextToken() == ',' {
		any.cache = append(any.cache, iter.readAny(iter))
	}
	any.remaining = nil
	any.err = iter.Error
}

func (any *arrayLazyAny) LastError() error {
	return any.err
}

func (any *arrayLazyAny) ToBool() bool {
	if any.cache == nil {
		any.IterateArray() // trigger first element read
	}
	return len(any.cache) != 0
}

func (any *arrayLazyAny) ToInt() int {
	if any.cache == nil {
		any.IterateArray() // trigger first element read
	}
	if len(any.cache) == 0 {
		return 0
	}
	return 1
}

func (any *arrayLazyAny) ToInt32() int32 {
	if any.cache == nil {
		any.IterateArray() // trigger first element read
	}
	if len(any.cache) == 0 {
		return 0
	}
	return 1
}

func (any *arrayLazyAny) ToInt64() int64 {
	if any.cache == nil {
		any.IterateArray() // trigger first element read
	}
	if len(any.cache) == 0 {
		return 0
	}
	return 1
}

func (any *arrayLazyAny) ToUint() uint {
	if any.cache == nil {
		any.IterateArray() // trigger first element read
	}
	if len(any.cache) == 0 {
		return 0
	}
	return 1
}

func (any *arrayLazyAny) ToUint32() uint32 {
	if any.cache == nil {
		any.IterateArray() // trigger first element read
	}
	if len(any.cache) == 0 {
		return 0
	}
	return 1
}

func (any *arrayLazyAny) ToUint64() uint64 {
	if any.cache == nil {
		any.IterateArray() // trigger first element read
	}
	if len(any.cache) == 0 {
		return 0
	}
	return 1
}

func (any *arrayLazyAny) ToFloat32() float32 {
	if any.cache == nil {
		any.IterateArray() // trigger first element read
	}
	if len(any.cache) == 0 {
		return 0
	}
	return 1
}

func (any *arrayLazyAny) ToFloat64() float64 {
	if any.cache == nil {
		any.IterateArray() // trigger first element read
	}
	if len(any.cache) == 0 {
		return 0
	}
	return 1
}

func (any *arrayLazyAny) ToString() string {
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

func (any *arrayLazyAny) Get(path ...interface{}) Any {
	if len(path) == 0 {
		return any
	}
	var element Any
	switch firstPath := path[0].(type) {
	case int:
		element = any.fillCacheUntil(firstPath)
		if element == nil {
			element = &invalidAny{baseAny{}, fmt.Errorf("%v not found in %v", firstPath, any.cache)}
		}
	case int32:
		if '*' == firstPath {
			any.fillCache()
			arr := make([]Any, 0, len(any.cache))
			for _, element := range any.cache {
				found := element.Get(path[1:]...)
				if found.ValueType() != Invalid {
					arr = append(arr, found)
				}
			}
			return wrapArray(arr)
		} else {
			element = &invalidAny{baseAny{}, fmt.Errorf("%v not found in %v", path[0], any.cache)}
		}
	default:
		element = &invalidAny{baseAny{}, fmt.Errorf("%v not found in %v", path[0], any.cache)}
	}
	if len(path) == 1 {
		return element
	} else {
		return element.Get(path[1:]...)
	}
}

func (any *arrayLazyAny) Size() int {
	any.fillCache()
	return len(any.cache)
}

func (any *arrayLazyAny) IterateArray() (func() (Any, bool), bool) {
	if any.cache == nil {
		any.cache = make([]Any, 0, 8)
	}
	remaining := any.remaining
	if len(remaining) == len(any.buf) {
		iter := any.Parse()
		iter.head++
		c := iter.nextToken()
		if c != ']' {
			iter.unreadByte()
			v := iter.readAny(iter)
			any.cache = append(any.cache, v)
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
	arr := any.cache
	nextValue := arr[0]
	i := 1
	return func() (Any, bool) {
		value := nextValue
		if i < len(arr) {
			// read from cache
			nextValue = arr[i]
			i++
			return value, true
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
				nextValue = iter.readAny(iter)
				any.cache = append(any.cache, nextValue)
				remaining = iter.buf[iter.head:]
				any.remaining = remaining
				any.err = iter.Error
				return value, true
			} else {
				remaining = nil
				any.remaining = nil
				any.err = iter.Error
				return value, false
			}
		}
	}, true
}

func (any *arrayLazyAny) GetArray() []Any {
	any.fillCache()
	return any.cache
}

func (any *arrayLazyAny) SetArray(newList []Any) bool {
	any.fillCache()
	any.cache = newList
	return true
}

func (any *arrayLazyAny) WriteTo(stream *Stream) {
	if len(any.remaining) == len(any.buf) {
		// nothing has been parsed yet
		stream.Write(any.buf)
	} else {
		any.fillCache()
		stream.WriteVal(any.cache)
	}
}

func (any *arrayLazyAny) GetInterface() interface{} {
	any.fillCache()
	return any.cache
}

type arrayAny struct {
	baseAny
	err      error
	cache    []Any
	val      reflect.Value
}

func wrapArray(val interface{}) *arrayAny {
	return &arrayAny{baseAny{}, nil, nil, reflect.ValueOf(val)}
}

func (any *arrayAny) ValueType() ValueType {
	return Array
}

func (any *arrayAny) Parse() *Iterator {
	return nil
}

func (any *arrayAny) LastError() error {
	return any.err
}

func (any *arrayAny) ToBool() bool {
	return any.val.Len() != 0
}

func (any *arrayAny) ToInt() int {
	if any.val.Len() == 0 {
		return 0
	}
	return 1
}

func (any *arrayAny) ToInt32() int32 {
	if any.val.Len() == 0 {
		return 0
	}
	return 1
}

func (any *arrayAny) ToInt64() int64 {
	if any.val.Len() == 0 {
		return 0
	}
	return 1
}

func (any *arrayAny) ToUint() uint {
	if any.val.Len() == 0 {
		return 0
	}
	return 1
}

func (any *arrayAny) ToUint32() uint32 {
	if any.val.Len() == 0 {
		return 0
	}
	return 1
}

func (any *arrayAny) ToUint64() uint64 {
	if any.val.Len() == 0 {
		return 0
	}
	return 1
}

func (any *arrayAny) ToFloat32() float32 {
	if any.val.Len() == 0 {
		return 0
	}
	return 1
}

func (any *arrayAny) ToFloat64() float64 {
	if any.val.Len() == 0 {
		return 0
	}
	return 1
}

func (any *arrayAny) ToString() string {
	if len(any.cache) == 0 {
		// nothing has been parsed yet
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

func (any *arrayAny) fillCacheUntil(idx int) Any {
	if idx < len(any.cache) {
		return any.cache[idx]
	} else {
		for i := len(any.cache); i < any.val.Len(); i++ {
			element := Wrap(any.val.Index(i).Interface())
			any.cache = append(any.cache, element)
			if idx == i {
				return element
			}
		}
		return nil
	}
}

func (any *arrayAny) fillCache() {
	any.cache = make([]Any, any.val.Len())
	for i := 0; i < any.val.Len(); i++ {
		any.cache[i] = Wrap(any.val.Index(i).Interface())
	}
}

func (any *arrayAny) Get(path ...interface{}) Any {
	if len(path) == 0 {
		return any
	}
	var element Any
	switch firstPath := path[0].(type) {
	case int:
		element = any.fillCacheUntil(firstPath)
		if element == nil {
			element = &invalidAny{baseAny{}, fmt.Errorf("%v not found in %v", firstPath, any.cache)}
		}
	case int32:
		if '*' == firstPath {
			any.fillCache()
			mappedAll := make([]Any, 0, len(any.cache))
			for _, element := range any.cache {
				mapped := element.Get(path[1:]...)
				if mapped.ValueType() != Invalid {
					mappedAll = append(mappedAll, mapped)
				}
			}
			return wrapArray(mappedAll)
		} else {
			element = &invalidAny{baseAny{}, fmt.Errorf("%v not found in %v", path[0], any.cache)}
		}
	default:
		element = &invalidAny{baseAny{}, fmt.Errorf("%v not found in %v", path[0], any.cache)}
	}
	if len(path) == 1 {
		return element
	} else {
		return element.Get(path[1:]...)
	}
}

func (any *arrayAny) Size() int {
	any.fillCache()
	return len(any.cache)
}

func (any *arrayAny) IterateArray() (func() (Any, bool), bool) {
	if any.val.Len() == 0 {
		return nil, false
	}
	i := 0
	return func() (Any, bool) {
		if i == any.val.Len() {
			return nil, false
		}
		if i == len(any.cache) {
			any.cache = append(any.cache, Wrap(any.val.Index(i).Interface()))
		}
		val := any.cache[i]
		i++
		return val, i != any.val.Len()
	}, true
}

func (any *arrayAny) GetArray() []Any {
	any.fillCache()
	return any.cache
}

func (any *arrayAny) SetArray(newList []Any) bool {
	any.fillCache()
	any.cache = newList
	return true
}

func (any *arrayAny) WriteTo(stream *Stream) {
	if len(any.cache) == 0 {
		// nothing has been parsed yet
		stream.WriteVal(any.val)
	} else {
		any.fillCache()
		stream.WriteVal(any.cache)
	}
}

func (any *arrayAny) GetInterface() interface{} {
	any.fillCache()
	return any.cache
}