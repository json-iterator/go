package jsoniter

import "io"

// Unmarshal adapts to json/encoding APIs
func Unmarshal(data []byte, v interface{}) error {
	iter := ParseBytes(data)
	iter.Read(v)
	if iter.Error == io.EOF {
		return nil
	}
	return iter.Error
}
