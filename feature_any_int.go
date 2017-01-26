package jsoniter

import (
	"io"
	"unsafe"
	"strconv"
)

type intLazyAny struct {
	baseAny
	buf   []byte
	iter  *Iterator
	err   error
	cache int64
}

func (any *intLazyAny) Parse() *Iterator {
	iter := any.iter
	if iter == nil {
		iter = NewIterator()
	}
	iter.ResetBytes(any.buf)
	return iter
}

func (any *intLazyAny) fillCache() {
	if any.err != nil {
		return
	}
	iter := any.Parse()
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

func (any *intLazyAny) WriteTo(stream *Stream) {
	stream.Write(any.buf)
}

type intAny struct {
	baseAny
	err   error
	val int64
}

func (any *intAny) LastError() error {
	return any.err
}

func (any *intAny) ToBool() bool {
	return any.ToInt64() != 0
}

func (any *intAny) ToInt() int {
	return int(any.val)
}

func (any *intAny) ToInt32() int32 {
	return int32(any.val)
}

func (any *intAny) ToInt64() int64 {
	return any.val
}

func (any *intAny) ToFloat32() float32 {
	return float32(any.val)
}

func (any *intAny) ToFloat64() float64 {
	return float64(any.val)
}

func (any *intAny) ToString() string {
	return strconv.FormatInt(any.val, 10)
}

func (any *intAny) WriteTo(stream *Stream) {
	stream.WriteInt64(any.val)
}

func (any *intAny) Parse() *Iterator {
	return nil
}