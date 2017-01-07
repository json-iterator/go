package jsoniter

import (
	"io"
)

var bytesNull []byte
var digits []uint8;

func init() {
	bytesNull = []byte("null")
	digits = []uint8{
		'0', '1', '2', '3', '4', '5',
		'6', '7', '8', '9', 'a', 'b',
		'c', 'd', 'e', 'f', 'g', 'h',
		'i', 'j', 'k', 'l', 'm', 'n',
		'o', 'p', 'q', 'r', 's', 't',
		'u', 'v', 'w', 'x', 'y', 'z',
	}
}

type Stream struct {
	out   io.Writer
	buf   []byte
	n     int
	Error error
}

func NewStream(out io.Writer, bufSize int) *Stream {
	return &Stream{out, make([]byte, bufSize), 0, nil}
}


// Available returns how many bytes are unused in the buffer.
func (b *Stream) Available() int {
	return len(b.buf) - b.n
}

// Buffered returns the number of bytes that have been written into the current buffer.
func (b *Stream) Buffered() int {
	return b.n
}

// Write writes the contents of p into the buffer.
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (b *Stream) Write(p []byte) (nn int, err error) {
	for len(p) > b.Available() && b.Error == nil {
		var n int
		if b.Buffered() == 0 {
			// Large write, empty buffer.
			// Write directly from p to avoid copy.
			n, b.Error = b.out.Write(p)
		} else {
			n = copy(b.buf[b.n:], p)
			b.n += n
			b.Flush()
		}
		nn += n
		p = p[n:]
	}
	if b.Error != nil {
		return nn, b.Error
	}
	n := copy(b.buf[b.n:], p)
	b.n += n
	nn += n
	return nn, nil
}


// WriteByte writes a single byte.
func (b *Stream) WriteByte(c byte) error {
	if b.Error != nil {
		return b.Error
	}
	if b.Available() <= 0 && b.Flush() != nil {
		return b.Error
	}
	b.buf[b.n] = c
	b.n++
	return nil
}

// Flush writes any buffered data to the underlying io.Writer.
func (b *Stream) Flush() error {
	if b.Error != nil {
		return b.Error
	}
	if b.n == 0 {
		return nil
	}
	n, err := b.out.Write(b.buf[0:b.n])
	if n < b.n && err == nil {
		err = io.ErrShortWrite
	}
	if err != nil {
		if n > 0 && n < b.n {
			copy(b.buf[0:b.n - n], b.buf[n:b.n])
		}
		b.n -= n
		b.Error = err
		return err
	}
	b.n = 0
	return nil
}

func (b *Stream) WriteString(s string) {
	for len(s) > b.Available() && b.Error == nil {
		n := copy(b.buf[b.n:], s)
		b.n += n
		s = s[n:]
		b.Flush()
	}
	if b.Error != nil {
		return
	}
	n := copy(b.buf[b.n:], s)
	b.n += n
}

func (stream *Stream) WriteNull() {
	stream.Write(bytesNull)
}

func (stream *Stream) WriteUint8(val uint8) {
	if stream.Available() < 3 {
		stream.Flush()
	}
	charPos := stream.n
	if val <= 9 {
		charPos += 1;
	} else {
		if val <= 99 {
			charPos += 2;
		} else {
			charPos += 3;
		}
	}
	stream.n = charPos
	var q uint8
	var r uint8
	for {
		q = val / 10
		r = val - ((q << 3) + (q << 1))  // r = i-(q*10) ...
		charPos--
		stream.buf[charPos] = digits[r]
		val = q;
		if val == 0 {
			break
		}
	}
}

func (stream *Stream) WriteInt8(val int8) {
	if stream.Available() < 4 {
		stream.Flush()
	}
	charPos := stream.n
	if (val < 0) {
		charPos += 1
		val = -val
		stream.buf[stream.n] = '-'
	}
	if val <= 9 {
		charPos += 1;
	} else {
		if val <= 99 {
			charPos += 2;
		} else {
			charPos += 3;
		}
	}
	stream.n = charPos
	var q int8
	var r int8
	for {
		q = val / 10
		r = val - ((q << 3) + (q << 1))  // r = i-(q*10) ...
		charPos--
		stream.buf[charPos] = digits[r]
		val = q;
		if val == 0 {
			break
		}
	}
}

func (stream *Stream) WriteUint16(val uint16) {
	if stream.Available() < 5 {
		stream.Flush()
	}
	charPos := stream.n
	if val <= 99 {
		if val <= 9 {
			charPos += 1;
		} else {
			charPos += 2;
		}
	} else {
		if val <= 999 {
			charPos += 3;
		} else {
			if val <= 9999 {
				charPos += 4;
			} else {
				charPos += 5;
			}
		}
	}
	stream.n = charPos
	var q uint16
	var r uint16
	for {
		q = val / 10
		r = val - ((q << 3) + (q << 1))  // r = i-(q*10) ...
		charPos--
		stream.buf[charPos] = digits[r]
		val = q;
		if val == 0 {
			break
		}
	}
}

func (stream *Stream) WriteVal(val interface{}) {
}