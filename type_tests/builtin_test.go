package test

func init() {
	testCases = append(testCases,
		(*bool)(nil),
		(*boolAlias)(nil),
		(*byte)(nil),
		(*byteAlias)(nil),
		(*float32)(nil),
		(*float32Alias)(nil),
		(*float64)(nil),
		(*float64Alias)(nil),
		(*int8)(nil),
		(*int8Alias)(nil),
		(*int16)(nil),
		(*int16Alias)(nil),
		(*int32)(nil),
		(*int32Alias)(nil),
		(*int64)(nil),
		(*int64Alias)(nil),
		(*string)(nil),
		(*stringAlias)(nil),
		(*uint8)(nil),
		(*uint8Alias)(nil),
		(*uint16)(nil),
		(*uint16Alias)(nil),
		(*uint32)(nil),
		(*uint32Alias)(nil),
		(*uintptr)(nil),
		(*uintptrAlias)(nil),
	)
}

type boolAlias bool
type byteAlias byte
type float32Alias float32
type float64Alias float64
type ptrFloat64Alias *float64
type int8Alias int8
type int16Alias int16
type int32Alias int32
type ptrInt32Alias *int32
type int64Alias int64
type stringAlias string
type ptrStringAlias *string
type uint8Alias uint8
type uint16Alias uint16
type uint32Alias uint32
type uintptrAlias uintptr