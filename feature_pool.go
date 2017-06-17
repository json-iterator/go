package jsoniter

func (cfg *frozenConfig) borrowStream() *Stream {
	select {
	case stream := <-cfg.streamPool:
		stream.Reset(nil)
		return stream
	default:
		return NewStream(cfg, nil, 512)
	}
}

func (cfg *frozenConfig) returnStream(stream *Stream) {
	select {
	case cfg.streamPool <- stream:
		return
	default:
		return
	}
}

func (cfg *frozenConfig) borrowIterator(data []byte) *Iterator {
	select {
	case iter := <- cfg.iteratorPool:
		iter.ResetBytes(data)
		return iter
	default:
		return ParseBytes(cfg, data)
	}
}

func (cfg *frozenConfig) returnIterator(iter *Iterator) {
	select {
	case cfg.iteratorPool <- iter:
		return
	default:
		return
	}
}
