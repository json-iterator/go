package jsoniter

import "fmt"

type Any interface {
	LastError() error
	ToBool() bool
	ToInt() int
	ToInt32() int32
	ToInt64() int64
	ToFloat32() float32
	ToFloat64() float64
	ToString() string
}

func (iter *Iterator) ReadAny() Any {
	valueType := iter.WhatIsNext()
	switch valueType {
	case Nil:
		iter.skipFixedBytes(4)
		return &nilAny{}
	case Number:
		return iter.readNumberAny()
	case String:
		return iter.readStringAny()
	}
	iter.reportError("ReadAny", fmt.Sprintf("unexpected value type: %v", valueType))
	return &invalidAny{}
}

func (iter *Iterator) readNumberAny() Any {
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
					return &floatLazyAny{lazyBuf, nil, nil, 0}
				} else {
					return &intLazyAny{lazyBuf, nil, nil, 0}
				}
			}
		}
		lazyBuf = append(lazyBuf, iter.buf[iter.head:iter.tail]...)
		if !iter.loadMore() {
			iter.head = iter.tail
			if dotFound {
				return &floatLazyAny{lazyBuf, nil, nil, 0}
			} else {
				return &intLazyAny{lazyBuf, nil, nil, 0}
			}
		}
	}
}

func (iter *Iterator) readStringAny() Any {
	iter.head++
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
			return &stringLazyAny{lazyBuf, nil, nil, ""}
		}
	}
}
