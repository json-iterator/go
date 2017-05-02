package jsoniter

import (
	"io"
	"strconv"
	"unsafe"
)

var floatDigits []int8
const invalidCharForNumber = int8(-1)
const endOfNumber = int8(-2)
const dotInNumber = int8(-3)

func init() {
	floatDigits = make([]int8, 256)
	for i := 0; i < len(floatDigits); i++ {
		floatDigits[i] = invalidCharForNumber
	}
	for i := int8('0'); i <= int8('9'); i++ {
		floatDigits[i] = i - int8('0')
	}
	floatDigits[','] = endOfNumber
	floatDigits[']'] = endOfNumber
	floatDigits['}'] = endOfNumber
	floatDigits[' '] = endOfNumber
	floatDigits['\t'] = endOfNumber
	floatDigits['\n'] = endOfNumber
	floatDigits['.'] = dotInNumber
}

func (iter *Iterator) ReadFloat32() (ret float32) {
	c := iter.nextToken()
	if c == '-' {
		return -iter.readPositiveFloat32()
	} else {
		iter.unreadByte()
		return iter.readPositiveFloat32()
	}
}

func (iter *Iterator) readPositiveFloat32() (ret float32) {
	value := uint64(0)
	c := byte(' ')
	i := iter.head
	non_decimal_loop:
	for ; i < iter.tail; i++ {
		c = iter.buf[i]
		ind := floatDigits[c]
		switch ind {
		case invalidCharForNumber:
			return iter.readFloat32SlowPath()
		case endOfNumber:
			iter.head = i
			return float32(value)
		case dotInNumber:
			break non_decimal_loop
		}
		if value > uint64SafeToMultiple10 {
			return iter.readFloat32SlowPath()
		}
		value = (value << 3) + (value << 1) + uint64(ind); // value = value * 10 + ind;
	}
	if c == '.' {
		i++
		decimalPlaces := 0;
		for ; i < iter.tail; i++ {
			c = iter.buf[i]
			ind := floatDigits[c];
			switch ind {
			case endOfNumber:
				if decimalPlaces > 0 && decimalPlaces < len(POW10) {
					iter.head = i
					return float32(float64(value) / float64(POW10[decimalPlaces]))
				}
				// too many decimal places
			return iter.readFloat32SlowPath()
			case invalidCharForNumber:
				fallthrough
			case dotInNumber:
				return iter.readFloat32SlowPath()
			}
			decimalPlaces++
			if value > uint64SafeToMultiple10 {
				return iter.readFloat32SlowPath()
			}
			value = (value << 3) + (value << 1) + uint64(ind)
		}
	}
	return iter.readFloat32SlowPath()
}

func (iter *Iterator) readFloat32SlowPath() (ret float32) {
	strBuf := [16]byte{}
	str := strBuf[0:0]
	load_loop:
	for {
		for i := iter.head; i < iter.tail; i++ {
			c := iter.buf[i]
			switch c {
			case '-', '.', 'e', 'E', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				str = append(str, c)
				continue
			default:
				iter.head = i
				break load_loop
			}
		}
		if !iter.loadMore() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		return
	}
	val, err := strconv.ParseFloat(*(*string)(unsafe.Pointer(&str)), 32)
	if err != nil {
		iter.Error = err
		return
	}
	return float32(val)
}

func (iter *Iterator) ReadFloat64() (ret float64) {
	c := iter.nextToken()
	if c == '-' {
		return -iter.readPositiveFloat64()
	} else {
		iter.unreadByte()
		return iter.readPositiveFloat64()
	}
}

func (iter *Iterator) readPositiveFloat64() (ret float64) {
	value := uint64(0)
	c := byte(' ')
	i := iter.head
	non_decimal_loop:
	for ; i < iter.tail; i++ {
		c = iter.buf[i]
		ind := floatDigits[c]
		switch ind {
		case invalidCharForNumber:
			return iter.readFloat64SlowPath()
		case endOfNumber:
			iter.head = i
			return float64(value)
		case dotInNumber:
			break non_decimal_loop
		}
		if value > uint64SafeToMultiple10 {
			return iter.readFloat64SlowPath()
		}
		value = (value << 3) + (value << 1) + uint64(ind); // value = value * 10 + ind;
	}
	if c == '.' {
		i++
		decimalPlaces := 0;
		for ; i < iter.tail; i++ {
			c = iter.buf[i]
			ind := floatDigits[c];
			switch ind {
			case endOfNumber:
				if decimalPlaces > 0 && decimalPlaces < len(POW10) {
					iter.head = i
					return float64(value) / float64(POW10[decimalPlaces])
				}
				// too many decimal places
				return iter.readFloat64SlowPath()
			case invalidCharForNumber:
				fallthrough
			case dotInNumber:
				return iter.readFloat64SlowPath()
			}
			decimalPlaces++
			if value > uint64SafeToMultiple10 {
				return iter.readFloat64SlowPath()
			}
			value = (value << 3) + (value << 1) + uint64(ind)
		}
	}
	return iter.readFloat64SlowPath()
}

func (iter *Iterator) readFloat64SlowPath() (ret float64) {
	strBuf := [16]byte{}
	str := strBuf[0:0]
	load_loop:
	for {
		for i := iter.head; i < iter.tail; i++ {
			c := iter.buf[i]
			switch c {
			case '-', '.', 'e', 'E', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				str = append(str, c)
				continue
			default:
				iter.head = i
				break load_loop
			}
		}
		if !iter.loadMore() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		return
	}
	val, err := strconv.ParseFloat(*(*string)(unsafe.Pointer(&str)), 64)
	if err != nil {
		iter.Error = err
		return
	}
	return val
}
