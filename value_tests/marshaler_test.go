package test

import (
	"encoding/json"
	"encoding"
)

func init() {
	jsonMarshaler := json.Marshaler(fakeJsonMarshaler{})
	textMarshaler := encoding.TextMarshaler(fakeTextMarshaler{})
	marshalCases = append(marshalCases,
		fakeJsonMarshaler{},
		&jsonMarshaler,
		fakeTextMarshaler{},
		&textMarshaler,
	)
}

type fakeJsonMarshaler struct {
	F2 chan []byte
}

func (q fakeJsonMarshaler) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}

func (q *fakeJsonMarshaler) UnmarshalJSON(value []byte) error {
	return nil
}


type fakeTextMarshaler struct {
	F2 chan []byte
}

func (q fakeTextMarshaler) MarshalText() ([]byte, error) {
	return []byte(`""`), nil
}

func (q *fakeTextMarshaler) UnmarshalText(value []byte) error {
	return nil
}
