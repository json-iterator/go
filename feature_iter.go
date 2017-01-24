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
	Nil
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
	valueTypes['n'] = Nil
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
		if iter.Error != io.EOF {
			return
		}
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
		if iter.Error == nil {
			iter.Error = io.EOF
		}
		return false
	}
	for {
		n, err := iter.reader.Read(iter.buf)
		if n == 0 {
			if err != nil {
				if iter.Error == nil {
					iter.Error = err
				}
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

func (iter *Iterator) Read() interface{} {
	valueType := iter.WhatIsNext()
	switch valueType {
	case String:
		return iter.ReadString()
	case Number:
		return iter.ReadFloat64()
	case Nil:
		iter.skipFixedBytes(4) // null
		return nil
	case Bool:
		return iter.ReadBool()
	case Array:
		arr := []interface{}{}
		iter.ReadArrayCB(func(iter *Iterator) bool {
			arr = append(arr, iter.Read())
			return true
		})
		return arr
	case Object:
		obj := map[string]interface{}{}
		iter.ReadObjectCB(func(Iter *Iterator, field string) bool {
			obj[field] = iter.Read()
			return true
		})
		return obj
	default:
		iter.reportError("Read", fmt.Sprintf("unexpected value type: %v", valueType))
		return nil
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

