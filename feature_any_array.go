package jsoniter

import (
	"fmt"
	"reflect"
	"unsafe"
)

type arrayLazyAny struct {
	baseAny
	cfg       *frozenConfig
	buf       []byte
	err       error
}

func (any *arrayLazyAny) ValueType() ValueType {
	return Array
}

func (any *arrayLazyAny) LastError() error {
	return any.err
}

func (any *arrayLazyAny) ToBool() bool {
	iter := any.cfg.BorrowIterator(any.buf)
	defer any.cfg.ReturnIterator(iter)
	return iter.ReadArray()
}

func (any *arrayLazyAny) ToInt() int {
	if any.ToBool() {
		return 1
	} else {
		return 0
	}
}

func (any *arrayLazyAny) ToInt32() int32 {
	if any.ToBool() {
		return 1
	} else {
		return 0
	}
}

func (any *arrayLazyAny) ToInt64() int64 {
	if any.ToBool() {
		return 1
	} else {
		return 0
	}
}

func (any *arrayLazyAny) ToUint() uint {
	if any.ToBool() {
		return 1
	} else {
		return 0
	}
}

func (any *arrayLazyAny) ToUint32() uint32 {
	if any.ToBool() {
		return 1
	} else {
		return 0
	}
}

func (any *arrayLazyAny) ToUint64() uint64 {
	if any.ToBool() {
		return 1
	} else {
		return 0
	}
}

func (any *arrayLazyAny) ToFloat32() float32 {
	if any.ToBool() {
		return 1
	} else {
		return 0
	}
}

func (any *arrayLazyAny) ToFloat64() float64 {
	if any.ToBool() {
		return 1
	} else {
		return 0
	}
}

func (any *arrayLazyAny) ToString() string {
	return *(*string)(unsafe.Pointer(&any.buf))
}

func (any *arrayLazyAny) Get(path ...interface{}) Any {
	if len(path) == 0 {
		return any
	}
	switch firstPath := path[0].(type) {
	case int:
		iter := any.cfg.BorrowIterator(any.buf)
		defer any.cfg.ReturnIterator(iter)
		valueBytes := locateArrayElement(iter, firstPath)
		if valueBytes == nil {
			return newInvalidAny(path)
		} else {
			iter.ResetBytes(valueBytes)
			return locatePath(iter, path[1:])
		}
	case int32:
		if '*' == firstPath {
			iter := any.cfg.BorrowIterator(any.buf)
			defer any.cfg.ReturnIterator(iter)
			arr := make([]Any, 0)
			iter.ReadArrayCB(func(iter *Iterator) bool {
				found := iter.readAny().Get(path[1:]...)
				if found.ValueType() != Invalid {
					arr = append(arr, found)
				}
				return true
			})
			return wrapArray(arr)
		} else {
			return newInvalidAny(path)
		}
	default:
		return newInvalidAny(path)
	}
}

func (any *arrayLazyAny) Size() int {
	size := 0
	iter := any.cfg.BorrowIterator(any.buf)
	defer any.cfg.ReturnIterator(iter)
	iter.ReadArrayCB(func(iter *Iterator) bool {
		size++
		iter.Skip()
		return true
	})
	return size
}

func (any *arrayLazyAny) GetArray() []Any {
	elements := make([]Any, 0)
	iter := any.cfg.BorrowIterator(any.buf)
	defer any.cfg.ReturnIterator(iter)
	iter.ReadArrayCB(func(iter *Iterator) bool {
		elements = append(elements, iter.ReadAny())
		return true
	})
	return elements
}

func (any *arrayLazyAny) WriteTo(stream *Stream) {
	stream.Write(any.buf)
}

func (any *arrayLazyAny) GetInterface() interface{} {
	iter := any.cfg.BorrowIterator(any.buf)
	defer any.cfg.ReturnIterator(iter)
	return iter.Read()
}

type arrayAny struct {
	baseAny
	err   error
	cache []Any
	val   reflect.Value
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
