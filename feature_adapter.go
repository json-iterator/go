package jsoniter

import (
	"io"
	"unsafe"
	"bytes"
)

// Unmarshal adapts to json/encoding APIs
func Unmarshal(data []byte, v interface{}) error {
	iter := ParseBytes(data)
	iter.Read(v)
	if iter.Error == io.EOF {
		return nil
	}
	return iter.Error
}

func UnmarshalFromString(str string, v interface{}) error {
	// safe to do the unsafe cast here, as str is always referenced in this scope
	data := *(*[]byte)(unsafe.Pointer(&str))
	iter := ParseBytes(data)
	iter.Read(v)
	if iter.Error == io.EOF {
		return nil
	}
	return iter.Error
}

func Marshal(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	stream := NewStream(buf, 4096)
	stream.WriteVal(v)
	if stream.Error != nil {
		return nil, stream.Error
	}
	return buf.Bytes(), nil
}

func MarshalToString(v interface{}) (string, error) {
	buf, err := Marshal(v)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}