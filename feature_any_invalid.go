package jsoniter

type invalidAny struct {
	baseAny
}

func (any *invalidAny) LastError() error {
	return nil
}

func (any *invalidAny) ToBool() bool {
	return false
}

func (any *invalidAny) ToInt() int {
	return 0
}

func (any *invalidAny) ToInt32() int32 {
	return 0
}

func (any *invalidAny) ToInt64() int64 {
	return 0
}

func (any *invalidAny) ToFloat32() float32 {
	return 0
}

func (any *invalidAny) ToFloat64() float64 {
	return 0
}

func (any *invalidAny) ToString() string {
	return ""
}

func (any *invalidAny) WriteTo(stream *Stream) {
}

func (any *invalidAny) Get(path ...interface{}) Any {
	return any
}

func (any *invalidAny) Parse() *Iterator {
	return nil
}
