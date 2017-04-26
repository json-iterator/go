package jsoniter

import (
	"errors"
)

var (
	// ErrInvalidPath represent a invalid path
	ErrInvalidPath = errors.New("Invalid Path")

	// ErrInvalidKeyReader represent not set a correct reader
	ErrInvalidKeyReader = errors.New("Invalid Key Reader")
)

var (
	_ KeyReader = &KeyString{}
	_ KeyReader = &KeyInt64{}
	_ KeyReader = &KeyBool{}
	_ KeyReader = &KeyFloat64{}
)

// KeyReader represent a read to read the key
type KeyReader interface {
	Read(iter *Iterator, key IKey) (bool, error)
	HasRead() bool
	MaybeNull() bool
}

// KeyString --------------------------------
type KeyString struct {
	Key     IKey
	MayNull bool
	hasRead bool
	String  string
}

// Read implement KeyReader
func (k *KeyString) Read(iter *Iterator, key IKey) (bool, error) {
	if !k.hasRead && k.Key.Equal(key) {
		k.String = iter.ReadString()
		if k.String == "" && iter.Error != nil {
			return false, ErrInvalidKeyReader
		}
		k.hasRead = true
		return true, nil
	}
	return false, nil
}

// HasRead implement the KeyReader
func (k *KeyString) HasRead() bool {
	return k.hasRead
}

// MaybeNull implement the KeyReader
func (k *KeyString) MaybeNull() bool {
	return k.MayNull
}

// KeyInt64 ----------------------------------------------------------
type KeyInt64 struct {
	Key     IKey
	MayNull bool
	hasRead bool
	Int64   int64
}

func (ki *KeyInt64) Read(iter *Iterator, key IKey) (bool, error) {
	if !ki.hasRead && ki.Key.Equal(key) {
		ki.hasRead = true
		ki.Int64 = iter.ReadInt64()
		if ki.Int64 == 0 && iter.Error != nil {
			return false, ErrInvalidKeyReader
		}
		return true, nil
	}
	return false, nil
}

// HasRead implement the KeyReader
func (ki *KeyInt64) HasRead() bool {
	return ki.hasRead
}

// MaybeNull implement the KeyReader
func (ki *KeyInt64) MaybeNull() bool {
	return ki.MayNull
}

// KeyFloat64 -------------------------------------------------------
type KeyFloat64 struct {
	Key     IKey
	MayNull bool
	hasRead bool
	Float64 float64
}

func (kf *KeyFloat64) Read(iter *Iterator, key IKey) (bool, error) {
	if !kf.hasRead && kf.Key.Equal(key) {
		kf.hasRead = true
		kf.Float64 = iter.ReadFloat64()
		if kf.Float64 == 0 && iter.Error != nil {
			return false, ErrInvalidKeyReader
		}
		return true, nil
	}

	return false, nil
}

// HasRead implement the KeyReader
func (kf *KeyFloat64) HasRead() bool {
	return kf.hasRead
}

// MaybeNull implement the KeyReader
func (kf *KeyFloat64) MaybeNull() bool {
	return kf.MayNull
}

// KeyBool ---------------------------------------------------------
type KeyBool struct {
	Key     IKey
	MayNull bool
	hasRead bool
	Bool    bool
}

func (kb *KeyBool) Read(iter *Iterator, key IKey) (bool, error) {
	if !kb.hasRead && kb.Key.Equal(key) {
		kb.hasRead = true
		kb.Bool = iter.ReadBool()
		if kb.Bool == false && iter.Error != nil {
			return false, ErrInvalidKeyReader
		}
		return true, nil
	}
	return false, nil
}

// HasRead implement the KeyReader
func (kb *KeyBool) HasRead() bool {
	return kb.hasRead
}

// MaybeNull implement the KeyReader
func (kb *KeyBool) MaybeNull() bool {
	return kb.MayNull
}
