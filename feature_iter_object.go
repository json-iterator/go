package jsoniter

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
		iter.reportError("ReadObject", `expect { or , or } or n`)
		return
	}
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
			c = iter.nextToken()
			for c == ',' {
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
