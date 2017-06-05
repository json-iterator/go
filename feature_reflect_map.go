package jsoniter

import (
	"unsafe"
	"reflect"
	"encoding/json"
	"encoding"
)

type mapDecoder struct {
	mapType      reflect.Type
	elemType     reflect.Type
	elemDecoder  Decoder
	mapInterface emptyInterface
}

func (decoder *mapDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	// dark magic to cast unsafe.Pointer back to interface{} using reflect.Type
	mapInterface := decoder.mapInterface
	mapInterface.word = ptr
	realInterface := (*interface{})(unsafe.Pointer(&mapInterface))
	realVal := reflect.ValueOf(*realInterface).Elem()
	if realVal.IsNil() {
		realVal.Set(reflect.MakeMap(realVal.Type()))
	}
	for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
		elem := reflect.New(decoder.elemType)
		decoder.elemDecoder.decode(unsafe.Pointer(elem.Pointer()), iter)
		// to put into map, we have to use reflection
		realVal.SetMapIndex(reflect.ValueOf(string([]byte(field))), elem.Elem())
	}
}

type mapEncoder struct {
	mapType      reflect.Type
	elemType     reflect.Type
	elemEncoder  Encoder
	mapInterface emptyInterface
}

func (encoder *mapEncoder) encode(ptr unsafe.Pointer, stream *Stream) {
	mapInterface := encoder.mapInterface
	mapInterface.word = ptr
	realInterface := (*interface{})(unsafe.Pointer(&mapInterface))
	realVal := reflect.ValueOf(*realInterface)

	stream.WriteObjectStart()
	for i, key := range realVal.MapKeys() {
		if i != 0 {
			stream.WriteMore()
		}
		encodeMapKey(key, stream)
		stream.writeByte(':')
		val := realVal.MapIndex(key).Interface()
		encoder.elemEncoder.encodeInterface(val, stream)
	}
	stream.WriteObjectEnd()
}

func encodeMapKey(key reflect.Value, stream *Stream) {
	if key.Kind() == reflect.String {
		stream.WriteString(key.String())
		return
	}
	if tm, ok := key.Interface().(encoding.TextMarshaler); ok {
		buf, err := tm.MarshalText()
		if err != nil {
			stream.Error = err
			return
		}
		stream.writeByte('"')
		stream.Write(buf)
		stream.writeByte('"')
		return
	}
	switch key.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		stream.writeByte('"')
		stream.WriteInt64(key.Int())
		stream.writeByte('"')
		return
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		stream.writeByte('"')
		stream.WriteUint64(key.Uint())
		stream.writeByte('"')
		return
	}
	stream.Error = &json.UnsupportedTypeError{key.Type()}
}

func (encoder *mapEncoder) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (encoder *mapEncoder) isEmpty(ptr unsafe.Pointer) bool {
	mapInterface := encoder.mapInterface
	mapInterface.word = ptr
	realInterface := (*interface{})(unsafe.Pointer(&mapInterface))
	realVal := reflect.ValueOf(*realInterface)
	return realVal.Len() == 0
}

type mapInterfaceEncoder struct {
	mapType      reflect.Type
	elemType     reflect.Type
	elemEncoder  Encoder
	mapInterface emptyInterface
}

func (encoder *mapInterfaceEncoder) encode(ptr unsafe.Pointer, stream *Stream) {
	mapInterface := encoder.mapInterface
	mapInterface.word = ptr
	realInterface := (*interface{})(unsafe.Pointer(&mapInterface))
	realVal := reflect.ValueOf(*realInterface)

	stream.WriteObjectStart()
	for i, key := range realVal.MapKeys() {
		if i != 0 {
			stream.WriteMore()
		}
		switch key.Interface().(type) {
		case string:
			stream.WriteObjectField(key.String())
		default:
			stream.Error = &json.UnsupportedTypeError{key.Type()}
			return
		}
		val := realVal.MapIndex(key).Interface()
		encoder.elemEncoder.encode(unsafe.Pointer(&val), stream)
	}
	stream.WriteObjectEnd()
}

func (encoder *mapInterfaceEncoder) encodeInterface(val interface{}, stream *Stream) {
	writeToStream(val, stream, encoder)
}

func (encoder *mapInterfaceEncoder) isEmpty(ptr unsafe.Pointer) bool {
	mapInterface := encoder.mapInterface
	mapInterface.word = ptr
	realInterface := (*interface{})(unsafe.Pointer(&mapInterface))
	realVal := reflect.ValueOf(*realInterface)

	return realVal.Len() == 0
}