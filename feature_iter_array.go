package jsoniter

func (iter *Iterator) ReadArray() (ret bool) {
	c := iter.nextToken()
	switch c {
	case 'n':
		iter.skipFixedBytes(3)
		return false // null
	case '[':
		c = iter.nextToken()
		if c != ']' {
			iter.unreadByte()
			return true
		}
		return false
	case ']':
		return false
	case ',':
		return true
	default:
		iter.reportError("ReadArray", "expect [ or , or ] or n, but found: " + string([]byte{c}))
		return
	}
}


func (iter *Iterator) ReadArrayCB(callback func(*Iterator) bool) (ret bool) {
	c := iter.nextToken()
	if c == '[' {
		c = iter.nextToken()
		if c != ']' {
			iter.unreadByte()
			if !callback(iter) {
				return false
			}
			for iter.nextToken() == ',' {
				if !callback(iter) {
					return false
				}
			}
			return true
		}
		return true
	}
	if c == 'n' {
		iter.skipFixedBytes(3)
		return true // null
	}
	iter.reportError("ReadArrayCB", "expect [ or n, but found: " + string([]byte{c}))
	return false
}