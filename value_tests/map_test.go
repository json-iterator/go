package test

import "math/big"

func init() {
	marshalCases = append(marshalCases,
		map[string]interface{}{"abc": 1},
		map[string]MyInterface{"hello": MyString("world")},
		map[*big.Float]string{big.NewFloat(1.2): "2"},
		map[string]interface{}{
			"3": 3,
			"1": 1,
			"2": 2,
		},
		map[uint64]interface{}{
			uint64(1): "a",
			uint64(2): "a",
			uint64(4): "a",
		},
	)
	unmarshalCases = append(unmarshalCases, unmarshalCase{
		ptr: (*map[string]string)(nil),
		input: `{"k\"ey": "val"}`,
	})
}

type MyInterface interface {
	Hello() string
}

type MyString string

func (ms MyString) Hello() string {
	return string(ms)
}
