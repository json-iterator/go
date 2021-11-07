package jsoniter

import (
	"encoding/json"
	"github.com/modern-go/reflect2"
	"unsafe"
)

var jsonRawMessageType = reflect2.TypeOfPtr((*json.RawMessage)(nil)).Elem()

func createEncoderOfJsonRawMessage(ctx *ctx, typ reflect2.Type) ValEncoder {
	if typ == jsonRawMessageType {
		return &jsonRawMessageCodec{}
	}
	return nil
}

func createDecoderOfJsonRawMessage(ctx *ctx, typ reflect2.Type) ValDecoder {
	if typ == jsonRawMessageType {
		return &jsonRawMessageCodec{}
	}
	return nil
}

type jsonRawMessageCodec struct {
}

func (codec *jsonRawMessageCodec) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if iter.ReadNil() {
		*((*json.RawMessage)(ptr)) = nil
	} else {
		*((*json.RawMessage)(ptr)) = iter.SkipAndReturnBytes()
	}
}

func (codec *jsonRawMessageCodec) Encode(ptr unsafe.Pointer, stream *Stream) {
	if *((*json.RawMessage)(ptr)) == nil {
		stream.WriteNil()
	} else {
		stream.WriteRaw(string(*((*json.RawMessage)(ptr))))
	}
}

func (codec *jsonRawMessageCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return len(*((*json.RawMessage)(ptr))) == 0
}
