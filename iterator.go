package jsoniter

import (
	"encoding/base64"
	"fmt"
	"io"
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

var hexDigits []byte
var valueTypes []ValueType

func init() {
	hexDigits = make([]byte, 256)
	for i := 0; i < len(hexDigits); i++ {
		hexDigits[i] = 255
	}
	for i := '0'; i <= '9'; i++ {
		hexDigits[i] = byte(i - '0')
	}
	for i := 'a'; i <= 'f'; i++ {
		hexDigits[i] = byte((i - 'a') + 10)
	}
	for i := 'A'; i <= 'F'; i++ {
		hexDigits[i] = byte((i - 'A') + 10)
	}
	valueTypes = make([]ValueType, 256)
	for i := 0; i < len(valueTypes); i++ {
		valueTypes[i] = Invalid
	}
	valueTypes['"'] = String
	valueTypes['-'] = Number
	valueTypes['0'] = Number
	valueTypes['1'] = Number
	valueTypes['2'] = Number
	valueTypes['3'] = Number
	valueTypes['4'] = Number
	valueTypes['5'] = Number
	valueTypes['6'] = Number
	valueTypes['7'] = Number
	valueTypes['8'] = Number
	valueTypes['9'] = Number
	valueTypes['t'] = Bool
	valueTypes['f'] = Bool
	valueTypes['n'] = Null
	valueTypes['['] = Array
	valueTypes['{'] = Object
}

// Iterator is a fast and flexible JSON parser
type Iterator struct {
	reader io.Reader
	buf    []byte
	head   int
	tail   int
	Error  error
}

// Create creates an empty Iterator instance
func NewIterator() *Iterator {
	return &Iterator{
		reader: nil,
		buf:    nil,
		head:   0,
		tail:   0,
	}
}

// Parse parses a json buffer in io.Reader into an Iterator instance
func Parse(reader io.Reader, bufSize int) *Iterator {
	return &Iterator{
		reader: reader,
		buf:    make([]byte, bufSize),
		head:   0,
		tail:   0,
	}
}

// ParseBytes parses a json byte slice into an Iterator instance
func ParseBytes(input []byte) *Iterator {
	return &Iterator{
		reader: nil,
		buf:    input,
		head:   0,
		tail:   len(input),
	}
}

// ParseString parses a json string into an Iterator instance
func ParseString(input string) *Iterator {
	return ParseBytes([]byte(input))
}

// Reset can reset an Iterator instance for another json buffer in io.Reader
func (iter *Iterator) Reset(reader io.Reader) *Iterator {
	iter.reader = reader
	iter.head = 0
	iter.tail = 0
	return iter
}

// ResetBytes can reset an Iterator instance for another json byte slice
func (iter *Iterator) ResetBytes(input []byte) *Iterator {
	iter.reader = nil
	iter.Error = nil
	iter.buf = input
	iter.head = 0
	iter.tail = len(input)
	return iter
}

// WhatIsNext gets ValueType of relatively next json object
func (iter *Iterator) WhatIsNext() ValueType {
	valueType := valueTypes[iter.nextToken()]
	iter.unreadByte()
	return valueType
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
			iter.head = i + 1
			return c
		}
		if !iter.loadMore() {
			return 0
		}
	}
}

func (iter *Iterator) reportError(operation string, msg string) {
	if iter.Error != nil {
		return
	}
	peekStart := iter.head - 10
	if peekStart < 0 {
		peekStart = 0
	}
	iter.Error = fmt.Errorf("%s: %s, parsing %v ...%s... at %s", operation, msg, iter.head,
		string(iter.buf[peekStart:iter.head]), string(iter.buf[0:iter.tail]))
}

// CurrentBuffer gets current buffer as string
func (iter *Iterator) CurrentBuffer() string {
	peekStart := iter.head - 10
	if peekStart < 0 {
		peekStart = 0
	}
	return fmt.Sprintf("parsing %v ...|%s|... at %s", iter.head,
		string(iter.buf[peekStart:iter.head]), string(iter.buf[0:iter.tail]))
}

func (iter *Iterator) readByte() (ret byte) {
	if iter.head == iter.tail {
		if iter.loadMore() {
			ret = iter.buf[iter.head]
			iter.head++
			return ret
		}
		return 0
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
		iter.reportError("unreadByte", "unread too many bytes")
		return
	}
	iter.head--
	return
}

// ReadArray reads a json object as Array
func (iter *Iterator) ReadArray() (ret bool) {
	c := iter.nextToken()
	if iter.Error != nil {
		return
	}
	switch c {
	case 'n':
		iter.skipUntilBreak()
		return false // null
	case '[':
		c = iter.nextToken()
		if iter.Error != nil {
			return
		}
		if c == ']' {
			return false
		}
		iter.unreadByte()
		return true
	case ']':
		return false
	case ',':
		return true
	default:
		iter.reportError("ReadArray", "expect [ or , or ] or n, but found: " + string([]byte{c}))
		return
	}
}


// ReadBool reads a json object as Bool
func (iter *Iterator) ReadBool() (ret bool) {
	c := iter.nextToken()
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
		iter.reportError("ReadBool", "expect t or f")
		return
	}
}

// ReadBase64 reads a json object as Base64 in byte slice
func (iter *Iterator) ReadBase64() (ret []byte) {
	src := iter.ReadStringAsSlice()
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

// ReadNil reads a json object as nil and
// returns whether it's a nil or not
func (iter *Iterator) ReadNil() (ret bool) {
	c := iter.nextToken()
	if c == 'n' {
		iter.skipUntilBreak()
		return true
	}
	iter.unreadByte()
	return false
}

// Skip skips a json object and positions to relatively the next json object
func (iter *Iterator) Skip() {
	c := iter.nextToken()
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
		iter.reportError("Skip", fmt.Sprintf("do not know how to skip: %v", c))
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
			}
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
		if !iter.loadMore() {
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
		if !iter.loadMore() {
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
		if !iter.loadMore() {
			return
		}
	}
}
