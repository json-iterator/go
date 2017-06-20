package jsoniter

import (
	"reflect"
	"fmt"
	"unsafe"
	"strings"
	"unicode"
)

var typeDecoders = map[string]ValDecoder{}
var fieldDecoders = map[string]ValDecoder{}
var typeEncoders = map[string]ValEncoder{}
var fieldEncoders = map[string]ValEncoder{}
var extensions = []Extension{}

type StructDescriptor struct {
	Type   reflect.Type
	Fields []*Binding
}

func (structDescriptor *StructDescriptor) GetField(fieldName string) *Binding {
	for _, binding := range structDescriptor.Fields {
		if binding.Field.Name == fieldName {
			return binding
		}
	}
	return nil
}

type Binding struct {
	Field           *reflect.StructField
	FromNames       []string
	ToNames         []string
	Encoder         ValEncoder
	Decoder         ValDecoder
}

type Extension interface {
	UpdateStructDescriptor(structDescriptor *StructDescriptor)
	CreateDecoder(typ reflect.Type) ValDecoder
	CreateEncoder(typ reflect.Type) ValEncoder
}

type DummyExtension struct {
}

func (extension *DummyExtension) UpdateStructDescriptor(structDescriptor *StructDescriptor) {
}

func (extension *DummyExtension) CreateDecoder(typ reflect.Type) ValDecoder {
	return nil
}

func (extension *DummyExtension) CreateEncoder(typ reflect.Type) ValEncoder {
	return nil
}

type funcDecoder struct {
	fun DecoderFunc
}

func RegisterTypeDecoderFunc(typ string, fun DecoderFunc) {
	typeDecoders[typ] = &funcDecoder{fun}
}

func RegisterTypeDecoder(typ string, decoder ValDecoder) {
	typeDecoders[typ] = decoder
}

func RegisterFieldDecoderFunc(typ string, field string, fun DecoderFunc) {
	RegisterFieldDecoder(typ, field, &funcDecoder{fun})
}

func RegisterFieldDecoder(typ string, field string, decoder ValDecoder) {
	fieldDecoders[fmt.Sprintf("%s/%s", typ, field)] = decoder
}

func RegisterTypeEncoderFunc(typ string, fun EncoderFunc, isEmptyFunc func(unsafe.Pointer) bool) {
	typeEncoders[typ] = &funcEncoder{fun, isEmptyFunc}
}

func RegisterTypeEncoder(typ string, encoder ValEncoder) {
	typeEncoders[typ] = encoder
}

func RegisterFieldEncoderFunc(typ string, field string, fun EncoderFunc, isEmptyFunc func(unsafe.Pointer) bool) {
	RegisterFieldEncoder(typ, field, &funcEncoder{fun, isEmptyFunc})
}

func RegisterFieldEncoder(typ string, field string, encoder ValEncoder) {
	fieldEncoders[fmt.Sprintf("%s/%s", typ, field)] = encoder
}

func RegisterExtension(extension Extension) {
	extensions = append(extensions, extension)
}

func getTypeDecoderFromExtension(typ reflect.Type) ValDecoder {
	for _, extension := range extensions {
		decoder := extension.CreateDecoder(typ)
		if decoder != nil {
			return decoder
		}
	}
	typeName := typ.String()
	decoder := typeDecoders[typeName]
	if decoder != nil {
		return decoder
	}
	if typ.Kind() == reflect.Ptr {
		decoder := typeDecoders[typ.Elem().String()]
		if decoder != nil {
			return &optionalDecoder{typ.Elem(), decoder}
		}
	}
	return nil
}

func getTypeEncoderFromExtension(typ reflect.Type) ValEncoder {
	for _, extension := range extensions {
		encoder := extension.CreateEncoder(typ)
		if encoder != nil {
			return encoder
		}
	}
	typeName := typ.String()
	encoder := typeEncoders[typeName]
	if encoder != nil {
		return encoder
	}
	if typ.Kind() == reflect.Ptr {
		encoder := typeEncoders[typ.Elem().String()]
		if encoder != nil {
			return &optionalEncoder{encoder}
		}
	}
	return nil
}

func describeStruct(cfg *frozenConfig, typ reflect.Type) (*StructDescriptor, error) {
	bindings := []*Binding{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Anonymous {
			if field.Type.Kind() == reflect.Struct {
				structDescriptor, err := describeStruct(cfg, field.Type)
				if err != nil {
					return nil, err
				}
				for _, binding := range structDescriptor.Fields {
					bindings = append(bindings, binding)
				}
			} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
				structDescriptor, err := describeStruct(cfg, field.Type.Elem())
				if err != nil {
					return nil, err
				}
				for _, binding := range structDescriptor.Fields {
					binding.Encoder = &optionalEncoder{binding.Encoder}
					binding.Encoder = &structFieldEncoder{&field, binding.Encoder, false}
					binding.Decoder = &optionalDecoder{field.Type, binding.Decoder}
					binding.Decoder = &structFieldDecoder{&field, binding.Decoder}
					bindings = append(bindings, binding)
				}
			}
		} else {
			tagParts := strings.Split(field.Tag.Get("json"), ",")
			fieldNames := calcFieldNames(field.Name, tagParts[0])
			fieldCacheKey := fmt.Sprintf("%s/%s", typ.String(), field.Name)
			decoder := fieldDecoders[fieldCacheKey]
			if decoder == nil && len(fieldNames) > 0 {
				var err error
				decoder, err = decoderOfType(cfg, field.Type)
				if err != nil {
					return nil, err
				}
			}
			encoder := fieldEncoders[fieldCacheKey]
			if encoder == nil && len(fieldNames) > 0 {
				var err error
				encoder, err = encoderOfType(cfg, field.Type)
				if err != nil {
					return nil, err
				}
				// map is stored as pointer in the struct
				if field.Type.Kind() == reflect.Map {
					encoder = &optionalEncoder{encoder}
				}
			}
			binding := &Binding{
				Field:     &field,
				FromNames: fieldNames,
				ToNames:   fieldNames,
				Decoder:   decoder,
				Encoder:   encoder,
			}
			shouldOmitEmpty := false
			for _, tagPart := range tagParts[1:] {
				if tagPart == "omitempty" {
					shouldOmitEmpty = true
				} else if tagPart == "string" {
					binding.Decoder = &stringModeDecoder{binding.Decoder}
					binding.Encoder = &stringModeEncoder{binding.Encoder}
				}
			}
			binding.Decoder = &structFieldDecoder{&field, binding.Decoder}
			binding.Encoder = &structFieldEncoder{&field, binding.Encoder, shouldOmitEmpty}
			bindings = append(bindings, binding)
		}
	}
	structDescriptor := &StructDescriptor{
		Type:   typ,
		Fields: bindings,
	}
	for _, extension := range extensions {
		extension.UpdateStructDescriptor(structDescriptor)
	}
	return structDescriptor, nil
}

func listStructFields(typ reflect.Type) []*reflect.StructField {
	fields := []*reflect.StructField{}
	return fields
}

func calcFieldNames(originalFieldName string, tagProvidedFieldName string) []string {
	// tag => exported? => original
	isNotExported := unicode.IsLower(rune(originalFieldName[0]))
	var fieldNames []string
	/// tagParts[0] always present, even if no tags
	switch tagProvidedFieldName {
	case "":
		if isNotExported {
			fieldNames = []string{}
		} else {
			fieldNames = []string{originalFieldName}
		}
	case "-":
		fieldNames = []string{}
	default:
		fieldNames = []string{tagProvidedFieldName}
	}
	return fieldNames
}