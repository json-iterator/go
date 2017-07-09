package test

import (
	"encoding"
	"strings"
)

type KeyType string

func (k KeyType) MarshalText() ([]byte, error) {
	return []byte("MANUAL__" + k), nil
}

func (k *KeyType) UnmarshalText(text []byte) error {
	*k = KeyType(strings.TrimPrefix(string(text), "MANUAL__"))
	return nil
}

var _ encoding.TextMarshaler = KeyType("")
var _ encoding.TextUnmarshaler = new(KeyType)

type typeForTest map[KeyType]string
