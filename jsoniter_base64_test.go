package jsoniter

import (
	"testing"
	"bytes"
)

func Test_read_base64(t *testing.T) {
	iter := ParseString(`"YWJj"`)
	if !bytes.Equal(iter.ReadBase64(), []byte("abc")) {
		t.FailNow()
	}
}
