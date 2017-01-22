package jsoniter

import (
	"io"
	"unsafe"
	"bytes"
)

// Unmarshal adapts to json/encoding APIs
func Unmarshal(data []byte, v interface{}) error {
	iter := ParseBytes(data)
	iter.ReadVal(v)
	if iter.Error == io.EOF {
		return nil
	}
	if iter.Error == nil {
		iter.reportError("UnmarshalAny", "there are bytes left after unmarshal")
	}
	return iter.Error
}

func UnmarshalAny(data []byte) (Any, error) {
	iter := ParseBytes(data)
	any := iter.ReadAny()
	if iter.Error == io.EOF {
		return any, nil
	}
	if iter.Error == nil {
		iter.reportError("UnmarshalAny", "there are bytes left after unmarshal")
	}
	return any, iter.Error
}

func UnmarshalFromString(str string, v interface{}) error {
	// safe to do the unsafe cast here, as str is always referenced in this scope
	data := *(*[]byte)(unsafe.Pointer(&str))
	iter := ParseBytes(data)
	iter.ReadVal(v)
	if iter.Error == io.EOF {
		return nil
	}
	if iter.Error == nil {
		iter.reportError("UnmarshalFromString", "there are bytes left after unmarshal")
	}
	return iter.Error
}

func UnmarshalAnyFromString(str string) (Any, error) {
	// safe to do the unsafe cast here, as str is always referenced in this scope
	data := *(*[]byte)(unsafe.Pointer(&str))
	iter := ParseBytes(data)
	any := iter.ReadAny()
	if iter.Error == io.EOF {
		return any, nil
	}
	if iter.Error == nil {
		iter.reportError("UnmarshalAnyFromString", "there are bytes left after unmarshal")
	}
	return nil, iter.Error
}

func Marshal(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	stream := NewStream(buf, 4096)
	stream.WriteVal(v)
	stream.Flush()
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