package jsoniter

import (
	"strconv"
	"unsafe"
	"io"
)


type uint64LazyAny struct {
	baseAny
	buf   []byte
	iter  *Iterator
	err   error
	cache uint64
}

func (any *uint64LazyAny) ValueType() ValueType {
	return Number
}

func (any *uint64LazyAny) Parse() *Iterator {
	iter := any.iter
	if iter == nil {
		iter = NewIterator()
	}
	iter.ResetBytes(any.buf)
	return iter
}

func (any *uint64LazyAny) fillCache() {
	if any.err != nil {
		return
	}
	iter := any.Parse()
	any.cache = iter.ReadUint64()
	if iter.Error != io.EOF {
		iter.reportError("intLazyAny", "there are bytes left")
	}
	any.err = iter.Error
}

func (any *uint64LazyAny) LastError() error {
	return any.err
}

func (any *uint64LazyAny) ToBool() bool {
	return any.ToInt64() != 0
}

func (any *uint64LazyAny) ToInt() int {
	any.fillCache()
	return int(any.cache)
}

func (any *uint64LazyAny) ToInt32() int32 {
	any.fillCache()
	return int32(any.cache)
}

func (any *uint64LazyAny) ToInt64() int64 {
	any.fillCache()
	return int64(any.cache)
}

func (any *uint64LazyAny) ToUint() uint {
	any.fillCache()
	return uint(any.cache)
}

func (any *uint64LazyAny) ToUint32() uint32 {
	any.fillCache()
	return uint32(any.cache)
}

func (any *uint64LazyAny) ToUint64() uint64 {
	any.fillCache()
	return any.cache
}

func (any *uint64LazyAny) ToFloat32() float32 {
	any.fillCache()
	return float32(any.cache)
}

func (any *uint64LazyAny) ToFloat64() float64 {
	any.fillCache()
	return float64(any.cache)
}

func (any *uint64LazyAny) ToString() string {
	return *(*string)(unsafe.Pointer(&any.buf))
}

func (any *uint64LazyAny) WriteTo(stream *Stream) {
	stream.Write(any.buf)
}

func (any *uint64LazyAny) GetInterface() interface{} {
	any.fillCache()
	return any.cache
}

type uint64Any struct {
	baseAny
	val uint64
}

func (any *uint64Any) LastError() error {
	return nil
}

func (any *uint64Any) ValueType() ValueType {
	return Number
}

func (any *uint64Any) ToBool() bool {
	return any.val != 0
}

func (any *uint64Any) ToInt() int {
	return int(any.val)
}

func (any *uint64Any) ToInt32() int32 {
	return int32(any.val)
}

func (any *uint64Any) ToInt64() int64 {
	return int64(any.val)
}

func (any *uint64Any) ToUint() uint {
	return uint(any.val)
}

func (any *uint64Any) ToUint32() uint32 {
	return uint32(any.val)
}

func (any *uint64Any) ToUint64() uint64 {
	return any.val
}

func (any *uint64Any) ToFloat32() float32 {
	return float32(any.val)
}

func (any *uint64Any) ToFloat64() float64 {
	return float64(any.val)
}

func (any *uint64Any) ToString() string {
	return strconv.FormatUint(any.val, 10)
}

func (any *uint64Any) WriteTo(stream *Stream) {
	stream.WriteUint64(any.val)
}

func (any *uint64Any) Parse() *Iterator {
	return nil
}

func (any *uint64Any) GetInterface() interface{} {
	return any.val
}