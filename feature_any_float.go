package jsoniter

import (
	"io"
	"unsafe"
)

type floatLazyAny struct {
	baseAny
	buf []byte
	iter *Iterator
	err error
	cache float64
}

func (any *floatLazyAny) Parse() *Iterator {
	iter := any.iter
	if iter == nil {
		iter = NewIterator()
	}
	iter.ResetBytes(any.buf)
	return iter
}

func (any *floatLazyAny) fillCache() {
	if any.err != nil {
		return
	}
	iter := any.Parse()
	any.cache = iter.ReadFloat64()
	if iter.Error != io.EOF {
		iter.reportError("floatLazyAny", "there are bytes left")
	}
	any.err = iter.Error
}

func (any *floatLazyAny) LastError() error {
	return any.err
}

func (any *floatLazyAny) ToBool() bool {
	return any.ToFloat64() != 0
}

func (any *floatLazyAny) ToInt() int {
	any.fillCache()
	return int(any.cache)
}

func (any *floatLazyAny) ToInt32() int32 {
	any.fillCache()
	return int32(any.cache)
}

func (any *floatLazyAny) ToInt64() int64 {
	any.fillCache()
	return int64(any.cache)
}

func (any *floatLazyAny) ToFloat32() float32 {
	any.fillCache()
	return float32(any.cache)
}

func (any *floatLazyAny) ToFloat64() float64 {
	any.fillCache()
	return any.cache
}

func (any *floatLazyAny) ToString() string {
	return *(*string)(unsafe.Pointer(&any.buf))
}

func (any *floatLazyAny) WriteTo(stream *Stream) {
	stream.Write(any.buf)
}