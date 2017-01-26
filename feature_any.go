package jsoniter

import "fmt"

type Any interface {
	LastError() error
	ValueType() ValueType
	ToBool() bool
	ToInt() int
	ToInt32() int32
	ToInt64() int64
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

func WrapInt64(val int64) Any {
	return &intAny{baseAny{}, val}
}

func WrapFloat64(val float64) Any {
	return &floatAny{baseAny{}, val}
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
		iter.unreadByte()
		return iter.readNumberAny(reusableIter)
	}
}

func (iter *Iterator) readNumberAny(reusableIter *Iterator) Any {
	dotFound := false
	var lazyBuf []byte
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
					return &floatLazyAny{baseAny{}, lazyBuf, reusableIter, nil, 0}
				} else {
					return &intLazyAny{baseAny{}, lazyBuf, reusableIter, nil, 0}
				}
			}
		}
		lazyBuf = append(lazyBuf, iter.buf[iter.head:iter.tail]...)
		if !iter.loadMore() {
			iter.head = iter.tail
			if dotFound {
				return &floatLazyAny{baseAny{}, lazyBuf, reusableIter, nil, 0}
			} else {
				return &intLazyAny{baseAny{}, lazyBuf, reusableIter, nil, 0}
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
