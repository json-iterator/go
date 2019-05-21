package jsoniter

import (
	"github.com/junchih/stringx"
	"io"
	"strings"
)

// IteratorPool a thread safe pool of iterators with same configuration
type IteratorPool interface {
	BorrowIterator(data []byte) *Iterator
	ReturnIterator(iter *Iterator)
}

// StreamPool a thread safe pool of streams with same configuration
type StreamPool interface {
	BorrowStream(writer io.Writer) *Stream
	ReturnStream(stream *Stream)
}

func (cfg *frozenConfig) BorrowStream(writer io.Writer) *Stream {
	stream := cfg.streamPool.Get().(*Stream)
	stream.Reset(writer)
	return stream
}

func (cfg *frozenConfig) ReturnStream(stream *Stream) {
	stream.out = nil
	stream.Error = nil
	stream.Attachment = nil
	cfg.streamPool.Put(stream)
}

func (cfg *frozenConfig) BorrowIterator(data []byte) *Iterator {
	iter := cfg.iteratorPool.Get().(*Iterator)
	iter.ResetBytes(data)
	return iter
}

func (cfg *frozenConfig) ReturnIterator(iter *Iterator) {
	iter.Error = nil
	iter.Attachment = nil
	cfg.iteratorPool.Put(iter)
}

func (cfg *frozenConfig) borrowStringBuilder() *strings.Builder {
	return cfg.stringBuilderPool.Get().(*strings.Builder)
}

func (cfg *frozenConfig) returnStringBuilder(builder *strings.Builder) {
	builder.Reset()
	cfg.stringBuilderPool.Put(builder)
}

func (cfg *frozenConfig) borrowStringFactory() *stringx.Factory {
	return cfg.stringFactoryPool.Get().(*stringx.Factory)
}

func (cfg *frozenConfig) returnStringFactory(factory *stringx.Factory) {
	cfg.stringFactoryPool.Put(factory)
}
