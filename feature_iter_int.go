package jsoniter

import (
	"strconv"
)

var intDigits []int8

const uint32SafeToMultiply10 = uint32(0xffffffff) / 10 - 1
const uint64SafeToMultiple10 = uint64(0xffffffffffffffff) / 10 - 1
const int64Max = uint64(0x7fffffffffffffff)
const int32Max = uint32(0x7fffffff)
const int16Max = uint32(0x7fff)
const uint16Max = uint32(0xffff)
const int8Max = uint32(0x7fff)
const uint8Max = uint32(0xffff)

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

func (iter *Iterator) ReadInt8() (ret int8) {
	c := iter.nextToken()
	if c == '-' {
		val := iter.readUint32(iter.readByte())
		if val > int8Max + 1 {
			iter.reportError("ReadInt8", "overflow: " + strconv.FormatInt(int64(val), 10))
			return
		}
		return -int8(val)
	} else {
		val := iter.readUint32(c)
		if val > int8Max {
			iter.reportError("ReadInt8", "overflow: " + strconv.FormatInt(int64(val), 10))
			return
		}
		return int8(val)
	}
}

func (iter *Iterator) ReadUint8() (ret uint8) {
	val := iter.readUint32(iter.nextToken())
	if val > uint8Max {
		iter.reportError("ReadUint8", "overflow: " + strconv.FormatInt(int64(val), 10))
		return
	}
	return uint8(val)
}

func (iter *Iterator) ReadInt16() (ret int16) {
	c := iter.nextToken()
	if c == '-' {
		val := iter.readUint32(iter.readByte())
		if val > int16Max + 1 {
			iter.reportError("ReadInt16", "overflow: " + strconv.FormatInt(int64(val), 10))
			return
		}
		return -int16(val)
	} else {
		val := iter.readUint32(c)
		if val > int16Max {
			iter.reportError("ReadInt16", "overflow: " + strconv.FormatInt(int64(val), 10))
			return
		}
		return int16(val)
	}
}

func (iter *Iterator) ReadUint16() (ret uint16) {
	val := iter.readUint32(iter.nextToken())
	if val > uint16Max {
		iter.reportError("ReadUint16", "overflow: " + strconv.FormatInt(int64(val), 10))
		return
	}
	return uint16(val)
}

func (iter *Iterator) ReadInt32() (ret int32) {
	c := iter.nextToken()
	if c == '-' {
		val := iter.readUint32(iter.readByte())
		if val > int32Max + 1 {
			iter.reportError("ReadInt32", "overflow: " + strconv.FormatInt(int64(val), 10))
			return
		}
		return -int32(val)
	} else {
		val := iter.readUint32(c)
		if val > int32Max {
			iter.reportError("ReadInt32", "overflow: " + strconv.FormatInt(int64(val), 10))
			return
		}
		return int32(val)
	}
}

func (iter *Iterator) ReadUint32() (ret uint32) {
	return iter.readUint32(iter.nextToken())
}

func (iter *Iterator) readUint32(c byte) (ret uint32) {
	ind := intDigits[c]
	if ind == 0 {
		return 0 // single zero
	}
	if ind == invalidCharForNumber {
		iter.reportError("readUint32", "unexpected character: " + string([]byte{byte(ind)}))
		return
	}
	value := uint32(ind)
	if iter.tail - iter.head > 10 {
		i := iter.head
		ind2 := intDigits[iter.buf[i]]
		if ind2 == invalidCharForNumber {
			iter.head = i
			return value
		}
		i++
		ind3 := intDigits[iter.buf[i]]
		if ind3 == invalidCharForNumber {
			iter.head = i
			return value * 10 + uint32(ind2)
		}
		//iter.head = i + 1
		//value = value * 100 + uint32(ind2) * 10 + uint32(ind3)
		i++
		ind4 := intDigits[iter.buf[i]]
		if ind4 == invalidCharForNumber {
			iter.head = i
			return value * 100 + uint32(ind2) * 10 + uint32(ind3)
		}
		i++
		ind5 := intDigits[iter.buf[i]]
		if ind5 == invalidCharForNumber {
			iter.head = i
			return value * 1000 + uint32(ind2) * 100 + uint32(ind3) * 10 + uint32(ind4)
		}
		i++
		ind6 := intDigits[iter.buf[i]]
		if ind6 == invalidCharForNumber {
			iter.head = i
			return value * 10000 + uint32(ind2) * 1000 + uint32(ind3) * 100 + uint32(ind4) * 10 + uint32(ind5)
		}
		i++
		ind7 := intDigits[iter.buf[i]]
		if ind7 == invalidCharForNumber {
			iter.head = i
			return value * 100000 + uint32(ind2) * 10000 + uint32(ind3) * 1000 + uint32(ind4) * 100 + uint32(ind5) * 10 + uint32(ind6)
		}
		i++
		ind8 := intDigits[iter.buf[i]]
		if ind8 == invalidCharForNumber {
			iter.head = i
			return value * 1000000 + uint32(ind2) * 100000 + uint32(ind3) * 10000 + uint32(ind4) * 1000 + uint32(ind5) * 100 + uint32(ind6) * 10 + uint32(ind7)
		}
		i++
		ind9 := intDigits[iter.buf[i]]
		value = value * 10000000 + uint32(ind2) * 1000000 + uint32(ind3) * 100000 + uint32(ind4) * 10000 + uint32(ind5) * 1000 + uint32(ind6) * 100 + uint32(ind7) * 10 + uint32(ind8)
		iter.head = i
		if ind9 == invalidCharForNumber {
			return value
		}
	}
	for {
		for i := iter.head; i < iter.tail; i++ {
			ind = intDigits[iter.buf[i]]
			if ind == invalidCharForNumber {
				iter.head = i
				return value
			}
			if value > uint32SafeToMultiply10 {
				value2 := (value << 3) + (value << 1) + uint32(ind)
				if value2 < value {
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

func (iter *Iterator) ReadInt64() (ret int64) {
	c := iter.nextToken()
	if c == '-' {
		val := iter.readUint64(iter.readByte())
		if val > int64Max + 1 {
			iter.reportError("ReadInt64", "overflow: " + strconv.FormatUint(uint64(val), 10))
			return
		}
		return -int64(val)
	} else {
		val := iter.readUint64(c)
		if val > int64Max {
			iter.reportError("ReadInt64", "overflow: " + strconv.FormatUint(uint64(val), 10))
			return
		}
		return int64(val)
	}
}

func (iter *Iterator) ReadUint64() uint64 {
	return iter.readUint64(iter.nextToken())
}

func (iter *Iterator) readUint64(c byte) (ret uint64) {
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
				iter.head = i
				return value
			}
			if value > uint64SafeToMultiple10 {
				value2 := (value << 3) + (value << 1) + uint64(ind)
				if value2 < value {
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
