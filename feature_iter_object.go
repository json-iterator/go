package jsoniter

import "fmt"

func (iter *Iterator) ReadObject() (ret string) {
	c := iter.nextToken()
	switch c {
	case 'n':
		iter.skipFixedBytes(3)
		return "" // null
	case '{':
		c = iter.nextToken()
		if c == '"' {
			iter.unreadByte()
			return string(iter.readObjectFieldAsBytes())
		}
		if c == '}' {
			return "" // end of object
		}
		iter.reportError("ReadObject", `expect " after {`)
		return
	case ',':
		return string(iter.readObjectFieldAsBytes())
	case '}':
		return "" // end of object
	default:
		iter.reportError("ReadObject", fmt.Sprintf(`expect { or , or } or n, but found %s`, string([]byte{c})))
		return
	}
}

func (iter *Iterator) readFieldHash() int32 {
	hash := int64(0x811c9dc5)
	c := iter.nextToken()
	if c == '"' {
		for {
			for i := iter.head; i < iter.tail; i++ {
				// require ascii string and no escape
				b := iter.buf[i]
				if b == '"' {
					iter.head = i+1
					c = iter.nextToken()
					if c != ':' {
						iter.reportError("readFieldHash", `expect :, but found ` + string([]byte{c}))
					}
					return int32(hash)
				}
				hash ^= int64(b)
				hash *= 0x1000193
			}
			if !iter.loadMore() {
				iter.reportError("readFieldHash", `incomplete field name`)
				return 0
			}
		}
	}
	iter.reportError("readFieldHash", `expect ", but found ` + string([]byte{c}))
	return 0
}

func calcHash(str string) int32 {
	hash := int64(0x811c9dc5)
	for _, b := range str {
		hash ^= int64(b)
		hash *= 0x1000193
	}
	return int32(hash)
}

func (iter *Iterator) ReadObjectCB(callback func(*Iterator, string) bool) bool {
	c := iter.nextToken()
	if c == '{' {
		c = iter.nextToken()
		if c == '"' {
			iter.unreadByte()
			field := string(iter.readObjectFieldAsBytes())
			if !callback(iter, field) {
				return false
			}
			for iter.nextToken() == ',' {
				field := string(iter.readObjectFieldAsBytes())
				if !callback(iter, field) {
					return false
				}
			}
			return true
		}
		if c == '}' {
			return true
		}
		iter.reportError("ReadObjectCB", `expect " after }`)
		return false
	}
	if c == 'n' {
		iter.skipFixedBytes(3)
		return true // null
	}
	iter.reportError("ReadObjectCB", `expect { or n`)
	return false
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

func (iter *Iterator) readObjectFieldAsBytes() (ret []byte) {
	str := iter.ReadStringAsSlice()
	if iter.skipWhitespacesWithoutLoadMore() {
		if ret == nil {
			ret = make([]byte, len(str))
			copy(ret, str)
		}
		if !iter.loadMore() {
			return
		}
	}
	if iter.buf[iter.head] != ':' {
		iter.reportError("readObjectFieldAsBytes", "expect : after object field")
		return
	}
	iter.head++
	if iter.skipWhitespacesWithoutLoadMore() {
		if ret == nil {
			ret = make([]byte, len(str))
			copy(ret, str)
		}
		if !iter.loadMore() {
			return
		}
	}
	if ret == nil {
		return str
	}
	return ret
}
