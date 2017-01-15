package jsoniter

var intDigits []int8

const int8SafeToMultiply10 = uint32(int8(0x7f) / 10 - 10)
const uint8SafeToMultiply10 = uint32(0xff) / 10 - 10
const int16SafeToMultiply10 = uint32(int16(0x7fff) / 10 - 10)
const uint16SafeToMultiply10 = uint32(0xffff) / 10 - 10
const int32SafeToMultiply10 = uint32(int32(0x7fffffff) / 10 - 10)
const uint32SafeToMultiply10 = uint32(0xffffffff) / 10 - 10
const uint64SafeToMultiple10 = uint64(0xffffffffffffffff) / 10 - 10
const int64SafeToMultiple10 = uint64(int64(0x7fffffffffffffff) / 10 - 10)

func init() {
	intDigits = make([]int8, 256)
	for i := 0; i < len(floatDigits); i++ {
		intDigits[i] = invalidCharForNumber
	}
	for i := int8('0'); i <= int8('9'); i++ {
		intDigits[i] = i - int8('0')
	}
}

func (iter *Iterator) ReadUint() uint {
	return uint(iter.ReadUint64())
}

func (iter *Iterator) ReadInt() int {
	return int(iter.ReadInt64())
}

func (iter *Iterator) ReadInt8() int8 {
	c := iter.nextToken()
	if c == '-' {
		return -int8(iter.readUint32(int8SafeToMultiply10, iter.readByte()))
	} else {
		return int8(iter.readUint32(int8SafeToMultiply10, c))
	}
}

func (iter *Iterator) ReadUint8() (ret uint8) {
	return uint8(iter.readUint32(uint8SafeToMultiply10, iter.nextToken()))
}

func (iter *Iterator) ReadInt16() int16 {
	c := iter.nextToken()
	if c == '-' {
		return -int16(iter.readUint32(int16SafeToMultiply10, iter.readByte()))
	} else {
		return int16(iter.readUint32(int16SafeToMultiply10, c))
	}
}

func (iter *Iterator) ReadUint16() uint16 {
	return uint16(iter.readUint32(uint16SafeToMultiply10, iter.nextToken()))
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
			ind = intDigits[iter.buf[i]]
			if ind == invalidCharForNumber {
				return value
			}
			if value > safeToMultiply10 {
				value2 := (value << 3) + (value << 1) + uint32(ind)
				if value2 < safeToMultiply10 * 10 {
					iter.reportError("readUint32", "overflow")
					return
				} else {
					value = value2
					continue
				}
			}
			value = (value << 3) + (value << 1) + uint32(ind)
		}
		if (!iter.loadMore()) {
			return value
		}
	}
}

func (iter *Iterator) ReadInt64() int64 {
	c := iter.nextToken()
	if c == '-' {
		return -int64(iter.readUint64(int64SafeToMultiple10, iter.readByte()))
	} else {
		return int64(iter.readUint64(int64SafeToMultiple10, c))
	}
}

func (iter *Iterator) ReadUint64() uint64 {
	return iter.readUint64(uint64SafeToMultiple10, iter.nextToken())
}

func (iter *Iterator) readUint64(safeToMultiply10 uint64, c byte) (ret uint64) {
	ind := intDigits[c]
	if ind == 0 {
		return 0 // single zero
	}
	if ind == invalidCharForNumber {
		iter.reportError("readUint64", "unexpected character: " + string([]byte{byte(ind)}))
		return
	}
	value := uint64(ind)
	for {
		for i := iter.head; i < iter.tail; i++ {
			ind = intDigits[iter.buf[i]]
			if ind == invalidCharForNumber {
				return value
			}
			if value > safeToMultiply10 {
				value2 := (value << 3) + (value << 1) + uint64(ind)
				if value2 < safeToMultiply10 * 10 {
					iter.reportError("readUint64", "overflow")
					return
				} else {
					value = value2
					continue
				}
			}
			value = (value << 3) + (value << 1) + uint64(ind)
		}
		if (!iter.loadMore()) {
			return value
		}
	}
}
