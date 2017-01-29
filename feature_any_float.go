package jsoniter

import (
	"io"
	"unsafe"
	"strconv"
)

type float64LazyAny struct {
	baseAny
	buf []byte
	iter *Iterator
	err error
	cache float64
}

func (any *float64LazyAny) Parse() *Iterator {
	iter := any.iter
	if iter == nil {
		iter = NewIterator()
	}
	iter.ResetBytes(any.buf)
	return iter
}

func (any *float64LazyAny) ValueType() ValueType {
	return Number
}

func (any *float64LazyAny) fillCache() {
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

func (any *float64LazyAny) LastError() error {
	return any.err
}

func (any *float64LazyAny) ToBool() bool {
	return any.ToFloat64() != 0
}

func (any *float64LazyAny) ToInt() int {
	any.fillCache()
	return int(any.cache)
}

func (any *float64LazyAny) ToInt32() int32 {
	any.fillCache()
	return int32(any.cache)
}

func (any *float64LazyAny) ToInt64() int64 {
	any.fillCache()
	return int64(any.cache)
}

func (any *float64LazyAny) ToUint() uint {
	any.fillCache()
	return uint(any.cache)
}

func (any *float64LazyAny) ToUint32() uint32 {
	any.fillCache()
	return uint32(any.cache)
}

func (any *float64LazyAny) ToUint64() uint64 {
	any.fillCache()
	return uint64(any.cache)
}

func (any *float64LazyAny) ToFloat32() float32 {
	any.fillCache()
	return float32(any.cache)
}

func (any *float64LazyAny) ToFloat64() float64 {
	any.fillCache()
	return any.cache
}

func (any *float64LazyAny) ToString() string {
	return *(*string)(unsafe.Pointer(&any.buf))
}

func (any *float64LazyAny) WriteTo(stream *Stream) {
	stream.Write(any.buf)
}

func (any *float64LazyAny) GetInterface() interface{} {
	any.fillCache()
	return any.cache
}

type floatAny struct {
	baseAny
	val float64
}

func (any *floatAny) Parse() *Iterator {
	return nil
}

func (any *floatAny) ValueType() ValueType {
	return Number
}

func (any *floatAny) LastError() error {
	return nil
}

func (any *floatAny) ToBool() bool {
	return any.ToFloat64() != 0
}

func (any *floatAny) ToInt() int {
	return int(any.val)
}

func (any *floatAny) ToInt32() int32 {
	return int32(any.val)
}

func (any *floatAny) ToInt64() int64 {
	return int64(any.val)
}

func (any *floatAny) ToUint() uint {
	return uint(any.val)
}

func (any *floatAny) ToUint32() uint32 {
	return uint32(any.val)
}

func (any *floatAny) ToUint64() uint64 {
	return uint64(any.val)
}

func (any *floatAny) ToFloat32() float32 {
	return float32(any.val)
}

func (any *floatAny) ToFloat64() float64 {
	return any.val
}

func (any *floatAny) ToString() string {
	return strconv.FormatFloat(any.val, 'E', -1, 64)
}

func (any *floatAny) WriteTo(stream *Stream) {
	stream.WriteFloat64(any.val)
}

func (any *floatAny) GetInterface() interface{} {
	return any.val
}