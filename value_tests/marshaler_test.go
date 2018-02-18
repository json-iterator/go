package test

import (
	"encoding/json"
	"encoding"
)

func init() {
	jsonMarshaler := json.Marshaler(fakeJsonMarshaler{})
	textMarshaler := encoding.TextMarshaler(fakeTextMarshaler{})
	textMarshaler2 := encoding.TextMarshaler(&fakeTextMarshaler2{})
	marshalCases = append(marshalCases,
		fakeJsonMarshaler{},
		&jsonMarshaler,
		fakeTextMarshaler{},
		&textMarshaler,
		fakeTextMarshaler2{},
		&textMarshaler2,
		map[fakeTextMarshaler]int{
			fakeTextMarshaler{}: 100,
		},
		map[*fakeTextMarshaler]int{
			&fakeTextMarshaler{}: 100,
		},
		map[encoding.TextMarshaler]int{
			textMarshaler: 100,
		},
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

type fakeTextMarshaler2 struct {
	Field2 int
}

func (q *fakeTextMarshaler2) MarshalText() ([]byte, error) {
	return []byte(`"abc"`), nil
}

func (q *fakeTextMarshaler2) UnmarshalText(value []byte) error {
	return nil
}
