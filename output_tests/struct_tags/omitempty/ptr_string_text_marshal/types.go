package test

type T struct {
	F *TM `json:"f,omitempty"`
}

type TM string

func (t *TM) UnmarshalText(b []byte) error {
	return nil
}

func (t TM) MarshalText() ([]byte, error) {
	return []byte(`""`), nil
}
