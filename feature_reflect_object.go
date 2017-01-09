package jsoniter

import (
	"io"
	"fmt"
	"reflect"
	"unsafe"
	"strings"
)


func encoderOfStruct(typ reflect.Type) (Encoder, error) {
	structEncoder_ := &structEncoder{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		var fieldNames []string
		for _, extension := range extensions {
			alternativeFieldNames, _ := extension(typ, &field)
			if alternativeFieldNames != nil {
				fieldNames = alternativeFieldNames
			}
		}
		tagParts := strings.Split(field.Tag.Get("json"), ",")
		// if fieldNames set by extension, use theirs, otherwise try tags
		if fieldNames == nil {
			/// tagParts[0] always present, even if no tags
			switch tagParts[0] {
			case "":
				fieldNames = []string{field.Name}
			case "-":
				fieldNames = []string{}
			default:
				fieldNames = []string{tagParts[0]}
			}
		}
		encoder, err := encoderOfType(field.Type)
		if err != nil {
			return prefix(fmt.Sprintf("{%s}", field.Name)).addToEncoder(encoder, err)
		}
		for _, fieldName := range fieldNames {
			if structEncoder_.firstField == nil {
				structEncoder_.firstField = &structFieldEncoder{&field, fieldName, encoder}
			} else {
				structEncoder_.fields = append(structEncoder_.fields, &structFieldEncoder{&field, fieldName, encoder})
			}
		}
	}
	if structEncoder_.firstField == nil {
		return &emptyStructEncoder{}, nil
	}
	return structEncoder_, nil
}

func decoderOfStruct(typ reflect.Type) (Decoder, error) {
	fields := map[string]*structFieldDecoder{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldDecoderKey := fmt.Sprintf("%s/%s", typ.String(), field.Name)
		var fieldNames []string
		for _, extension := range extensions {
			alternativeFieldNames, fun := extension(typ, &field)
			if alternativeFieldNames != nil {
				fieldNames = alternativeFieldNames
			}
			if fun != nil {
				fieldDecoders[fieldDecoderKey] = &funcDecoder{fun}
			}
		}
		decoder := fieldDecoders[fieldDecoderKey]
		tagParts := strings.Split(field.Tag.Get("json"), ",")
		// if fieldNames set by extension, use theirs, otherwise try tags
		if fieldNames == nil {
			/// tagParts[0] always present, even if no tags
			switch tagParts[0] {
			case "":
				fieldNames = []string{field.Name}
			case "-":
				fieldNames = []string{}
			default:
				fieldNames = []string{tagParts[0]}
			}
		}
		if decoder == nil {
			var err error
			decoder, err = decoderOfType(field.Type)
			if err != nil {
				return prefix(fmt.Sprintf("{%s}", field.Name)).addToDecoder(decoder, err)
			}
		}
		if len(tagParts) > 1 && tagParts[1] == "string" {
			decoder = &stringNumberDecoder{decoder}
		}
		for _, fieldName := range fieldNames {
			fields[fieldName] = &structFieldDecoder{&field, decoder}
		}
	}
	switch len(fields) {
	case 0:
		return &skipDecoder{typ}, nil
	case 1:
		for fieldName, fieldDecoder := range fields {
			return &oneFieldStructDecoder{typ, fieldName, fieldDecoder}, nil
		}
	case 2:
		var fieldName1 string
		var fieldName2 string
		var fieldDecoder1 *structFieldDecoder
		var fieldDecoder2 *structFieldDecoder
		for fieldName, fieldDecoder := range fields {
			if fieldName1 == "" {
				fieldName1 = fieldName
				fieldDecoder1 = fieldDecoder
			} else {
				fieldName2 = fieldName
				fieldDecoder2 = fieldDecoder
			}
		}
		return &twoFieldsStructDecoder{typ, fieldName1, fieldDecoder1, fieldName2, fieldDecoder2}, nil
	case 3:
		var fieldName1 string
		var fieldName2 string
		var fieldName3 string
		var fieldDecoder1 *structFieldDecoder
		var fieldDecoder2 *structFieldDecoder
		var fieldDecoder3 *structFieldDecoder
		for fieldName, fieldDecoder := range fields {
			if fieldName1 == "" {
				fieldName1 = fieldName
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == "" {
				fieldName2 = fieldName
				fieldDecoder2 = fieldDecoder
			} else {
				fieldName3 = fieldName
				fieldDecoder3 = fieldDecoder
			}
		}
		return &threeFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1, fieldName2, fieldDecoder2, fieldName3, fieldDecoder3}, nil
	case 4:
		var fieldName1 string
		var fieldName2 string
		var fieldName3 string
		var fieldName4 string
		var fieldDecoder1 *structFieldDecoder
		var fieldDecoder2 *structFieldDecoder
		var fieldDecoder3 *structFieldDecoder
		var fieldDecoder4 *structFieldDecoder
		for fieldName, fieldDecoder := range fields {
			if fieldName1 == "" {
				fieldName1 = fieldName
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == "" {
				fieldName2 = fieldName
				fieldDecoder2 = fieldDecoder
			} else if fieldName3 == "" {
				fieldName3 = fieldName
				fieldDecoder3 = fieldDecoder
			} else {
				fieldName4 = fieldName
				fieldDecoder4 = fieldDecoder
			}
		}
		return &fourFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1, fieldName2, fieldDecoder2, fieldName3, fieldDecoder3,
			fieldName4, fieldDecoder4}, nil
	}
	return &generalStructDecoder{typ, fields}, nil
}

type generalStructDecoder struct {
	typ    reflect.Type
	fields map[string]*structFieldDecoder
}

func (decoder *generalStructDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	field := iter.readObjectField()
	fieldDecoder := decoder.fields[field]
	if fieldDecoder == nil {
		iter.Skip()
	} else {
		fieldDecoder.decode(ptr, iter)
	}
	for iter.nextToken() == ',' {
		field = iter.readObjectField()
		fieldDecoder = decoder.fields[field]
		if fieldDecoder == nil {
			iter.Skip()
		} else {
			fieldDecoder.decode(ptr, iter)
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type skipDecoder struct {
	typ reflect.Type
}

func (decoder *skipDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	iter.Skip()
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type oneFieldStructDecoder struct {
	typ          reflect.Type
	fieldName    string
	fieldDecoder *structFieldDecoder
}

func (decoder *oneFieldStructDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	field := iter.readObjectField()
	if field == decoder.fieldName {
		decoder.fieldDecoder.decode(ptr, iter)
	} else {
		iter.Skip()
	}
	for iter.nextToken() == ',' {
		field = iter.readObjectField()
		if field == decoder.fieldName {
			decoder.fieldDecoder.decode(ptr, iter)
		} else {
			iter.Skip()
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type twoFieldsStructDecoder struct {
	typ           reflect.Type
	fieldName1    string
	fieldDecoder1 *structFieldDecoder
	fieldName2    string
	fieldDecoder2 *structFieldDecoder
}

func (decoder *twoFieldsStructDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	field := iter.readObjectField()
	switch field {
	case decoder.fieldName1:
		decoder.fieldDecoder1.decode(ptr, iter)
	case decoder.fieldName2:
		decoder.fieldDecoder2.decode(ptr, iter)
	default:
		iter.Skip()
	}
	for iter.nextToken() == ',' {
		field = iter.readObjectField()
		switch field {
		case decoder.fieldName1:
			decoder.fieldDecoder1.decode(ptr, iter)
		case decoder.fieldName2:
			decoder.fieldDecoder2.decode(ptr, iter)
		default:
			iter.Skip()
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type threeFieldsStructDecoder struct {
	typ           reflect.Type
	fieldName1    string
	fieldDecoder1 *structFieldDecoder
	fieldName2    string
	fieldDecoder2 *structFieldDecoder
	fieldName3    string
	fieldDecoder3 *structFieldDecoder
}

func (decoder *threeFieldsStructDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	field := iter.readObjectField()
	switch field {
	case decoder.fieldName1:
		decoder.fieldDecoder1.decode(ptr, iter)
	case decoder.fieldName2:
		decoder.fieldDecoder2.decode(ptr, iter)
	case decoder.fieldName3:
		decoder.fieldDecoder3.decode(ptr, iter)
	default:
		iter.Skip()
	}
	for iter.nextToken() == ',' {
		field = iter.readObjectField()
		switch field {
		case decoder.fieldName1:
			decoder.fieldDecoder1.decode(ptr, iter)
		case decoder.fieldName2:
			decoder.fieldDecoder2.decode(ptr, iter)
		case decoder.fieldName3:
			decoder.fieldDecoder3.decode(ptr, iter)
		default:
			iter.Skip()
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type fourFieldsStructDecoder struct {
	typ           reflect.Type
	fieldName1    string
	fieldDecoder1 *structFieldDecoder
	fieldName2    string
	fieldDecoder2 *structFieldDecoder
	fieldName3    string
	fieldDecoder3 *structFieldDecoder
	fieldName4    string
	fieldDecoder4 *structFieldDecoder
}

func (decoder *fourFieldsStructDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	field := iter.readObjectField()
	switch field {
	case decoder.fieldName1:
		decoder.fieldDecoder1.decode(ptr, iter)
	case decoder.fieldName2:
		decoder.fieldDecoder2.decode(ptr, iter)
	case decoder.fieldName3:
		decoder.fieldDecoder3.decode(ptr, iter)
	case decoder.fieldName4:
		decoder.fieldDecoder4.decode(ptr, iter)
	default:
		iter.Skip()
	}
	for iter.nextToken() == ',' {
		field = iter.readObjectField()
		switch field {
		case decoder.fieldName1:
			decoder.fieldDecoder1.decode(ptr, iter)
		case decoder.fieldName2:
			decoder.fieldDecoder2.decode(ptr, iter)
		case decoder.fieldName3:
			decoder.fieldDecoder3.decode(ptr, iter)
		case decoder.fieldName4:
			decoder.fieldDecoder4.decode(ptr, iter)
		default:
			iter.Skip()
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type structFieldDecoder struct {
	field        *reflect.StructField
	fieldDecoder Decoder
}

func (decoder *structFieldDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	fieldPtr := uintptr(ptr) + decoder.field.Offset
	decoder.fieldDecoder.decode(unsafe.Pointer(fieldPtr), iter)
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%s: %s", decoder.field.Name, iter.Error.Error())
	}
}

type structFieldEncoder struct {
	field        *reflect.StructField
	fieldName    string
	fieldEncoder Encoder
}

func (encoder *structFieldEncoder) encode(ptr unsafe.Pointer, stream *Stream) {
	fieldPtr := uintptr(ptr) + encoder.field.Offset
	stream.WriteObjectField(encoder.fieldName)
	encoder.fieldEncoder.encode(unsafe.Pointer(fieldPtr), stream)
	if stream.Error != nil && stream.Error != io.EOF {
		stream.Error = fmt.Errorf("%s: %s", encoder.field.Name, stream.Error.Error())
	}
}


type structEncoder struct {
	firstField *structFieldEncoder
	fields []*structFieldEncoder
}

func (encoder *structEncoder) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteObjectStart()
	encoder.firstField.encode(ptr, stream)
	for _, field := range encoder.fields {
		stream.WriteMore()
		field.encode(ptr, stream)
	}
	stream.WriteObjectEnd()
}

type emptyStructEncoder struct {
}

func (encoder *emptyStructEncoder) encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteEmptyObject()
}