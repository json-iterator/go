package jsoniter

import (
	"io"
	"unsafe"
	"strconv"
)

type int64LazyAny struct {
	baseAny
	buf   []byte
	iter  *Iterator
	err   error
	cache int64
}

func (any *int64LazyAny) ValueType() ValueType {
	return Number
}

func (any *int64LazyAny) Parse() *Iterator {
	iter := any.iter
	if iter == nil {
		iter = NewIterator()
	}
	iter.ResetBytes(any.buf)
	return iter
}

func (any *int64LazyAny) fillCache() {
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

func (any *int64LazyAny) LastError() error {
	return any.err
}

func (any *int64LazyAny) ToBool() bool {
	return any.ToInt64() != 0
}

func (any *int64LazyAny) ToInt() int {
	any.fillCache()
	return int(any.cache)
}

func (any *int64LazyAny) ToInt32() int32 {
	any.fillCache()
	return int32(any.cache)
}

func (any *int64LazyAny) ToInt64() int64 {
	any.fillCache()
	return any.cache
}

func (any *int64LazyAny) ToUint() uint {
	any.fillCache()
	return uint(any.cache)
}

func (any *int64LazyAny) ToUint32() uint32 {
	any.fillCache()
	return uint32(any.cache)
}

func (any *int64LazyAny) ToUint64() uint64 {
	any.fillCache()
	return uint64(any.cache)
}

func (any *int64LazyAny) ToFloat32() float32 {
	any.fillCache()
	return float32(any.cache)
}

func (any *int64LazyAny) ToFloat64() float64 {
	any.fillCache()
	return float64(any.cache)
}

func (any *int64LazyAny) ToString() string {
	return *(*string)(unsafe.Pointer(&any.buf))
}

func (any *int64LazyAny) WriteTo(stream *Stream) {
	stream.Write(any.buf)
}

func (any *int64LazyAny) GetInterface() interface{} {
	any.fillCache()
	return any.cache
}

type int64Any struct {
	baseAny
	val int64
}

func (any *int64Any) LastError() error {
	return nil
}

func (any *int64Any) ValueType() ValueType {
	return Number
}

func (any *int64Any) ToBool() bool {
	return any.val != 0
}

func (any *int64Any) ToInt() int {
	return int(any.val)
}

func (any *int64Any) ToInt32() int32 {
	return int32(any.val)
}

func (any *int64Any) ToInt64() int64 {
	return any.val
}

func (any *int64Any) ToUint() uint {
	return uint(any.val)
}

func (any *int64Any) ToUint32() uint32 {
	return uint32(any.val)
}

func (any *int64Any) ToUint64() uint64 {
	return uint64(any.val)
}

func (any *int64Any) ToFloat32() float32 {
	return float32(any.val)
}

func (any *int64Any) ToFloat64() float64 {
	return float64(any.val)
}

func (any *int64Any) ToString() string {
	return strconv.FormatInt(any.val, 10)
}

func (any *int64Any) WriteTo(stream *Stream) {
	stream.WriteInt64(any.val)
}

func (any *int64Any) Parse() *Iterator {
	return nil
}

func (any *int64Any) GetInterface() interface{} {
	return any.val
}