package jsoniter

type floatLazyAny struct {
	buf []byte
	iter *Iterator
	err error
}

func (any *floatLazyAny) LastError() error {
	return any.err
}

func (any *floatLazyAny) ToBool() bool {
	return false
}

func (any *floatLazyAny) ToInt() int {
	return 0
}

func (any *floatLazyAny) ToInt32() int32 {
	return 0
}

func (any *floatLazyAny) ToInt64() int64 {
	return 0
}

func (any *floatLazyAny) ToFloat32() float32 {
	return 0
}

func (any *floatLazyAny) ToFloat64() float64 {
	return 0
}

func (any *floatLazyAny) ToString() string {
	return ""
}