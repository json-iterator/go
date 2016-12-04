package jsoniter

// adapt to json/encoding api

func Unmarshal(data []byte, v interface{}) error {
	iter := ParseBytes(data)
	iter.Read(v)
	return iter.Error
}
