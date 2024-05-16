package jsoniter

import (
	"bytes"
	"fmt"
	"io"
)

func ExampleIterator_ReadStringAsReader() {
	json := bytes.NewBufferString(`{"foo":"abcdefghijklmnopqrstuvwxyz"}`)
	iter := Parse(ConfigDefault, json, 16)

	_ = iter.ReadObject()
	r := iter.ReadStringAsReader()
	buf := make([]byte, 8)

	for {
		n, err := r.Read(buf)
		fmt.Println(string(buf[:n]))

		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}
	}

	// Output:
	//
	// abcdefgh
	// ijklmnop
	// qrstuvwx
	// yz
}
