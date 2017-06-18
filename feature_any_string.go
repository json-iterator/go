package jsoniter

import (
	"strconv"
)

type stringAny struct {
	baseAny
	err error
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
