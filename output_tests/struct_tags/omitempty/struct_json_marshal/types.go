package test

type T struct {
	F JM `json:"f,omitempty"`
}

type JM struct{}

func (t *JM) UnmarshalJSON(b []byte) error {
	return nil
}

func (t JM) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}
