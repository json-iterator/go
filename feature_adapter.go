package jsoniter

import (
	"io"
	"unsafe"
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

func UnmarshalString(str string, v interface{}) error {
	data := *(*[]byte)(unsafe.Pointer(&str))
	iter := ParseBytes(data)
	iter.Read(v)
	if iter.Error == io.EOF {
		return nil
	}
	return iter.Error
}
