package extra

import (
	"github.com/json-iterator/go"
	"unsafe"
	"time"
)

// keep epoch milliseconds
func RegisterTimeAsInt64Codec() {
	jsoniter.RegisterTypeEncoder("time.Time", &timeAsInt64Codec{})
}

type timeAsInt64Codec struct {
}

func (codec *timeAsInt64Codec) IsEmpty(ptr unsafe.Pointer) bool {
	ts := *((*time.Time)(ptr))
	return ts.UnixNano() == 0
}
func (codec *timeAsInt64Codec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	ts := *((*time.Time)(ptr))
	stream.WriteInt64(ts.UnixNano())
}
func (codec *timeAsInt64Codec) EncodeInterface(val interface{}, stream *jsoniter.Stream) {
	jsoniter.WriteToStream(val, stream, codec)
}
