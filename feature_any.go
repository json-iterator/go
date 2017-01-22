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
		dotFound, lazyBuf := iter.skipNumber()
		if dotFound {
			return &floatLazyAny{lazyBuf, nil, nil}
		} else {
			return &intLazyAny{lazyBuf, nil, nil, 0}
		}
	}
	iter.reportError("ReadAny", fmt.Sprintf("unexpected value type: %v", valueType))
	return nil
}

func (iter *Iterator) skipNumber() (bool, []byte) {
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
				return dotFound, lazyBuf
			}
		}
		lazyBuf = append(lazyBuf, iter.buf[iter.head:iter.tail]...)
		if !iter.loadMore() {
			iter.head = iter.tail;
			return dotFound, lazyBuf
		}
	}
}
