package test

func init() {
	marshalCases = append(marshalCases,
		withChan{},
	)
}

type withChan struct {
	F2 chan []byte
}

func (q withChan) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}

func (q *withChan) UnmarshalJSON(value []byte) error {
	return nil
}
