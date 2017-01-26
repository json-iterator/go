package jsoniter

import (
	"unsafe"
)

type arrayLazyAny struct {
	baseAny
	buf       []byte
	iter      *Iterator
	err       error
	cache     []Any
	remaining []byte
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
	if len(path) == 1 {
		idx := path[0].(int)
		return any.fillCacheUntil(idx)
	} else {
		idx := path[0].(int)
		return any.fillCacheUntil(idx).Get(path[1:]...)
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