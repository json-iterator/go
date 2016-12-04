package jsoniter

import (
	"reflect"
	"errors"
)

type Decoder interface {
	decode(iter *Iterator, obj interface{})
}

type stringDecoder struct {
}

func (decoder *stringDecoder) decode(iter *Iterator, obj interface{}) {
	ptr := obj.(*string)
	*ptr = iter.ReadString()
}

var DECODER_STRING *stringDecoder

func init() {
	DECODER_STRING = &stringDecoder{}
}

func (iter *Iterator) Read(obj interface{}) {
	type_ := reflect.TypeOf(obj)
	decoder, err := decoderOfType(type_)
	if err != nil {
		iter.Error = err
		return
	}
	decoder.decode(iter, obj)
}

func decoderOfType(type_ reflect.Type) (Decoder, error) {
	switch type_.Kind() {
	case reflect.Ptr:
		return decoderOfPtr(type_.Elem())
	default:
		return nil, errors.New("expect ptr")
	}
}

func decoderOfPtr(type_ reflect.Type) (Decoder, error) {
	switch type_.Kind() {
	case reflect.String:
		return DECODER_STRING, nil
	default:
		return nil, errors.New("expect string")
	}
}

