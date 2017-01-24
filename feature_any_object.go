package jsoniter

import (
	"unsafe"
)

type objectLazyAny struct {
	baseAny
	buf       []byte
	iter      *Iterator
	err       error
	cache     map[string]Any
	remaining []byte
}

func (any *objectLazyAny) parse() *Iterator {
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
	iter := any.parse()
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
				return v
			}
		} else {
			any.remaining = nil
			return nil
		}
	}
	for iter.nextToken() == ',' {
		k := string(iter.readObjectFieldAsBytes())
		v := iter.readAny(iter)
		any.cache[k] = v
		if target == k {
			any.remaining = iter.buf[iter.head:]
			return v
		}
	}
	any.remaining = nil
	return nil
}

func (any *objectLazyAny) fillCache() {
	if any.remaining == nil {
		return
	}
	if any.cache == nil {
		any.cache = map[string]Any{}
	}
	iter := any.parse()
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
			return
		}
	}
	for iter.nextToken() == ',' {
		k := string(iter.readObjectFieldAsBytes())
		v := iter.readAny(iter)
		any.cache[k] = v
	}
	any.remaining = nil
	return
}

func (any *objectLazyAny) LastError() error {
	return nil
}

func (any *objectLazyAny) ToBool() bool {
	return false
}

func (any *objectLazyAny) ToInt() int {
	return 0
}

func (any *objectLazyAny) ToInt32() int32 {
	return 0
}

func (any *objectLazyAny) ToInt64() int64 {
	return 0
}

func (any *objectLazyAny) ToFloat32() float32 {
	return 0
}

func (any *objectLazyAny) ToFloat64() float64 {
	return 0
}

func (any *objectLazyAny) ToString() string {
	if len(any.remaining) == len(any.buf) {
		// nothing has been parsed yet
		return *(*string)(unsafe.Pointer(&any.buf))
	} else {
		// TODO: serialize the cache
		return ""
	}
}

func (any *objectLazyAny) Get(path ...interface{}) Any {
	key := path[0].(string)
	return any.fillCacheUntil(key)
}

func (any *objectLazyAny) Keys() []string {
	any.fillCache()
	keys := make([]string, 0, len(any.cache))
	for key := range any.cache {
		keys = append(keys, key)
	}
	return keys
}

func (any *objectLazyAny) IterateObject() (func() (string, Any, bool), bool) {
	if any.cache == nil {
		any.cache = map[string]Any{}
	}
	remaining := any.remaining
	if len(remaining) == len(any.buf) {
		iter := any.parse()
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
				return key, value, true
			} else {
				remaining = nil
				any.remaining = nil
				return key, value, false
			}
		}
	}, true
}

