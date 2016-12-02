package jsoniter

import (
	"io"
	"fmt"
	"unicode/utf16"
	"strconv"
)

type Iterator struct {
	reader io.Reader
	buf    []byte
	head   int
	tail   int
	Error  error
}

func Parse(reader io.Reader, bufSize int) *Iterator {
	iter := &Iterator{
		reader: reader,
		buf: make([]byte, bufSize),
		head: 0,
		tail: 0,
	}
	iter.skipWhitespaces()
	return iter
}

func ParseBytes(input []byte) *Iterator {
	iter := &Iterator{
		reader: nil,
		buf: input,
		head: 0,
		tail: len(input),
	}
	iter.skipWhitespaces()
	return iter
}

func ParseString(input string) *Iterator {
	return ParseBytes([]byte(input))
}

func (iter *Iterator) skipWhitespaces() {
	c := iter.readByte()
	for c == ' ' || c == '\n' {
		c = iter.readByte()
	}
	iter.unreadByte()
}

func (iter *Iterator) ReportError(operation string, msg string) {
	peekStart := iter.head - 10
	if peekStart < 0 {
		peekStart = 0
	}
	iter.Error = fmt.Errorf("%s: %s, parsing %v ...%s... at %s", operation, msg, iter.head,
		string(iter.buf[peekStart: iter.head]), string(iter.buf[0:iter.tail]))
}

func (iter *Iterator) readByte() (ret byte) {
	if iter.head == iter.tail {
		if iter.reader == nil {
			iter.Error = io.EOF
			return
		}
		n, err := iter.reader.Read(iter.buf)
		if err != nil {
			iter.Error = err
			return
		}
		if n == 0 {
			iter.Error = io.EOF
			return
		}
		iter.head = 0
		iter.tail = n
	}
	ret = iter.buf[iter.head]
	iter.head += 1
	return ret
}

func (iter *Iterator) unreadByte() {
	if iter.head == 0 {
		iter.ReportError("unreadByte", "unread too many bytes")
		return
	}
	iter.head -= 1
	return
}

const maxUint64 = (1 << 64 - 1)
const cutoffUint64 = maxUint64 / 10 + 1
const maxUint32 = (1 << 32 - 1)
const cutoffUint32 = maxUint32 / 10 + 1

func (iter *Iterator) ReadUint64() (ret uint64) {
	c := iter.readByte()
	if iter.Error != nil {
		return
	}

	/* a single zero, or a series of integers */
	if c == '0' {
		return 0
	} else if c >= '1' && c <= '9' {
		for c >= '0' && c <= '9' {
			var v byte
			v = c - '0'
			if ret >= cutoffUint64 {
				iter.ReportError("ReadUint64", "overflow")
				return
			}
			ret = ret * uint64(10) + uint64(v)
			c = iter.readByte()
			if iter.Error != nil {
				if iter.Error == io.EOF {
					break
				} else {
					return 0
				}
			}
		}
		if iter.Error != io.EOF {
			iter.unreadByte()
		}
	} else {
		iter.ReportError("ReadUint64", "expects 0~9")
		return
	}
	return ret
}

func (iter *Iterator) ReadInt64() (ret int64) {
	c := iter.readByte()
	if iter.Error != nil {
		return
	}

	/* optional leading minus */
	if c == '-' {
		n := iter.ReadUint64()
		return -int64(n)
	} else {
		iter.unreadByte()
		n := iter.ReadUint64()
		return int64(n)
	}
}

func (iter *Iterator) ReadString() (ret string) {
	str := make([]byte, 0, 10)
	c := iter.readByte()
	if iter.Error != nil {
		return
	}
	switch c {
	case 'n':
		iter.skipNull()
		if iter.Error != nil {
			return
		}
		return ""
	case '"':
		// nothing
	default:
		iter.ReportError("ReadString", `expects " or n`)
		return
	}
	for {
		c = iter.readByte()
		if iter.Error != nil {
			return
		}
		switch c {
		case '\\':
			c = iter.readByte()
			if iter.Error != nil {
				return
			}
			switch c {
			case 'u':
				r := iter.readU4()
				if iter.Error != nil {
					return
				}
				if utf16.IsSurrogate(r) {
					c = iter.readByte()
					if iter.Error != nil {
						return
					}
					if c != '\\' {
						iter.ReportError("ReadString",
							`expects \u after utf16 surrogate, but \ not found`)
						return
					}
					c = iter.readByte()
					if iter.Error != nil {
						return
					}
					if c != 'u' {
						iter.ReportError("ReadString",
							`expects \u after utf16 surrogate, but \u not found`)
						return
					}
					r2 := iter.readU4()
					if iter.Error != nil {
						return
					}
					combined := utf16.DecodeRune(r, r2)
					str = appendRune(str, combined)
				} else {
					str = appendRune(str, r)
				}
			case '"':
				str = append(str, '"')
			case '\\':
				str = append(str, '\\')
			case '/':
				str = append(str, '/')
			case 'b':
				str = append(str, '\b')
			case 'f':
				str = append(str, '\f')
			case 'n':
				str = append(str, '\n')
			case 'r':
				str = append(str, '\r')
			case 't':
				str = append(str, '\t')
			default:
				iter.ReportError("ReadString",
					`invalid escape char after \`)
				return
			}
		case '"':
			return string(str)
		default:
			str = append(str, c)
		}
	}
}

func (iter *Iterator) readU4() (ret rune) {
	for i := 0; i < 4; i++ {
		c := iter.readByte()
		if iter.Error != nil {
			return
		}
		if (c >= '0' && c <= '9') {
			if ret >= cutoffUint32 {
				iter.ReportError("readU4", "overflow")
				return
			}
			ret = ret * 16 + rune(c - '0')
		} else if ((c >= 'a' && c <= 'f') ) {
			if ret >= cutoffUint32 {
				iter.ReportError("readU4", "overflow")
				return
			}
			ret = ret * 16 + rune(c - 'a' + 10)
		} else {
			iter.ReportError("readU4", "expects 0~9 or a~f")
			return
		}
	}
	return ret
}

const (
	t1 = 0x00 // 0000 0000
	tx = 0x80 // 1000 0000
	t2 = 0xC0 // 1100 0000
	t3 = 0xE0 // 1110 0000
	t4 = 0xF0 // 1111 0000
	t5 = 0xF8 // 1111 1000

	maskx = 0x3F // 0011 1111
	mask2 = 0x1F // 0001 1111
	mask3 = 0x0F // 0000 1111
	mask4 = 0x07 // 0000 0111

	rune1Max = 1 << 7 - 1
	rune2Max = 1 << 11 - 1
	rune3Max = 1 << 16 - 1

	surrogateMin = 0xD800
	surrogateMax = 0xDFFF

	MaxRune = '\U0010FFFF' // Maximum valid Unicode code point.
	RuneError = '\uFFFD'     // the "error" Rune or "Unicode replacement character"
)

func appendRune(p []byte, r rune) []byte {
	// Negative values are erroneous. Making it unsigned addresses the problem.
	switch i := uint32(r); {
	case i <= rune1Max:
		p = append(p, byte(r))
		return p
	case i <= rune2Max:
		p = append(p, t2 | byte(r >> 6))
		p = append(p, tx | byte(r) & maskx)
		return p
	case i > MaxRune, surrogateMin <= i && i <= surrogateMax:
		r = RuneError
		fallthrough
	case i <= rune3Max:
		p = append(p, t3 | byte(r >> 12))
		p = append(p, tx | byte(r >> 6) & maskx)
		p = append(p, tx | byte(r) & maskx)
		return p
	default:
		p = append(p, t4 | byte(r >> 18))
		p = append(p, tx | byte(r >> 12) & maskx)
		p = append(p, tx | byte(r >> 6) & maskx)
		p = append(p, tx | byte(r) & maskx)
		return p
	}
}

func (iter *Iterator) ReadArray() (ret bool) {
	iter.skipWhitespaces()
	c := iter.readByte()
	if iter.Error != nil {
		return
	}
	switch c {
	case 'n': {
		iter.skipNull()
		if iter.Error != nil {
			return
		}
		return false // null
	}
	case '[': {
		iter.skipWhitespaces()
		c = iter.readByte()
		if iter.Error != nil {
			return
		}
		if c == ']' {
			return false
		} else {
			iter.unreadByte()
			return true
		}
	}
	case ']': return false
	case ',':
		iter.skipWhitespaces()
		return true
	default:
		iter.ReportError("ReadArray", "expect [ or , or ] or n")
		return
	}
}

func (iter *Iterator) ReadObject() (ret string) {
	iter.skipWhitespaces()
	c := iter.readByte()
	if iter.Error != nil {
		return
	}
	switch c {
	case 'n': {
		iter.skipNull()
		if iter.Error != nil {
			return
		}
		return "" // null
	}
	case '{': {
		iter.skipWhitespaces()
		c = iter.readByte()
		if iter.Error != nil {
			return
		}
		switch c {
		case '}':
			return "" // end of object
		case '"':
			iter.unreadByte()
			field := iter.readObjectField()
			if iter.Error != nil {
				return
			}
			return field
		default:
			iter.ReportError("ReadObject", `expect " after {`)
			return
		}
	}
	case ',':
		iter.skipWhitespaces()
		field := iter.readObjectField()
		if iter.Error != nil {
			return
		}
		return field
	case '}':
		return "" // end of object
	default:
		iter.ReportError("ReadObject", `expect { or , or } or n`)
		return
	}
}

func (iter *Iterator) readObjectField() (ret string) {
	field := iter.ReadString()
	if iter.Error != nil {
		return
	}
	iter.skipWhitespaces()
	c := iter.readByte()
	if iter.Error != nil {
		return
	}
	if c != ':' {
		iter.ReportError("ReadObject", "expect : after object field")
		return
	}
	iter.skipWhitespaces()
	return field
}

func (iter *Iterator) ReadFloat64() (ret float64) {
	str := make([]byte, 0, 10)
	for c := iter.readByte(); iter.Error == nil; c = iter.readByte() {
		switch c {
		case '-', '+', '.', 'e', 'E', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			str = append(str, c)
		default:
			iter.unreadByte()
			val, err := strconv.ParseFloat(string(str), 64)
			if err != nil {
				iter.Error = err
				return
			}
			return val
		}
	}
	if iter.Error == io.EOF {
		val, err := strconv.ParseFloat(string(str), 64)
		if err != nil {
			iter.Error = err
			return
		}
		return val
	}
	return
}

func (iter *Iterator) ReadBool() (ret bool) {
	c := iter.readByte()
	if iter.Error != nil {
		return
	}
	switch c {
	case 't':
		iter.skipTrue()
		if iter.Error != nil {
			return
		}
		return true
	case 'f':
		iter.skipFalse()
		if iter.Error != nil {
			return
		}
		return false
	default:
		iter.ReportError("ReadBool", "expect t or f")
		return
	}
}

func (iter *Iterator) skipTrue() {
	c := iter.readByte()
	if iter.Error != nil {
		return
	}
	if c != 'r' {
		iter.ReportError("skipTrue", "expect r of true")
		return
	}
	c = iter.readByte()
	if iter.Error != nil {
		return
	}
	if c != 'u' {
		iter.ReportError("skipTrue", "expect u of true")
		return
	}
	c = iter.readByte()
	if iter.Error != nil {
		return
	}
	if c != 'e' {
		iter.ReportError("skipTrue", "expect e of true")
		return
	}
}

func (iter *Iterator) skipFalse() {
	c := iter.readByte()
	if iter.Error != nil {
		return
	}
	if c != 'a' {
		iter.ReportError("skipFalse", "expect a of false")
		return
	}
	c = iter.readByte()
	if iter.Error != nil {
		return
	}
	if c != 'l' {
		iter.ReportError("skipFalse", "expect l of false")
		return
	}
	c = iter.readByte()
	if iter.Error != nil {
		return
	}
	if c != 's' {
		iter.ReportError("skipFalse", "expect s of false")
		return
	}
	c = iter.readByte()
	if iter.Error != nil {
		return
	}
	if c != 'e' {
		iter.ReportError("skipFalse", "expect e of false")
		return
	}
}

func (iter *Iterator) ReadNull() (ret bool) {
	c := iter.readByte()
	if iter.Error != nil {
		return
	}
	if c == 'n' {
		iter.skipNull()
		if iter.Error != nil {
			return
		}
		return true
	}
	iter.unreadByte()
	return false
}

func (iter *Iterator) skipNull() {
	c := iter.readByte()
	if iter.Error != nil {
		return
	}
	if c != 'u' {
		iter.ReportError("skipNull", "expect u of null")
		return
	}
	c = iter.readByte()
	if iter.Error != nil {
		return
	}
	if c != 'l' {
		iter.ReportError("skipNull", "expect l of null")
		return
	}
	c = iter.readByte()
	if iter.Error != nil {
		return
	}
	if c != 'l' {
		iter.ReportError("skipNull", "expect l of null")
		return
	}
}

func (iter *Iterator) Skip() {
	c := iter.readByte()
	if iter.Error != nil {
		return
	}
	switch c {
	case '"':
		iter.skipString()
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		iter.skipNumber()
	case '[':
		iter.skipArray()
	case '{':
		iter.skipObject()
	case 't':
		iter.skipTrue()
	case 'f':
		iter.skipFalse()
	case 'n':
		iter.skipNull()
	default:
		iter.ReportError("Skip", fmt.Sprintf("do not know how to skip: %v", c))
		return
	}
}

func (iter *Iterator) skipString() {
	for c := iter.readByte(); iter.Error == nil; c = iter.readByte() {
		switch c {
		case '"':
			return // end of string found
		case '\\':
			iter.readByte() // " after \\ does not count
			if iter.Error != nil {
				return
			}
		}
	}
}

func (iter *Iterator) skipNumber() {
	for c := iter.readByte(); iter.Error == nil; c = iter.readByte() {
		switch c {
		case '-', '+', '.', 'e', 'E', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			continue
		default:
			iter.unreadByte()
			return
		}
	}
}

func (iter *Iterator) skipArray() {
	for {
		iter.skipWhitespaces()
		c := iter.readByte()
		if iter.Error != nil {
			return
		}
		if c == ']' {
			return
		}
		iter.unreadByte()
		iter.Skip()
		iter.skipWhitespaces()
		c = iter.readByte()
		if iter.Error != nil {
			return
		}
		switch c {
		case ',':
			iter.skipWhitespaces()
			continue
		case ']':
			return
		default:
			iter.ReportError("skipArray", "expects , or ]")
			return
		}
	}
}

func (iter *Iterator) skipObject() {
	iter.skipWhitespaces()
	c := iter.readByte()
	if iter.Error != nil {
		return
	}
	if c == '}' {
		return // end of object
	} else {
		iter.unreadByte()
	}
	for {
		iter.skipWhitespaces()
		c := iter.readByte()
		if iter.Error != nil {
			return
		}
		if c != '"' {
			iter.ReportError("skipObject", `expects "`)
			return
		}
		iter.skipString()
		iter.skipWhitespaces()
		c = iter.readByte()
		if iter.Error != nil {
			return
		}
		if c != ':' {
			iter.ReportError("skipObject", `expects :`)
			return
		}
		iter.skipWhitespaces()
		iter.Skip()
		iter.skipWhitespaces()
		c = iter.readByte()
		if iter.Error != nil {
			return
		}
		switch c {
		case ',':
			iter.skipWhitespaces()
			continue
		case '}':
			return // end of object
		default:
			iter.ReportError("skipObject", "expects , or }")
			return
		}
	}
}