package jsoniter


var intDigits []int8
const maxUint64 = (1<<64 - 1)
const cutoffUint64 = maxUint64/10 + 1
const maxUint32 = (1<<32 - 1)
const cutoffUint32 = maxUint32/10 + 1
const int32SafeToMultiply10 = uint32(int32(0x7fffffff)/10 - 10)
const uint32SafeToMultiply10 = uint32(0xffffffff)/10 - 10

func init() {
	intDigits = make([]int8, 256)
	for i := 0; i < len(floatDigits); i++ {
		intDigits[i] = invalidCharForNumber
	}
	for i := int8('0'); i <= int8('9'); i++ {
		intDigits[i] = i - int8('0')
	}
}

// ReadUint reads a json object as Uint
func (iter *Iterator) ReadUint() (ret uint) {
	val := iter.ReadUint64()
	converted := uint(val)
	if uint64(converted) != val {
		iter.reportError("ReadUint", "int overflow")
		return
	}
	return converted
}

// ReadUint8 reads a json object as Uint8
func (iter *Iterator) ReadUint8() (ret uint8) {
	val := iter.ReadUint64()
	converted := uint8(val)
	if uint64(converted) != val {
		iter.reportError("ReadUint8", "int overflow")
		return
	}
	return converted
}

// ReadUint16 reads a json object as Uint16
func (iter *Iterator) ReadUint16() (ret uint16) {
	val := iter.ReadUint64()
	converted := uint16(val)
	if uint64(converted) != val {
		iter.reportError("ReadUint16", "int overflow")
		return
	}
	return converted
}

// ReadUint64 reads a json object as Uint64
func (iter *Iterator) ReadUint64() (ret uint64) {
	c := iter.nextToken()
	v := hexDigits[c]
	if v == 0 {
		return 0 // single zero
	}
	if v == 255 {
		iter.reportError("ReadUint64", "unexpected character")
		return
	}
	for {
		if ret >= cutoffUint64 {
			iter.reportError("ReadUint64", "overflow")
			return
		}
		ret = ret*10 + uint64(v)
		c = iter.readByte()
		v = hexDigits[c]
		if v == 255 {
			iter.unreadByte()
			break
		}
	}
	return ret
}

// ReadInt reads a json object as Int
func (iter *Iterator) ReadInt() (ret int) {
	val := iter.ReadInt64()
	converted := int(val)
	if int64(converted) != val {
		iter.reportError("ReadInt", "int overflow")
		return
	}
	return converted
}

// ReadInt8 reads a json object as Int8
func (iter *Iterator) ReadInt8() (ret int8) {
	val := iter.ReadInt64()
	converted := int8(val)
	if int64(converted) != val {
		iter.reportError("ReadInt8", "int overflow")
		return
	}
	return converted
}

// ReadInt16 reads a json object as Int16
func (iter *Iterator) ReadInt16() (ret int16) {
	val := iter.ReadInt64()
	converted := int16(val)
	if int64(converted) != val {
		iter.reportError("ReadInt16", "int overflow")
		return
	}
	return converted
}

func (iter *Iterator) ReadInt32() int32 {
	c := iter.nextToken()
	if c == '-' {
		return -int32(iter.readUint32(int32SafeToMultiply10, iter.readByte()))
	} else {
		return int32(iter.readUint32(int32SafeToMultiply10, c))
	}
}

func (iter *Iterator) ReadUint32() uint32 {
	return iter.readUint32(uint32SafeToMultiply10, iter.nextToken())
}

func (iter *Iterator) readUint32(safeToMultiply10 uint32, c byte) (ret uint32) {
	ind := intDigits[c]
	if ind == 0 {
		return 0 // single zero
	}
	if ind == invalidCharForNumber {
		iter.reportError("readUint32", "unexpected character: " + string([]byte{byte(ind)}))
		return
	}
	value := uint32(ind)
	for {
		for i := iter.head; i < iter.tail; i++ {
			if value > safeToMultiply10 {
				iter.reportError("readUint32", "overflow")
				return
			}
			ind = intDigits[iter.buf[i]]
			if ind == invalidCharForNumber {
				return value
			}
			value =  (value << 3) + (value << 1) + uint32(ind)
		}
		if (!iter.loadMore()) {
			return value
		}
	}
}

// ReadInt64 reads a json object as Int64
func (iter *Iterator) ReadInt64() (ret int64) {
	c := iter.nextToken()
	if iter.Error != nil {
		return
	}

	/* optional leading minus */
	if c == '-' {
		n := iter.ReadUint64()
		return -int64(n)
	}
	iter.unreadByte()
	n := iter.ReadUint64()
	return int64(n)
}
