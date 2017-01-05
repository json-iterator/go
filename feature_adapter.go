package jsoniter

import "io"

// adapt to json/encoding api

func Unmarshal(data []byte, v interface{}) error {
	iter := ParseBytes(data)
	iter.Read(v)
	if iter.Error == io.EOF {
		return nil
	}
	return iter.Error
}
