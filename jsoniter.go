package jsoniter

import (
	"io"
	"fmt"
	"unicode/utf16"
	"strconv"
	"unsafe"
	"encoding/base64"
)

type ValueType int

const (
	Invalid ValueType = iota
	String
	Number
	Null
	Bool
	Array
	Object
)

var digits []byte
var valueTypes []ValueType

func init() {
	digits = make([]byte, 256)
	for i := 0; i < len(digits); i++ {
		digits[i] = 255
	}
	for i := '0'; i <= '9'; i++ {
		digits[i] = byte(i - '0');
	}
	for i := 'a'; i <= 'f'; i++ {
		digits[i] = byte((i - 'a') + 10);
	}
	for i := 'A'; i <= 'F'; i++ {
		digits[i] = byte((i - 'A') + 10);
	}
	valueTypes = make([]ValueType, 256)
	for i := 0; i < len(valueTypes); i++ {
		valueTypes[i] = Invalid
	}
	valueTypes['"'] = String;
	valueTypes['-'] = Number;
	valueTypes['0'] = Number;
	valueTypes['1'] = Number;
	valueTypes['2'] = Number;
	valueTypes['3'] = Number;
	valueTypes['4'] = Number;
	valueTypes['5'] = Number;
	valueTypes['6'] = Number;
	valueTypes['7'] = Number;
	valueTypes['8'] = Number;
	valueTypes['9'] = Number;
	valueTypes['t'] = Bool;
	valueTypes['f'] = Bool;
	valueTypes['n'] = Null;
	valueTypes['['] = Array;
	valueTypes['{'] = Object;
}

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

func (iter *Iterator) Reset(reader io.Reader) *Iterator {
	iter.reader = reader
	iter.head = 0
	iter.tail = 0
	iter.skipWhitespaces()
	return iter
}

func (iter *Iterator) ResetBytes(input []byte) *Iterator {
	// only for benchmarking
	iter.reader = nil
	iter.Error = nil
	iter.buf = input
	iter.head = 0
	iter.tail = len(input)
	iter.skipWhitespaces()
	return iter
}

func (iter *Iterator) WhatIsNext() ValueType {
	valueType := valueTypes[iter.readByte()];
	iter.unreadByte();
	return valueType;
}

func (iter *Iterator) skipWhitespaces() {
	for {
		for i := iter.head; i < iter.tail; i++ {
			c := iter.buf[i]
			switch c {
			case ' ', '\n', '\t', '\r':
				continue
			}
			iter.head = i
			return
		}
		if !iter.loadMore() {
			return
		}
	}
}

func (iter *Iterator) skipWhitespacesWithoutLoadMore() bool {
	for i := iter.head; i < iter.tail; i++ {
		c := iter.buf[i]
		switch c {
		case ' ', '\n', '\t', '\r':
			continue
		}
		iter.head = i
		return false
	}
	return true
}

func (iter *Iterator) nextToken() byte {
	// a variation of skip whitespaces, returning the next non-whitespace token
	for {
		for i := iter.head; i < iter.tail; i++ {
			c := iter.buf[i]
			switch c {
			case ' ', '\n', '\t', '\r':
				continue
			}
			iter.head = i+1
			return c
		}
		if !iter.loadMore() {
			return 0
		}
	}
}

func (iter *Iterator) ReportError(operation string, msg string) {
	if iter.Error != nil {
		return
	}
	peekStart := iter.head - 10
	if peekStart < 0 {
		peekStart = 0
	}
	iter.Error = fmt.Errorf("%s: %s, parsing %v ...%s... at %s", operation, msg, iter.head,
		string(iter.buf[peekStart: iter.head]), string(iter.buf[0:iter.tail]))
}

func (iter *Iterator) CurrentBuffer() string {
	peekStart := iter.head - 10
	if peekStart < 0 {
		peekStart = 0
	}
	return fmt.Sprintf("parsing %v ...|%s|... at %s", iter.head,
		string(iter.buf[peekStart: iter.head]), string(iter.buf[0:iter.tail]))
}

func (iter *Iterator) readByte() (ret byte) {
	if iter.head == iter.tail {
		if iter.loadMore() {
			ret = iter.buf[iter.head]
			iter.head++
			return ret
		} else {
			return 0
		}
	}
	ret = iter.buf[iter.head]
	iter.head++
	return ret
}

func (iter *Iterator) loadMore() bool {
	if iter.reader == nil {
		iter.Error = io.EOF
		return false
	}
	for {
		n, err := iter.reader.Read(iter.buf)
		if n == 0 {
			if err != nil {
				iter.Error = err
				return false
			} else {
				// n == 0, err == nil is not EOF
				continue
			}
		} else {
			iter.head = 0
			iter.tail = n
			return true
		}
	}
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

func (iter *Iterator) ReadUint() (ret uint) {
	val := iter.ReadUint64()
	converted := uint(val)
	if uint64(converted) != val {
		iter.ReportError("ReadUint", "int overflow")
		return
	}
	return converted
}

func (iter *Iterator) ReadUint8() (ret uint8) {
	val := iter.ReadUint64()
	converted := uint8(val)
	if uint64(converted) != val {
		iter.ReportError("ReadUint8", "int overflow")
		return
	}
	return converted
}

func (iter *Iterator) ReadUint16() (ret uint16) {
	val := iter.ReadUint64()
	converted := uint16(val)
	if uint64(converted) != val {
		iter.ReportError("ReadUint16", "int overflow")
		return
	}
	return converted
}

func (iter *Iterator) ReadUint32() (ret uint32) {
	val := iter.ReadUint64()
	converted := uint32(val)
	if uint64(converted) != val {
		iter.ReportError("ReadUint32", "int overflow")
		return
	}
	return converted
}

func (iter *Iterator) ReadUint64() (ret uint64) {
	c := iter.readByte()
	v := digits[c]
	if v == 0 {
		return 0 // single zero
	}
	if v == 255 {
		iter.ReportError("ReadUint64", "unexpected character")
		return
	}
	for {
		if ret >= cutoffUint64 {
			iter.ReportError("ReadUint64", "overflow")
			return
		}
		ret = ret * 10 + uint64(v)
		c = iter.readByte()
		v = digits[c]
		if v == 255 {
			iter.unreadByte()
			break
		}
	}
	return ret
}

func (iter *Iterator) ReadInt() (ret int) {
	val := iter.ReadInt64()
	converted := int(val)
	if int64(converted) != val {
		iter.ReportError("ReadInt", "int overflow")
		return
	}
	return converted
}

func (iter *Iterator) ReadInt8() (ret int8) {
	val := iter.ReadInt64()
	converted := int8(val)
	if int64(converted) != val {
		iter.ReportError("ReadInt8", "int overflow")
		return
	}
	return converted
}

func (iter *Iterator) ReadInt16() (ret int16) {
	val := iter.ReadInt64()
	converted := int16(val)
	if int64(converted) != val {
		iter.ReportError("ReadInt16", "int overflow")
		return
	}
	return converted
}

func (iter *Iterator) ReadInt32() (ret int32) {
	val := iter.ReadInt64()
	converted := int32(val)
	if int64(converted) != val {
		iter.ReportError("ReadInt32", "int overflow")
		return
	}
	return converted
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
	return string(iter.readStringAsBytes())
}


func (iter *Iterator) readStringAsBytes() (ret []byte) {
	c := iter.readByte()
	if c == 'n' {
		iter.skipUntilBreak()
		return
	}
	if c != '"' {
		iter.ReportError("ReadString", `expects " or n`)
		return
	}
	end := iter.findStringEndWithoutEscape()
	if end != -1 {
		// fast path: reuse the underlying buffer
		ret = iter.buf[iter.head:end-1]
		iter.head = end
		return ret
	}
	str := make([]byte, 0, 8)
	for iter.Error == nil {
		c = iter.readByte()
		if c == '"' {
			return str
		}
		if c == '\\' {
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
		} else {
			str = append(str, c)
		}
	}
	return
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
	c := iter.nextToken()
	if iter.Error != nil {
		return
	}
	switch c {
	case 'n': {
		iter.skipUntilBreak()
		return false // null
	}
	case '[': {
		c = iter.nextToken()
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
	c := iter.nextToken()
	if iter.Error != nil {
		return
	}
	switch c {
	case 'n': {
		iter.skipUntilBreak()
		if iter.Error != nil {
			return
		}
		return "" // null
	}
	case '{': {
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
			iter.ReportError("ReadObject", `expect " after {`)
			return
		}
	}
	case ',':
		iter.skipWhitespaces()
		return iter.readObjectField()
	case '}':
		return "" // end of object
	default:
		iter.ReportError("ReadObject", `expect { or , or } or n`)
		return
	}
}

func (iter *Iterator) readObjectField() (ret string) {
	str := iter.readStringAsBytes()
	if iter.skipWhitespacesWithoutLoadMore() {
		if ret == "" {
			ret = string(str);
		}
		if !iter.loadMore() {
			return
		}
	}
	if iter.buf[iter.head] != ':' {
		iter.ReportError("ReadObject", "expect : after object field")
		return
	}
	iter.head++
	if iter.skipWhitespacesWithoutLoadMore() {
		if ret == "" {
			ret = string(str);
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

func (iter *Iterator) ReadFloat32() (ret float32) {
	strBuf := [8]byte{}
	str := strBuf[0:0]
	hasMore := true
	for(hasMore) {
		for i := iter.head; i < iter.tail; i++ {
			c := iter.buf[i]
			switch c {
			case '-', '+', '.', 'e', 'E', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				str = append(str, c)
				continue
			default:
				hasMore = false
				break
			}
		}
		if hasMore {
			if !iter.loadMore() {
				break
			}
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
	strBuf := [8]byte{}
	str := strBuf[0:0]
	hasMore := true
	for(hasMore) {
		for i := iter.head; i < iter.tail; i++ {
			c := iter.buf[i]
			switch c {
			case '-', '+', '.', 'e', 'E', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				str = append(str, c)
				continue
			default:
				hasMore = false
				break
			}
		}
		if hasMore {
			if !iter.loadMore() {
				break
			}
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

func (iter *Iterator) ReadBool() (ret bool) {
	c := iter.readByte()
	if iter.Error != nil {
		return
	}
	switch c {
	case 't':
		iter.skipUntilBreak()
		return true
	case 'f':
		iter.skipUntilBreak()
		return false
	default:
		iter.ReportError("ReadBool", "expect t or f")
		return
	}
}

func (iter *Iterator) ReadBase64() (ret []byte) {
	src := iter.readStringAsBytes()
	if iter.Error != nil {
		return
	}
	b64 := base64.StdEncoding
	ret = make([]byte, b64.DecodedLen(len(src)))
	n, err := b64.Decode(ret, src)
	if err != nil {
		iter.Error = err
		return
	}
	return ret[:n]
}

func (iter *Iterator) ReadNull() (ret bool) {
	c := iter.readByte()
	if c == 'n' {
		iter.skipUntilBreak()
		return true
	}
	iter.unreadByte()
	return false
}

func (iter *Iterator) Skip() {
	c := iter.readByte()
	switch c {
	case '"':
		iter.skipString()
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 't', 'f', 'n':
		iter.skipUntilBreak()
	case '[':
		iter.skipArray()
	case '{':
		iter.skipObject()
	default:
		iter.ReportError("Skip", fmt.Sprintf("do not know how to skip: %v", c))
		return
	}
}

func (iter *Iterator) skipString() {
	for {
		end, escaped := iter.findStringEnd()
		if end == -1 {
			if !iter.loadMore() {
				return
			}
			if escaped {
				iter.head = 1 // skip the first char as last char read is \
			}
		} else {
			iter.head = end
			return
		}
	}
}


// adapted from: https://github.com/buger/jsonparser/blob/master/parser.go
// Tries to find the end of string
// Support if string contains escaped quote symbols.
func (iter *Iterator) findStringEnd() (int, bool) {
	escaped := false
	for i := iter.head; i < iter.tail; i++ {
		c := iter.buf[i]
		if c == '"' {
			if !escaped {
				return i + 1, false
			} else {
				j := i - 1
				for {
					if j < iter.head || iter.buf[j] != '\\' {
						// even number of backslashes
						// either end of buffer, or " found
						return i + 1, true
					}
					j--
					if j < iter.head || iter.buf[j] != '\\' {
						// odd number of backslashes
						// it is \" or \\\"
						break
					}
					j--
				}
			}
		} else if c == '\\' {
			escaped = true
		}
	}
	j := iter.tail - 1
	for {
		if j < iter.head || iter.buf[j] != '\\' {
			// even number of backslashes
			// either end of buffer, or " found
			return -1, false // do not end with \
		}
		j--
		if j < iter.head || iter.buf[j] != '\\' {
			// odd number of backslashes
			// it is \" or \\\"
			break
		}
		j--

	}
	return -1, true // end with \
}


func (iter *Iterator) findStringEndWithoutEscape() int {
	for i := iter.head; i < iter.tail; i++ {
		c := iter.buf[i]
		if c == '"' {
			return i + 1
		} else if c == '\\' {
			return -1
		}
	}
	return -1
}

func (iter *Iterator) skipArray() {
	level := 1
	for {
		for i := iter.head; i < iter.tail; i++ {
			switch iter.buf[i] {
			case '"': // If inside string, skip it
				iter.head = i + 1
				iter.skipString()
				i = iter.head - 1 // it will be i++ soon
			case '[': // If open symbol, increase level
				level++
			case ']': // If close symbol, increase level
				level--

				// If we have returned to the original level, we're done
				if level == 0 {
					iter.head = i + 1
					return
				}
			}
		}
		if (!iter.loadMore()) {
			return
		}
	}
}

func (iter *Iterator) skipObject() {
	level := 1
	for {
		for i := iter.head; i < iter.tail; i++ {
			switch iter.buf[i] {
			case '"': // If inside string, skip it
				iter.head = i + 1
				iter.skipString()
				i = iter.head - 1 // it will be i++ soon
			case '{': // If open symbol, increase level
				level++
			case '}': // If close symbol, increase level
				level--

				// If we have returned to the original level, we're done
				if level == 0 {
					iter.head = i + 1
					return
				}
			}
		}
		if (!iter.loadMore()) {
			return
		}
	}
}

func (iter *Iterator) skipUntilBreak() {
	// true, false, null, number
	for {
		for i := iter.head; i < iter.tail; i++ {
			c := iter.buf[i]
			switch c {
			case ' ', '\n', '\r', '\t', ',', '}', ']':
				iter.head = i
				return
			}
		}
		if (!iter.loadMore()) {
			return
		}
	}
}
