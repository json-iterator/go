package jsoniter

import (
	"io"
	"unsafe"
)

type intLazyAny struct {
	buf   []byte
	iter  *Iterator
	err   error
	cache int64
}

func (any *intLazyAny) fillCache() {
	if any.err != nil {
		return
	}
	iter := any.iter
	if iter == nil {
		iter = NewIterator()
	}
	iter.ResetBytes(any.buf)
	any.cache = iter.ReadInt64()
	if iter.Error != io.EOF {
		iter.reportError("intLazyAny", "there are bytes left")
	}
	any.err = iter.Error
}

func (any *intLazyAny) LastError() error {
	return any.err
}

func (any *intLazyAny) ToBool() bool {
	return any.ToInt64() != 0
}

func (any *intLazyAny) ToInt() int {
	any.fillCache()
	return int(any.cache)
}

func (any *intLazyAny) ToInt32() int32 {
	any.fillCache()
	return int32(any.cache)
}

func (any *intLazyAny) ToInt64() int64 {
	any.fillCache()
	return any.cache
}

func (any *intLazyAny) ToFloat32() float32 {
	any.fillCache()
	return float32(any.cache)
}

func (any *intLazyAny) ToFloat64() float64 {
	any.fillCache()
	return float64(any.cache)
}

func (any *intLazyAny) ToString() string {
	return *(*string)(unsafe.Pointer(&any.buf))
}