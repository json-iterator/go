package jsoniter

import (
	"io"
)

var bytesNull []byte
var bytesTrue []byte
var bytesFalse []byte

func init() {
	bytesNull = []byte("null")
	bytesTrue = []byte("true")
	bytesFalse = []byte("false")
}

type Stream struct {
	out           io.Writer
	buf           []byte
	n             int
	Error         error
	indention     int
	IndentionStep int
}

func NewStream(out io.Writer, bufSize int) *Stream {
	return &Stream{out, make([]byte, bufSize), 0, nil, 0, 0}
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
func (b *Stream) writeByte(c byte) error {
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

func (b *Stream) WriteRaw(s string) {
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

func (b *Stream) WriteString(s string) {
	b.writeByte('"')
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
	b.writeByte('"')
}

func (stream *Stream) WriteNull() {
	stream.Write(bytesNull)
}

func (stream *Stream) WriteTrue() {
	stream.Write(bytesTrue)
}

func (stream *Stream) WriteFalse() {
	stream.Write(bytesFalse)
}

func (stream *Stream) WriteBool(val bool) {
	if val {
		stream.Write(bytesTrue)
	} else {
		stream.Write(bytesFalse)
	}
}

func (stream *Stream) WriteObjectStart() {
	stream.indention += stream.IndentionStep
	stream.writeByte('{')
	stream.writeIndention(0)
}

func (stream *Stream) WriteObjectField(field string) {
	stream.WriteString(field)
	stream.writeByte(':')
}

func (stream *Stream) WriteObjectEnd() {
	stream.writeIndention(stream.IndentionStep)
	stream.indention -= stream.IndentionStep
	stream.writeByte('}')
}

func (stream *Stream) WriteMore() {
	stream.writeByte(',')
	stream.writeIndention(0)
}

func (stream *Stream) WriteArrayStart() {
	stream.indention += stream.IndentionStep
	stream.writeByte('[')
	stream.writeIndention(0)
}

func (stream *Stream) WriteArrayEnd() {
	stream.writeIndention(stream.IndentionStep)
	stream.indention -= stream.IndentionStep
	stream.writeByte(']')
}

func (stream *Stream) writeIndention(delta int) {
	if (stream.indention == 0) {
		return
	}
	stream.writeByte('\n')
	toWrite := stream.indention - delta
	i := 0
	for {
		for ; i < toWrite && stream.n < len(stream.buf); i++ {
			stream.buf[stream.n] = ' '
			stream.n ++
		}
		if i == toWrite {
			break;
		} else {
			stream.Flush()
		}
	}
}

func (stream *Stream) WriteVal(val interface{}) {
}