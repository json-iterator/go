package jsoniter

import (
	"io"
)

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

func (b *Stream) Reset(out io.Writer) {
	b.out = out
	b.n = 0
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
func (b *Stream) writeByte(c byte) {
	if b.Error != nil {
		return
	}
	if b.Available() <= 0 && b.Flush() != nil {
		return
	}
	b.buf[b.n] = c
	b.n++
}

func (b *Stream) writeTwoBytes(c1 byte, c2 byte) {
	if b.Error != nil {
		return
	}
	if b.Available() <= 1 && b.Flush() != nil {
		return
	}
	b.buf[b.n] = c1
	b.buf[b.n + 1] = c2
	b.n += 2
}

func (b *Stream) writeThreeBytes(c1 byte, c2 byte, c3 byte) {
	if b.Error != nil {
		return
	}
	if b.Available() <= 2 && b.Flush() != nil {
		return
	}
	b.buf[b.n] = c1
	b.buf[b.n + 1] = c2
	b.buf[b.n + 2] = c3
	b.n += 3
}

func (b *Stream) writeFourBytes(c1 byte, c2 byte, c3 byte, c4 byte) {
	if b.Error != nil {
		return
	}
	if b.Available() <= 3 && b.Flush() != nil {
		return
	}
	b.buf[b.n] = c1
	b.buf[b.n + 1] = c2
	b.buf[b.n + 2] = c3
	b.buf[b.n + 3] = c4
	b.n += 4
}

func (b *Stream) writeFiveBytes(c1 byte, c2 byte, c3 byte, c4 byte, c5 byte) {
	if b.Error != nil {
		return
	}
	if b.Available() <= 3 && b.Flush() != nil {
		return
	}
	b.buf[b.n] = c1
	b.buf[b.n + 1] = c2
	b.buf[b.n + 2] = c3
	b.buf[b.n + 3] = c4
	b.buf[b.n + 4] = c5
	b.n += 5
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

func (stream *Stream) WriteString(s string) {
	valLen := len(s)
	toWriteLen := valLen
	bufLengthMinusTwo := len(stream.buf) - 2 // make room for the quotes
	if stream.n + toWriteLen > bufLengthMinusTwo {
		toWriteLen = bufLengthMinusTwo - stream.n
	}
	if toWriteLen < 0 {
		stream.Flush()
		if stream.n + toWriteLen > bufLengthMinusTwo {
			toWriteLen = bufLengthMinusTwo - stream.n
		}
	}
	n := stream.n
	stream.buf[n] = '"'
	n++
	// write string, the fast path, without utf8 and escape support
	i := 0
	for ; i < toWriteLen; i++ {
		c := s[i]
		if c > 31 && c != '"' && c != '\\' {
			stream.buf[n] = c
			n++
		} else {
			break;
		}
	}
	if i == valLen {
		stream.buf[n] = '"'
		n++
		stream.n = n
		return
	}
	stream.n = n
	// for the remaining parts, we process them char by char
	stream.writeStringSlowPath(s, i, valLen);
	stream.writeByte('"')
}

func (stream *Stream) writeStringSlowPath(s string, i int, valLen int) {
	for ; i < valLen; i++ {
		c := s[i]
		switch (c) {
		case '"':
			stream.writeTwoBytes('\\', '"')
		case '\\':
			stream.writeTwoBytes('\\', '\\')
		case '\b':
			stream.writeTwoBytes('\\', 'b')
		case '\f':
			stream.writeTwoBytes('\\', 'f')
		case '\n':
			stream.writeTwoBytes('\\', 'n')
		case '\r':
			stream.writeTwoBytes('\\', 'r')
		case '\t':
			stream.writeTwoBytes('\\', 't')
		default:
			stream.writeByte(c);
		}
	}
}

func (stream *Stream) WriteNil() {
	stream.writeFourBytes('n', 'u', 'l', 'l')
}

func (stream *Stream) WriteTrue() {
	stream.writeFourBytes('t', 'r', 'u', 'e')
}

func (stream *Stream) WriteFalse() {
	stream.writeFiveBytes('f', 'a', 'l', 's', 'e')
}

func (stream *Stream) WriteBool(val bool) {
	if val {
		stream.WriteTrue()
	} else {
		stream.WriteFalse()
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

func (stream *Stream) WriteEmptyObject() {
	stream.writeByte('{')
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

func (stream *Stream) WriteEmptyArray() {
	stream.writeByte('[')
	stream.writeByte(']')
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