package jsoniter

import (
	"fmt"
	"io"
)

// stream is a io.Writer like object, with JSON specific write functions.
// Error is not returned as return value, but stored as Error member on this stream instance.
type Stream struct {
	cfg                         *frozenConfig
	out                         io.Writer
	buf                         []byte
	Error                       error
	indention                   int
	Attachment                  interface{} // open for customized encoder
	enforceMarshalledBytesLimit bool

	// Number of bytes remaining before marshalled size exceeds cfg.maxMarshalledBytes.
	// This is tracked as an amount remaining to account for bytes already flushed in Write().
	marshalledBytesLimitRemaining uint64
}

// NewStream create new stream instance.
// cfg can be jsoniter.ConfigDefault.
// out can be nil if write to internal buffer.
// bufSize is the initial size for the internal buffer in bytes.
func NewStream(cfg API, out io.Writer, bufSize int) *Stream {
	config := cfg.(*frozenConfig)
	return &Stream{
		cfg:                           config,
		out:                           out,
		buf:                           make([]byte, 0, bufSize),
		Error:                         nil,
		indention:                     0,
		enforceMarshalledBytesLimit:   config.maxMarshalledBytes > 0,
		marshalledBytesLimitRemaining: config.maxMarshalledBytes,
	}
}

// Pool returns a pool can provide more stream with same configuration
func (stream *Stream) Pool() StreamPool {
	return stream.cfg
}

// Reset reuse this stream instance by assign a new writer
func (stream *Stream) Reset(out io.Writer) {
	stream.out = out
	stream.buf = stream.buf[:0]
	stream.marshalledBytesLimitRemaining = stream.cfg.maxMarshalledBytes
}

// Available returns how many bytes are unused in the buffer.
func (stream *Stream) Available() int {
	return cap(stream.buf) - len(stream.buf)
}

// Buffered returns the number of bytes that have been written into the current buffer.
func (stream *Stream) Buffered() int {
	return len(stream.buf)
}

// Buffer if writer is nil, use this method to take the result
func (stream *Stream) Buffer() []byte {
	return stream.buf
}

// SetBuffer allows to append to the internal buffer directly
func (stream *Stream) SetBuffer(buf []byte) {
	stream.buf = buf
}

// Write writes the contents of p into the buffer.
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (stream *Stream) Write(p []byte) (nn int, err error) {
	stream.buf = append(stream.buf, p...)
	stream.enforceMaxBytes()

	if stream.out != nil {
		nn, err = stream.out.Write(stream.buf)
		stream.buf = stream.buf[nn:]
		stream.marshalledBytesLimitRemaining -= uint64(nn)
		return
	}
	return len(p), nil
}

// WriteByte writes a single byte.
func (stream *Stream) writeByte(c byte) {
	stream.buf = append(stream.buf, c)
	stream.enforceMaxBytes()
}

func (stream *Stream) writeTwoBytes(c1 byte, c2 byte) {
	stream.buf = append(stream.buf, c1, c2)
	stream.enforceMaxBytes()
}

func (stream *Stream) writeThreeBytes(c1 byte, c2 byte, c3 byte) {
	stream.buf = append(stream.buf, c1, c2, c3)
	stream.enforceMaxBytes()
}

func (stream *Stream) writeFourBytes(c1 byte, c2 byte, c3 byte, c4 byte) {
	stream.buf = append(stream.buf, c1, c2, c3, c4)
	stream.enforceMaxBytes()
}

func (stream *Stream) writeFiveBytes(c1 byte, c2 byte, c3 byte, c4 byte, c5 byte) {
	stream.buf = append(stream.buf, c1, c2, c3, c4, c5)
	stream.enforceMaxBytes()
}

func (stream *Stream) enforceMaxBytes() {
	if !stream.enforceMarshalledBytesLimit {
		return
	}

	if uint64(len(stream.buf)) > stream.marshalledBytesLimitRemaining {
		// Why do we do this rather than return an error?
		// Most of the writing methods on Stream do not return an error, and introducing this would be a
		// breaking change for custom encoders.
		// Furthermore, nothing checks if the stream has failed until the object has been completely written
		// so if we don't panic here, we'd continue writing the rest of the object, negating the purpose of
		// this limit.
		panic(ExceededMaxMarshalledBytesError{stream.cfg.maxMarshalledBytes})
	}
}

type ExceededMaxMarshalledBytesError struct {
	MaxMarshalledBytes uint64
}

func (err ExceededMaxMarshalledBytesError) Error() string {
	return fmt.Sprintf("marshalling produced a result over the configured limit of %d bytes", err.MaxMarshalledBytes)
}

// Flush writes any buffered data to the underlying io.Writer.
func (stream *Stream) Flush() error {
	if stream.out == nil {
		return nil
	}
	if stream.Error != nil {
		return stream.Error
	}
	_, err := stream.out.Write(stream.buf)
	if err != nil {
		if stream.Error == nil {
			stream.Error = err
		}
		return err
	}
	stream.buf = stream.buf[:0]
	return nil
}

// WriteRaw write string out without quotes, just like []byte
func (stream *Stream) WriteRaw(s string) {
	stream.buf = append(stream.buf, s...)
	stream.enforceMaxBytes()
}

// WriteNil write null to stream
func (stream *Stream) WriteNil() {
	stream.writeFourBytes('n', 'u', 'l', 'l')
}

// WriteTrue write true to stream
func (stream *Stream) WriteTrue() {
	stream.writeFourBytes('t', 'r', 'u', 'e')
}

// WriteFalse write false to stream
func (stream *Stream) WriteFalse() {
	stream.writeFiveBytes('f', 'a', 'l', 's', 'e')
}

// WriteBool write true or false into stream
func (stream *Stream) WriteBool(val bool) {
	if val {
		stream.WriteTrue()
	} else {
		stream.WriteFalse()
	}
}

// WriteObjectStart write { with possible indention
func (stream *Stream) WriteObjectStart() {
	stream.indention += stream.cfg.indentionStep
	stream.writeByte('{')
	stream.writeIndention(0)
}

// WriteObjectField write "field": with possible indention
func (stream *Stream) WriteObjectField(field string) {
	stream.WriteString(field)
	if stream.indention > 0 {
		stream.writeTwoBytes(':', ' ')
	} else {
		stream.writeByte(':')
	}
}

// WriteObjectEnd write } with possible indention
func (stream *Stream) WriteObjectEnd() {
	stream.writeIndention(stream.cfg.indentionStep)
	stream.indention -= stream.cfg.indentionStep
	stream.writeByte('}')
}

// WriteEmptyObject write {}
func (stream *Stream) WriteEmptyObject() {
	stream.writeByte('{')
	stream.writeByte('}')
}

// WriteMore write , with possible indention
func (stream *Stream) WriteMore() {
	stream.writeByte(',')
	stream.writeIndention(0)
}

// WriteArrayStart write [ with possible indention
func (stream *Stream) WriteArrayStart() {
	stream.indention += stream.cfg.indentionStep
	stream.writeByte('[')
	stream.writeIndention(0)
}

// WriteEmptyArray write []
func (stream *Stream) WriteEmptyArray() {
	stream.writeTwoBytes('[', ']')
}

// WriteArrayEnd write ] with possible indention
func (stream *Stream) WriteArrayEnd() {
	stream.writeIndention(stream.cfg.indentionStep)
	stream.indention -= stream.cfg.indentionStep
	stream.writeByte(']')
}

func (stream *Stream) writeIndention(delta int) {
	if stream.indention == 0 {
		return
	}
	stream.writeByte('\n')
	toWrite := stream.indention - delta
	for i := 0; i < toWrite; i++ {
		stream.buf = append(stream.buf, ' ')
	}
	stream.enforceMaxBytes()
}
