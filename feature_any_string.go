package jsoniter

import (
	"io"
)

type stringLazyAny struct{
	baseAny
	buf   []byte
	iter  *Iterator
	err   error
	cache string
}

func (any *stringLazyAny) fillCache() {
	if any.err != nil {
		return
	}
	iter := any.parse()
	any.cache = iter.ReadString()
	if iter.Error != io.EOF {
		iter.reportError("stringLazyAny", "there are bytes left")
	}
	any.err = iter.Error
}

func (any *stringLazyAny) parse() *Iterator {
	iter := any.iter
	if iter == nil {
		iter = NewIterator()
		any.iter = iter
	}
	iter.ResetBytes(any.buf)
	return iter
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
	iter := any.parse()
	iter.head++
	val := iter.ReadInt()
	any.err = iter.Error
	return val
}

func (any *stringLazyAny) ToInt32() int32 {
	iter := any.parse()
	iter.head++
	val := iter.ReadInt32()
	any.err = iter.Error
	return val
}

func (any *stringLazyAny) ToInt64() int64 {
	iter := any.parse()
	iter.head++
	val := iter.ReadInt64()
	any.err = iter.Error
	return val
}

func (any *stringLazyAny) ToFloat32() float32 {
	iter := any.parse()
	iter.head++
	val := iter.ReadFloat32()
	any.err = iter.Error
	return val
}

func (any *stringLazyAny) ToFloat64() float64 {
	iter := any.parse()
	iter.head++
	val := iter.ReadFloat64()
	any.err = iter.Error
	return val
}

func (any *stringLazyAny) ToString() string {
	any.fillCache()
	return any.cache
}

func (any *stringLazyAny) Get(path ...interface{}) Any {
	return &invalidAny{}
}