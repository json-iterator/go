package jsoniter

type trueAny struct {
}

func (any *trueAny) LastError() error {
	return nil
}

func (any *trueAny) ToBool() bool {
	return true
}

func (any *trueAny) ToInt() int {
	return 1
}

func (any *trueAny) ToInt32() int32 {
	return 1
}

func (any *trueAny) ToInt64() int64 {
	return 1
}

func (any *trueAny) ToFloat32() float32 {
	return 1
}

func (any *trueAny) ToFloat64() float64 {
	return 1
}

func (any *trueAny) ToString() string {
	return "true"
}

type falseAny struct {
}

func (any *falseAny) LastError() error {
	return nil
}

func (any *falseAny) ToBool() bool {
	return false
}

func (any *falseAny) ToInt() int {
	return 0
}

func (any *falseAny) ToInt32() int32 {
	return 0
}

func (any *falseAny) ToInt64() int64 {
	return 0
}

func (any *falseAny) ToFloat32() float32 {
	return 0
}

func (any *falseAny) ToFloat64() float64 {
	return 0
}

func (any *falseAny) ToString() string {
	return "false"
}
