package jsoniter

import (
	"fmt"
	"reflect"
)

type Any interface {
	LastError() error
	ValueType() ValueType
	ToBool() bool
	ToInt() int
	ToInt32() int32
	ToInt64() int64
	ToUint() uint
	ToUint32() uint32
	ToUint64() uint64
	ToFloat32() float32
	ToFloat64() float64
	ToString() string
	Get(path ...interface{}) Any
	Size() int
	Keys() []string
	IterateObject() (func() (string, Any, bool), bool)
	IterateArray() (func() (Any, bool), bool)
	GetArray() []Any
	SetArray(newList []Any) bool
	GetObject() map[string]Any
	SetObject(map[string]Any) bool
	GetInterface() interface{}
	WriteTo(stream *Stream)
	Parse() *Iterator
}

type baseAny struct{}

func (any *baseAny) Get(path ...interface{}) Any {
	return &invalidAny{baseAny{}, fmt.Errorf("Get %v from simple value", path)}
}

func (any *baseAny) Size() int {
	return 0
}

func (any *baseAny) Keys() []string {
	return []string{}
}

func (any *baseAny) IterateObject() (func() (string, Any, bool), bool) {
	return nil, false
}

func (any *baseAny) IterateArray() (func() (Any, bool), bool) {
	return nil, false
}

func (any *baseAny) GetArray() []Any {
	return []Any{}
}

func (any *baseAny) SetArray(newList []Any) bool {
	return false
}

func (any *baseAny) GetObject() map[string]Any {
	return map[string]Any{}
}

func (any *baseAny) SetObject(map[string]Any) bool {
	return false
}

func WrapInt32(val int32) Any {
	return &int32Any{baseAny{}, val}
}

func WrapInt64(val int64) Any {
	return &int64Any{baseAny{}, val}
}

func WrapUint32(val uint32) Any {
	return &uint32Any{baseAny{}, val}
}

func WrapUint64(val uint64) Any {
	return &uint64Any{baseAny{}, val}
}

func WrapFloat64(val float64) Any {
	return &floatAny{baseAny{}, val}
}

func WrapString(val string) Any {
	return &stringAny{baseAny{}, nil, val}
}

func Wrap(val interface{}) Any {
	if val == nil {
		return &nilAny{}
	}
	type_ := reflect.TypeOf(val)
	switch type_.Kind() {
	case reflect.Slice:
		return wrapArray(val)
	case reflect.Struct:
		return wrapStruct(val)
	case reflect.Map:
		return wrapMap(val)
	case reflect.String:
		return WrapString(val.(string))
	case reflect.Int:
		return WrapInt64(int64(val.(int)))
	case reflect.Int8:
		return WrapInt32(int32(val.(int8)))
	case reflect.Int16:
		return WrapInt32(int32(val.(int16)))
	case reflect.Int32:
		return WrapInt32(val.(int32))
	case reflect.Int64:
		return WrapInt64(val.(int64))
	case reflect.Uint:
		return WrapUint64(uint64(val.(uint)))
	case reflect.Uint8:
		return WrapUint32(uint32(val.(uint8)))
	case reflect.Uint16:
		return WrapUint32(uint32(val.(uint16)))
	case reflect.Uint32:
		return WrapUint32(uint32(val.(uint32)))
	case reflect.Uint64:
		return WrapUint64(val.(uint64))
	case reflect.Float32:
		return WrapFloat64(float64(val.(float32)))
	case reflect.Float64:
		return WrapFloat64(val.(float64))
	case reflect.Bool:
		if val.(bool) == true {
			return &trueAny{}
		} else {
			return &falseAny{}
		}
	}
	return &invalidAny{baseAny{}, fmt.Errorf("unsupported type: %v", type_)}
}

func (iter *Iterator) ReadAny() Any {
	return iter.readAny(nil)
}

func (iter *Iterator) readAny(reusableIter *Iterator) Any {
	c := iter.nextToken()
	switch c {
	case '"':
		return iter.readStringAny(reusableIter)
	case 'n':
		iter.skipFixedBytes(3) // null
		return &nilAny{}
	case 't':
		iter.skipFixedBytes(3) // true
		return &trueAny{}
	case 'f':
		iter.skipFixedBytes(4) // false
		return &falseAny{}
	case '{':
		return iter.readObjectAny(reusableIter)
	case '[':
		return iter.readArrayAny(reusableIter)
	default:
		return iter.readNumberAny(reusableIter, c)
	}
}

func (iter *Iterator) readNumberAny(reusableIter *Iterator, firstByte byte) Any {
	dotFound := false
	lazyBuf := make([]byte, 1, 8)
	lazyBuf[0] = firstByte
	for {
		for i := iter.head; i < iter.tail; i++ {
			c := iter.buf[i]
			if c == '.' {
				dotFound = true
				continue
			}
			switch c {
			case ' ', '\n', '\r', '\t', ',', '}', ']':
				lazyBuf = append(lazyBuf, iter.buf[iter.head:i]...)
				iter.head = i
				if dotFound {
					return &float64LazyAny{baseAny{}, lazyBuf, reusableIter, nil, 0}
				} else {
					if firstByte == '-' {
						return &int64LazyAny{baseAny{}, lazyBuf, reusableIter, nil, 0}
					} else {
						return &uint64LazyAny{baseAny{}, lazyBuf, reusableIter, nil, 0}
					}
				}
			}
		}
		lazyBuf = append(lazyBuf, iter.buf[iter.head:iter.tail]...)
		if !iter.loadMore() {
			iter.head = iter.tail
			if dotFound {
				return &float64LazyAny{baseAny{}, lazyBuf, reusableIter, nil, 0}
			} else {
				if firstByte == '-' {
					return &int64LazyAny{baseAny{}, lazyBuf, reusableIter, nil, 0}
				} else {
					return &uint64LazyAny{baseAny{}, lazyBuf, reusableIter, nil, 0}
				}
			}
		}
	}
}

func (iter *Iterator) readStringAny(reusableIter *Iterator) Any {
	lazyBuf := make([]byte, 1, 8)
	lazyBuf[0] = '"'
	for {
		end, escaped := iter.findStringEnd()
		if end == -1 {
			lazyBuf = append(lazyBuf, iter.buf[iter.head:iter.tail]...)
			if !iter.loadMore() {
				iter.reportError("readStringAny", "incomplete string")
				return &invalidAny{}
			}
			if escaped {
				iter.head = 1 // skip the first char as last char read is \
			}
		} else {
			lazyBuf = append(lazyBuf, iter.buf[iter.head:end]...)
			iter.head = end
			return &stringLazyAny{baseAny{}, lazyBuf, reusableIter, nil, ""}
		}
	}
}

func (iter *Iterator) readObjectAny(reusableIter *Iterator) Any {
	level := 1
	lazyBuf := make([]byte, 1, 32)
	lazyBuf[0] = '{'
	for {
		start := iter.head
		for i := iter.head; i < iter.tail; i++ {
			switch iter.buf[i] {
			case '"': // If inside string, skip it
				iter.head = i + 1
				iter.skipString()
				i = iter.head - 1 // it will be i++ soon
			case '{': // If open symbol, increase level
				level++
			case '}': // If close symbol, increase level
				level--

				// If we have returned to the original level, we're done
				if level == 0 {
					iter.head = i + 1
					lazyBuf = append(lazyBuf, iter.buf[start:iter.head]...)
					return &objectLazyAny{baseAny{}, lazyBuf, reusableIter, nil, nil, lazyBuf}
				}
			}
		}
		lazyBuf = append(lazyBuf, iter.buf[iter.head:iter.tail]...)
		if !iter.loadMore() {
			iter.reportError("skipObject", "incomplete object")
			return &invalidAny{}
		}
	}
}

func (iter *Iterator) readArrayAny(reusableIter *Iterator) Any {
	level := 1
	lazyBuf := make([]byte, 1, 32)
	lazyBuf[0] = '['
	for {
		start := iter.head
		for i := iter.head; i < iter.tail; i++ {
			switch iter.buf[i] {
			case '"': // If inside string, skip it
				iter.head = i + 1
				iter.skipString()
				i = iter.head - 1 // it will be i++ soon
			case '[': // If open symbol, increase level
				level++
			case ']': // If close symbol, increase level
				level--

				// If we have returned to the original level, we're done
				if level == 0 {
					iter.head = i + 1
					lazyBuf = append(lazyBuf, iter.buf[start:iter.head]...)
					return &arrayLazyAny{baseAny{}, lazyBuf, reusableIter, nil, nil, lazyBuf}
				}
			}
		}
		lazyBuf = append(lazyBuf, iter.buf[iter.head:iter.tail]...)
		if !iter.loadMore() {
			iter.reportError("skipArray", "incomplete array")
			return &invalidAny{}
		}
	}
}
