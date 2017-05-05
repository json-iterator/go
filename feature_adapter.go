package jsoniter

import (
	"io"
	"bytes"
)

// Unmarshal adapts to json/encoding APIs
func Unmarshal(data []byte, v interface{}) error {
	data = data[:lastNotSpacePos(data)]
	iter := ParseBytes(data)
	iter.ReadVal(v)
	if iter.head == iter.tail {
		iter.loadMore()
	}
	if iter.Error == io.EOF {
		return nil
	}
	if iter.Error == nil {
		iter.reportError("Unmarshal", "there are bytes left after unmarshal")
	}
	return iter.Error
}

func UnmarshalAny(data []byte) (Any, error) {
	data = data[:lastNotSpacePos(data)]
	iter := ParseBytes(data)
	any := iter.ReadAny()
	if iter.head == iter.tail {
		iter.loadMore()
	}
	if iter.Error == io.EOF {
		return any, nil
	}
	if iter.Error == nil {
		iter.reportError("UnmarshalAny", "there are bytes left after unmarshal")
	}
	return any, iter.Error
}

func lastNotSpacePos(data []byte) int {
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] != ' ' && data[i] != '\t' && data[i] != '\r' && data[i] != '\n' {
			return i + 1
		}
	}
	return 0
}

func UnmarshalFromString(str string, v interface{}) error {
	data := []byte(str)
	data = data[:lastNotSpacePos(data)]
	iter := ParseBytes(data)
	iter.ReadVal(v)
	if iter.head == iter.tail {
		iter.loadMore()
	}
	if iter.Error == io.EOF {
		return nil
	}
	if iter.Error == nil {
		iter.reportError("UnmarshalFromString", "there are bytes left after unmarshal")
	}
	return iter.Error
}

func UnmarshalAnyFromString(str string) (Any, error) {
	data := []byte(str)
	data = data[:lastNotSpacePos(data)]
	iter := ParseBytes(data)
	any := iter.ReadAny()
	if iter.head == iter.tail {
		iter.loadMore()
	}
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
