package jsoniter

import (
	"strconv"
	"fmt"
)

type stringAny struct {
	baseAny
	val string
}

func (any *stringAny) Get(path ...interface{}) Any {
	if len(path) == 0 {
		return any
	}
	return &invalidAny{baseAny{}, fmt.Errorf("Get %v from simple value", path)}
}

func (any *stringAny) Parse() *Iterator {
	return nil
}

func (any *stringAny) ValueType() ValueType {
	return String
}

func (any *stringAny) MustBeValid() Any {
	return any
}

func (any *stringAny) LastError() error {
	return nil
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
	parsed, _ := strconv.ParseInt(any.val, 10, 64)
	return int(parsed)
}

func (any *stringAny) ToInt32() int32 {
	parsed, _ := strconv.ParseInt(any.val, 10, 32)
	return int32(parsed)
}

func (any *stringAny) ToInt64() int64 {
	parsed, _ := strconv.ParseInt(any.val, 10, 64)
	return parsed
}

func (any *stringAny) ToUint() uint {
	parsed, _ := strconv.ParseUint(any.val, 10, 64)
	return uint(parsed)
}

func (any *stringAny) ToUint32() uint32 {
	parsed, _ := strconv.ParseUint(any.val, 10, 32)
	return uint32(parsed)
}

func (any *stringAny) ToUint64() uint64 {
	parsed, _ := strconv.ParseUint(any.val, 10, 64)
	return parsed
}

func (any *stringAny) ToFloat32() float32 {
	parsed, _ := strconv.ParseFloat(any.val, 32)
	return float32(parsed)
}

func (any *stringAny) ToFloat64() float64 {
	parsed, _ := strconv.ParseFloat(any.val, 64)
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
