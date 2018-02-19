package jsoniter

import (
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

func decoderOfArray(cfg *frozenConfig, prefix string, typ reflect.Type) ValDecoder {
	decoder := decoderOfType(cfg, prefix+"[array]->", typ.Elem())
	return &arrayDecoder{typ, typ.Elem(), decoder}
}

func encoderOfArray(cfg *frozenConfig, prefix string, typ reflect.Type) ValEncoder {
	if typ.Len() == 0 {
		return emptyArrayEncoder{}
	}
	encoder := encoderOfType(cfg, prefix+"[array]->", typ.Elem())
	return &arrayEncoder{typ, typ.Elem(), encoder}
}

type emptyArrayEncoder struct{}

func (encoder emptyArrayEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteEmptyArray()
}

func (encoder emptyArrayEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return true
}

type arrayEncoder struct {
	arrayType   reflect.Type
	elemType    reflect.Type
	elemEncoder ValEncoder
}

func (encoder *arrayEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteArrayStart()
	elemPtr := unsafe.Pointer(ptr)
	encoder.elemEncoder.Encode(elemPtr, stream)
	for i := 1; i < encoder.arrayType.Len(); i++ {
		stream.WriteMore()
		elemPtr = unsafe.Pointer(uintptr(elemPtr) + encoder.elemType.Size())
		encoder.elemEncoder.Encode(unsafe.Pointer(elemPtr), stream)
	}
	stream.WriteArrayEnd()
	if stream.Error != nil && stream.Error != io.EOF {
		stream.Error = fmt.Errorf("%v: %s", encoder.arrayType, stream.Error.Error())
	}
}

func (encoder *arrayEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

type arrayDecoder struct {
	arrayType   reflect.Type
	elemType    reflect.Type
	elemDecoder ValDecoder
}

func (decoder *arrayDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	decoder.doDecode(ptr, iter)
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.arrayType, iter.Error.Error())
	}
}

func (decoder *arrayDecoder) doDecode(ptr unsafe.Pointer, iter *Iterator) {
	offset := uintptr(0)
	iter.ReadArrayCB(func(iter *Iterator) bool {
		if offset < decoder.arrayType.Size() {
			decoder.elemDecoder.Decode(unsafe.Pointer(uintptr(ptr)+offset), iter)
			offset += decoder.elemType.Size()
		} else {
			iter.Skip()
		}
		return true
	})
}
