package test

func init() {
	two := float64(2)
	marshalCases = append(marshalCases,
		[1]*float64{nil},
		[1]*float64{&two},
		[2]*float64{},
	)
}
