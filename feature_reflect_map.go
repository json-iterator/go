package jsoniter

import (
	"encoding"
	"reflect"
	"sort"
	"strconv"
	"unsafe"
	"github.com/v2pro/plz/reflect2"
	"fmt"
)

func decoderOfMap(cfg *frozenConfig, prefix string, typ reflect.Type) ValDecoder {
	decoder := decoderOfType(cfg, prefix+"[map]->", typ.Elem())
	mapInterface := reflect.New(typ).Interface()
	return &mapDecoder{typ, typ.Key(), typ.Elem(), decoder, extractInterface(mapInterface)}
}

func encoderOfMap(cfg *frozenConfig, prefix string, typ reflect.Type) ValEncoder {
	if cfg.sortMapKeys {
		return &sortKeysMapEncoder{
			mapType:     reflect2.Type2(typ).(*reflect2.UnsafeMapType),
			keyEncoder:  encoderOfMapKey(cfg, prefix+" [mapKey]", typ.Key()),
			elemEncoder: encoderOfType(cfg, prefix+" [mapElem]", typ.Elem()),
		}
	}
	return &mapEncoder{
		mapType:     reflect2.Type2(typ).(*reflect2.UnsafeMapType),
		keyEncoder:  encoderOfMapKey(cfg, prefix+" [mapKey]", typ.Key()),
		elemEncoder: encoderOfType(cfg, prefix+" [mapElem]", typ.Elem()),
	}
}

func encoderOfMapKey(cfg *frozenConfig, prefix string, typ reflect.Type) ValEncoder {
	switch typ.Kind() {
	case reflect.String:
		return encoderOfType(cfg, prefix, reflect2.DefaultTypeOfKind(reflect.String).Type1())
	case reflect.Bool,
		reflect.Uint8, reflect.Int8,
		reflect.Uint16, reflect.Int16,
		reflect.Uint32, reflect.Int32,
		reflect.Uint64, reflect.Int64,
		reflect.Uint, reflect.Int,
		reflect.Float32, reflect.Float64,
		reflect.Uintptr:
		typ = reflect2.DefaultTypeOfKind(typ.Kind()).Type1()
		return &numericMapKeyEncoder{encoderOfType(cfg, prefix, typ)}
	default:
		if typ == textMarshalerType {
			return &directTextMarshalerEncoder{
				stringEncoder: cfg.EncoderOf(reflect.TypeOf("")),
			}
		}
		if typ.Implements(textMarshalerType) {
			return &textMarshalerEncoder{
				valType: reflect2.Type2(typ),
				stringEncoder: cfg.EncoderOf(reflect.TypeOf("")),
			}
		}
		return &lazyErrorEncoder{err: fmt.Errorf("unsupported map key type: %v", typ)}
	}
}

type mapDecoder struct {
	mapType      reflect.Type
	keyType      reflect.Type
	elemType     reflect.Type
	elemDecoder  ValDecoder
	mapInterface emptyInterface
}

func (decoder *mapDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	// dark magic to cast unsafe.Pointer back to interface{} using reflect.Type
	mapInterface := decoder.mapInterface
	mapInterface.word = ptr
	realInterface := (*interface{})(unsafe.Pointer(&mapInterface))
	realVal := reflect.ValueOf(*realInterface).Elem()
	if iter.ReadNil() {
		realVal.Set(reflect.Zero(decoder.mapType))
		return
	}
	if realVal.IsNil() {
		realVal.Set(reflect.MakeMap(realVal.Type()))
	}
	iter.ReadMapCB(func(iter *Iterator, keyStr string) bool {
		elem := reflect.New(decoder.elemType)
		decoder.elemDecoder.Decode(extractInterface(elem.Interface()).word, iter)
		// to put into map, we have to use reflection
		keyType := decoder.keyType
		// TODO: remove this from loop
		switch {
		case keyType.Kind() == reflect.String:
			realVal.SetMapIndex(reflect.ValueOf(keyStr).Convert(keyType), elem.Elem())
			return true
		case keyType.Implements(textUnmarshalerType):
			textUnmarshaler := reflect.New(keyType.Elem()).Interface().(encoding.TextUnmarshaler)
			err := textUnmarshaler.UnmarshalText([]byte(keyStr))
			if err != nil {
				iter.ReportError("read map key as TextUnmarshaler", err.Error())
				return false
			}
			realVal.SetMapIndex(reflect.ValueOf(textUnmarshaler), elem.Elem())
			return true
		case reflect.PtrTo(keyType).Implements(textUnmarshalerType):
			textUnmarshaler := reflect.New(keyType).Interface().(encoding.TextUnmarshaler)
			err := textUnmarshaler.UnmarshalText([]byte(keyStr))
			if err != nil {
				iter.ReportError("read map key as TextUnmarshaler", err.Error())
				return false
			}
			realVal.SetMapIndex(reflect.ValueOf(textUnmarshaler).Elem(), elem.Elem())
			return true
		default:
			switch keyType.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				n, err := strconv.ParseInt(keyStr, 10, 64)
				if err != nil || reflect.Zero(keyType).OverflowInt(n) {
					iter.ReportError("read map key as int64", "read int64 failed")
					return false
				}
				realVal.SetMapIndex(reflect.ValueOf(n).Convert(keyType), elem.Elem())
				return true
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				n, err := strconv.ParseUint(keyStr, 10, 64)
				if err != nil || reflect.Zero(keyType).OverflowUint(n) {
					iter.ReportError("read map key as uint64", "read uint64 failed")
					return false
				}
				realVal.SetMapIndex(reflect.ValueOf(n).Convert(keyType), elem.Elem())
				return true
			}
		}
		iter.ReportError("read map key", "unexpected map key type "+keyType.String())
		return true
	})
}

type numericMapKeyEncoder struct {
	encoder ValEncoder
}

func (encoder *numericMapKeyEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.writeByte('"')
	encoder.encoder.Encode(ptr, stream)
	stream.writeByte('"')
}

func (encoder *numericMapKeyEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

type mapEncoder struct {
	mapType     *reflect2.UnsafeMapType
	keyEncoder  ValEncoder
	elemEncoder ValEncoder
}

func (encoder *mapEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteObjectStart()
	iter := encoder.mapType.UnsafeIterate(ptr)
	for i := 0; iter.HasNext(); i++ {
		if i != 0 {
			stream.WriteMore()
		}
		key, elem := iter.UnsafeNext()
		encoder.keyEncoder.Encode(key, stream)
		if stream.indention > 0 {
			stream.writeTwoBytes(byte(':'), byte(' '))
		} else {
			stream.writeByte(':')
		}
		encoder.elemEncoder.Encode(elem, stream)
	}
	stream.WriteObjectEnd()
}

func (encoder *mapEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	iter := encoder.mapType.UnsafeIterate(ptr)
	return !iter.HasNext()
}

type sortKeysMapEncoder struct {
	mapType     *reflect2.UnsafeMapType
	keyEncoder  ValEncoder
	elemEncoder ValEncoder
}

func (encoder *sortKeysMapEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	if *(*unsafe.Pointer)(ptr) == nil {
		stream.WriteNil()
		return
	}
	stream.WriteObjectStart()
	mapIter := encoder.mapType.UnsafeIterate(ptr)
	subStream := stream.cfg.BorrowStream(nil)
	subIter := stream.cfg.BorrowIterator(nil)
	keyValues := encodedKeyValues{}
	for mapIter.HasNext() {
		subStream.buf = make([]byte, 0, 64)
		key, elem := mapIter.UnsafeNext()
		encoder.keyEncoder.Encode(key, subStream)
		encodedKey := subStream.Buffer()
		subIter.ResetBytes(encodedKey)
		decodedKey := subIter.ReadString()
		if stream.indention > 0 {
			subStream.writeTwoBytes(byte(':'), byte(' '))
		} else {
			subStream.writeByte(':')
		}
		encoder.elemEncoder.Encode(elem, subStream)
		keyValues = append(keyValues, encodedKV{
			key:      decodedKey,
			keyValue: subStream.Buffer(),
		})
	}
	sort.Sort(keyValues)
	for i, keyValue := range keyValues {
		if i != 0 {
			stream.WriteMore()
		}
		stream.Write(keyValue.keyValue)
	}
	stream.WriteObjectEnd()
	stream.cfg.ReturnStream(subStream)
	stream.cfg.ReturnIterator(subIter)
}

func (encoder *sortKeysMapEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	iter := encoder.mapType.UnsafeIterate(ptr)
	return !iter.HasNext()
}

type encodedKeyValues []encodedKV

type encodedKV struct {
	key      string
	keyValue []byte
}

func (sv encodedKeyValues) Len() int           { return len(sv) }
func (sv encodedKeyValues) Swap(i, j int)      { sv[i], sv[j] = sv[j], sv[i] }
func (sv encodedKeyValues) Less(i, j int) bool { return sv[i].key < sv[j].key }
