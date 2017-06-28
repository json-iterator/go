package test

import (
	"encoding"
	"strings"
)

type KeyType struct {
	X string
}

func (k KeyType) MarshalText() ([]byte, error) {
	return []byte("MANUAL__" + k.X), nil
}

func (k *KeyType) UnmarshalText(text []byte) error {
	k.X = strings.TrimPrefix(string(text), "MANUAL__")
	return nil
}

var _ encoding.TextMarshaler = KeyType{}
var _ encoding.TextUnmarshaler = &KeyType{}

type T map[KeyType]string
