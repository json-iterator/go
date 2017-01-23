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
	c := iter.nextToken()
	switch c {
	case '"':
		return iter.readStringAny()
	case 'n':
		iter.skipFixedBytes(3) // null
		return &nilAny{}
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		iter.unreadByte()
		return iter.readNumberAny()
	case 't':
		iter.skipFixedBytes(3) // true
		return &trueAny{}
	case 'f':
		iter.skipFixedBytes(4) // false
		return &falseAny{}
	}
	iter.reportError("ReadAny", fmt.Sprintf("unexpected character: %v", c))
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
