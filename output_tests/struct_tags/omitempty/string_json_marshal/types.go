package test

type typeForTest struct {
	F JM `json:"f,omitempty"`
}

type JM string

func (t *JM) UnmarshalJSON(b []byte) error {
	return nil
}

func (t JM) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}
