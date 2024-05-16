package jsoniter

import "io"

// Returns a reader that can be used to read the following string until EOF.
// Using the reader after EOF is undefined behavior, and so is using the reader
// when when the current token isn't a string.
func (iter *Iterator) ReadStringAsReader() (r io.Reader) {
	c := iter.nextToken()
	if c == '"' {
		return iterStrReader{iter}
	}
	iter.ReportError("ReadStringAsReader", `expects " or n, but found `+string([]byte{c}))
	return
}

var _ io.Reader = iterStrReader{}

type iterStrReader struct {
	iter *Iterator
}

func (r iterStrReader) Read(dst []byte) (n int, err error) {
	for i := r.iter.head; i < r.iter.tail; i++ {
		// require ascii string and no escape
		// for: field name, base64, number
		if r.iter.buf[i] == '"' {
			n = copy(dst, r.iter.buf[r.iter.head:i])
			r.iter.head += n + 1
			err = io.EOF
			return
		}
	}

	n = copy(dst, r.iter.buf[r.iter.head:])
	r.iter.head += n

	if r.iter.head == r.iter.tail {
		if !r.iter.loadMore() {
			err = io.ErrUnexpectedEOF
		}
	}

	return
}
