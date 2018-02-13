package test

func init() {
	var pEFace = func(val interface{}) *interface{} {
		return &val
	}
	unmarshalCases = append(unmarshalCases, unmarshalCase{
		ptr: (**interface{})(nil),
		input: `"hello"`,
	}, unmarshalCase{
		ptr: (**interface{})(nil),
		input: `1e1`,
	}, unmarshalCase{
		ptr: (**interface{})(nil),
		input: `1.0e1`,
	})
	marshalCases = append(marshalCases,
		pEFace("hello"),
	)
}