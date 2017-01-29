package jsoniter

import (
	"io"
	"strconv"
)

type stringLazyAny struct{
	baseAny
	buf   []byte
	iter  *Iterator
	err   error
	cache string
}

func (any *stringLazyAny) ValueType() ValueType {
	return String
}

func (any *stringLazyAny) Parse() *Iterator {
	iter := any.iter
	if iter == nil {
		iter = NewIterator()
		any.iter = iter
	}
	iter.ResetBytes(any.buf)
	return iter
}

func (any *stringLazyAny) fillCache() {
	if any.err != nil {
		return
	}
	iter := any.Parse()
	any.cache = iter.ReadString()
	if iter.Error != io.EOF {
		iter.reportError("stringLazyAny", "there are bytes left")
	}
	any.err = iter.Error
}

func (any *stringLazyAny) LastError() error {
	return any.err
}

func (any *stringLazyAny) ToBool() bool {
	str := any.ToString()
	if str == "false" {
		return false
	}
	for _, c := range str {
		switch c {
		case ' ', '\n', '\r', '\t':
		default:
			return true
		}
	}
	return false
}

func (any *stringLazyAny) ToInt() int {
	iter := any.Parse()
	iter.head++
	val := iter.ReadInt()
	any.err = iter.Error
	return val
}

func (any *stringLazyAny) ToInt32() int32 {
	iter := any.Parse()
	iter.head++
	val := iter.ReadInt32()
	any.err = iter.Error
	return val
}

func (any *stringLazyAny) ToInt64() int64 {
	iter := any.Parse()
	iter.head++
	val := iter.ReadInt64()
	any.err = iter.Error
	return val
}

func (any *stringLazyAny) ToUint() uint {
	iter := any.Parse()
	iter.head++
	val := iter.ReadUint()
	any.err = iter.Error
	return val
}

func (any *stringLazyAny) ToUint32() uint32 {
	iter := any.Parse()
	iter.head++
	val := iter.ReadUint32()
	any.err = iter.Error
	return val
}

func (any *stringLazyAny) ToUint64() uint64 {
	iter := any.Parse()
	iter.head++
	val := iter.ReadUint64()
	any.err = iter.Error
	return val
}

func (any *stringLazyAny) ToFloat32() float32 {
	iter := any.Parse()
	iter.head++
	val := iter.ReadFloat32()
	any.err = iter.Error
	return val
}

func (any *stringLazyAny) ToFloat64() float64 {
	iter := any.Parse()
	iter.head++
	val := iter.ReadFloat64()
	any.err = iter.Error
	return val
}

func (any *stringLazyAny) ToString() string {
	any.fillCache()
	return any.cache
}

func (any *stringLazyAny) WriteTo(stream *Stream) {
	stream.Write(any.buf)
}

func (any *stringLazyAny) GetInterface() interface{} {
	any.fillCache()
	return any.cache
}

type stringAny struct{
	baseAny
	err   error
	val string
}

func (any *stringAny) Parse() *Iterator {
	return nil
}


func (any *stringAny) ValueType() ValueType {
	return String
}

func (any *stringAny) LastError() error {
	return any.err
}

func (any *stringAny) ToBool() bool {
	str := any.ToString()
	if str == "false" {
		return false
	}
	for _, c := range str {
		switch c {
		case ' ', '\n', '\r', '\t':
		default:
			return true
		}
	}
	return false
}

func (any *stringAny) ToInt() int {
	parsed, err := strconv.ParseInt(any.val, 10, 64)
	any.err = err
	return int(parsed)
}

func (any *stringAny) ToInt32() int32 {
	parsed, err := strconv.ParseInt(any.val, 10, 32)
	any.err = err
	return int32(parsed)
}

func (any *stringAny) ToInt64() int64 {
	parsed, err := strconv.ParseInt(any.val, 10, 64)
	any.err = err
	return parsed
}

func (any *stringAny) ToUint() uint {
	parsed, err := strconv.ParseUint(any.val, 10, 64)
	any.err = err
	return uint(parsed)
}

func (any *stringAny) ToUint32() uint32 {
	parsed, err := strconv.ParseUint(any.val, 10, 32)
	any.err = err
	return uint32(parsed)
}

func (any *stringAny) ToUint64() uint64 {
	parsed, err := strconv.ParseUint(any.val, 10, 64)
	any.err = err
	return parsed
}

func (any *stringAny) ToFloat32() float32 {
	parsed, err := strconv.ParseFloat(any.val, 32)
	any.err = err
	return float32(parsed)
}

func (any *stringAny) ToFloat64() float64 {
	parsed, err := strconv.ParseFloat(any.val, 64)
	any.err = err
	return parsed
}

func (any *stringAny) ToString() string {
	return any.val
}

func (any *stringAny) WriteTo(stream *Stream) {
	stream.WriteString(any.val)
}

func (any *stringAny) GetInterface() interface{} {
	return any.val
}