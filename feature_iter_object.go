package jsoniter

import "unsafe"

// ReadObject is a implemented iterator for json
func (iter *Iterator) ReadObject() (ret string) {
	c := iter.nextToken()
	if iter.Error != nil {
		return
	}
	switch c {
	case 'n':
		iter.skipFixedBytes(3)
		if iter.Error != nil {
			return
		}
		return "" // null
	case '{':
		c = iter.nextToken()
		if iter.Error != nil {
			return
		}
		switch c {
		case '}':
			return "" // end of object
		case '"':
			iter.unreadByte()
			return iter.readObjectField()
		default:
			iter.reportError("ReadObject", `expect " after {`)
			return
		}
	case ',':
		return iter.readObjectField()
	case '}':
		return "" // end of object
	default:
		iter.reportError("ReadObject", `expect { or , or } or n`)
		return
	}
}

func (iter *Iterator) readObjectStart() bool {
	c := iter.nextToken()
	if c == '{' {
		c = iter.nextToken()
		if c == '}' {
			return false
		}
		iter.unreadByte()
		return true
	}
	iter.reportError("readObjectStart", "expect { ")
	return false
}

func (iter *Iterator) readObjectField() (ret string) {
	str := iter.ReadStringAsSlice()
	if iter.skipWhitespacesWithoutLoadMore() {
		if ret == "" {
			ret = string(str)
		}
		if !iter.loadMore() {
			return
		}
	}
	if iter.buf[iter.head] != ':' {
		iter.reportError("ReadObject", "expect : after object field")
		return
	}
	iter.head++
	if iter.skipWhitespacesWithoutLoadMore() {
		if ret == "" {
			ret = string(str)
		}
		if !iter.loadMore() {
			return
		}
	}
	if ret == "" {
		return *(*string)(unsafe.Pointer(&str))
	}
	return ret
}
