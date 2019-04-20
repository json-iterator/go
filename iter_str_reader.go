package jsoniter

import "io"

type stringReader struct {
	eof  bool
	iter *Iterator
}

// Read implements io.Reader. There is no provision for resetting the reader
// after EOF, nor for ensuring that the full JSON string is read to its end.
func (r *stringReader) Read(p []byte) (n int, err error) {
	if r.eof {
		return 0, io.EOF
	}
	p = p[:0]
	iter := r.iter
	for iter.Error == nil && len(p) < cap(p) {
		c := iter.readByte()
		if c == '"' {
			r.eof = true
			break
		}
		if c == '\\' {
			// if we don't have room for 2 runes, leave it for the next call
			if len(p)+8 > cap(p) {
				iter.unreadByte()
				break
			}
			c = iter.readByte()
			p = iter.readEscapedChar(c, p)
		} else {
			p = append(p, c)
		}
	}
	return len(p), nil
}

// StringReader provides an io.Reader for the upcoming JSON string value. The
// consumer absolutely MUST consume the entire string in order to preserve the
// state of the Iterator. If the next value in the JSON stream is not a string,
// the returned reader will be nil.
func (iter *Iterator) StringReader() io.Reader {
	c := iter.nextToken()
	if c == '"' {
		return &stringReader{iter: iter}
	} else if c == 'n' {
		iter.skipThreeBytes('u', 'l', 'l')
		return nil
	}
	iter.ReportError("StringReader", `expects " or n, but found `+string([]byte{c}))
	return nil
}
