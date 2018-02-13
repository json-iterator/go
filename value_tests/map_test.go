package test

func init() {
	marshalCases = append(marshalCases,
		map[string]interface{}{"abc": 1},
		map[string]MyInterface{"hello": MyString("world")},
	)
}

type MyInterface interface {
	Hello() string
}

type MyString string

func (ms MyString) Hello() string {
	return string(ms)
}
