package jsoniter

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"unsafe"

	"github.com/modern-go/reflect2"
)

type binding struct {
	binding *Binding
	name    string
	hasTag  bool
}

func encoderOfStruct(ctx *ctx, typ reflect2.Type) ValEncoder {

	orderedBindings := []*binding{}
	structDescriptor := describeStruct(ctx, typ)

	fields := flattenTo(structDescriptor.Fields, ctx.frozenConfig)

	orderedBindings = resolveBindings(fields)

	if len(orderedBindings) == 0 {
		return &emptyStructEncoder{}
	}

	finalOrderedFields := []structFieldTo{}
	for _, bindingTo := range orderedBindings {
		finalOrderedFields = append(finalOrderedFields, structFieldTo{
			encoder: bindingTo.binding.Encoder.(*structFieldEncoder),
			toName:  bindingTo.name,
		})
	}

	return &structEncoder{typ, finalOrderedFields}
}

func flattenTo(bindings []*Binding, cfg *frozenConfig) []*binding {
	flattened := make([]*binding, 0, len(bindings))

	for _, b := range bindings {
		for _, toName := range b.ToNames {
			flattened = append(flattened, &binding{
				binding: b,
				name:    toName,
				hasTag:  hasTag(b, cfg),
			})
		}
	}

	return flattened
}

func hasTag(b *Binding, cfg *frozenConfig) bool {
	before, _, _ := strings.Cut(b.Field.Tag().Get(cfg.getTagKey()), ",")
	return before != ""
}

func resolveBindings(fields []*binding) []*binding {
	sort.SliceStable(fields, func(i, j int) bool {
		// As per std's encoding/json,
		// it sorts fields by names, index depth(here we call it levels) and tags.
		// We've already sorted fields by index order in describeStruct.
		// By using stable sorting, we avoid sorting them again.
		if fields[i].name != fields[j].name {
			return fields[i].name < fields[j].name
		}
		if len(fields[i].binding.levels) != len(fields[j].binding.levels) {
			return len(fields[i].binding.levels) < len(fields[j].binding.levels)
		}
		if fields[i].hasTag != fields[j].hasTag {
			return fields[i].hasTag
		}
		return true // equal.
	})

	orderedBindings := trimOverlappingBindings(fields)

	sort.Slice(orderedBindings, func(i, j int) bool {
		left := orderedBindings[i].binding.levels
		right := orderedBindings[j].binding.levels
		k := 0
		for {
			if left[k] < right[k] {
				return true
			} else if left[k] > right[k] {
				return false
			}
			k++
		}
	})

	return orderedBindings
}

func trimOverlappingBindings(bindings []*binding) []*binding {
	out := bindings[:0]
	for nameRange, i := 0, 0; i < len(bindings); i += nameRange {
		for nameRange = 1; i+nameRange < len(bindings); nameRange++ {
			endOfRange := bindings[i+nameRange]
			if endOfRange.name != bindings[i].name {
				break
			}
		}
		if nameRange == 1 { // only one field for that name
			out = append(out, bindings[i])
		} else {
			fields := bindings[i : i+nameRange]
			if len(fields[0].binding.levels) == len(fields[1].binding.levels) &&
				fields[0].hasTag == fields[1].hasTag {
				continue
			}
			out = append(out, fields[0])
		}
	}

	return out
}

func createCheckIsEmpty(ctx *ctx, typ reflect2.Type) checkIsEmpty {
	encoder := createEncoderOfNative(ctx, typ)
	if encoder != nil {
		return encoder
	}
	kind := typ.Kind()
	switch kind {
	case reflect.Interface:
		return &dynamicEncoder{typ}
	case reflect.Struct:
		return &structEncoder{typ: typ}
	case reflect.Array:
		return &arrayEncoder{}
	case reflect.Slice:
		return &sliceEncoder{}
	case reflect.Map:
		return encoderOfMap(ctx, typ)
	case reflect.Ptr:
		return &OptionalEncoder{}
	default:
		return &lazyErrorEncoder{err: fmt.Errorf("unsupported type: %v", typ)}
	}
}

func resolveConflictBinding(cfg *frozenConfig, old, new *Binding) (ignoreOld, ignoreNew bool) {
	newTagged := new.Field.Tag().Get(cfg.getTagKey()) != ""
	oldTagged := old.Field.Tag().Get(cfg.getTagKey()) != ""
	if newTagged {
		if oldTagged {
			if len(old.levels) > len(new.levels) {
				return true, false
			} else if len(new.levels) > len(old.levels) {
				return false, true
			} else {
				return true, true
			}
		} else {
			return true, false
		}
	} else {
		if oldTagged {
			return true, false
		}
		if len(old.levels) > len(new.levels) {
			return true, false
		} else if len(new.levels) > len(old.levels) {
			return false, true
		} else {
			return true, true
		}
	}
}

type structFieldEncoder struct {
	field        reflect2.StructField
	fieldEncoder ValEncoder
	omitempty    bool
}

func (encoder *structFieldEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	fieldPtr := encoder.field.UnsafeGet(ptr)
	encoder.fieldEncoder.Encode(fieldPtr, stream)
	if stream.Error != nil && stream.Error != io.EOF {
		stream.Error = fmt.Errorf("%s: %s", encoder.field.Name(), stream.Error.Error())
	}
}

func (encoder *structFieldEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	fieldPtr := encoder.field.UnsafeGet(ptr)
	return encoder.fieldEncoder.IsEmpty(fieldPtr)
}

func (encoder *structFieldEncoder) IsEmbeddedPtrNil(ptr unsafe.Pointer) bool {
	isEmbeddedPtrNil, converted := encoder.fieldEncoder.(IsEmbeddedPtrNil)
	if !converted {
		return false
	}
	fieldPtr := encoder.field.UnsafeGet(ptr)
	return isEmbeddedPtrNil.IsEmbeddedPtrNil(fieldPtr)
}

type IsEmbeddedPtrNil interface {
	IsEmbeddedPtrNil(ptr unsafe.Pointer) bool
}

type structEncoder struct {
	typ    reflect2.Type
	fields []structFieldTo
}

type structFieldTo struct {
	encoder *structFieldEncoder
	toName  string
}

func (encoder *structEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteObjectStart()
	isNotFirst := false
	for _, field := range encoder.fields {
		if field.encoder.omitempty && field.encoder.IsEmpty(ptr) {
			continue
		}
		if field.encoder.IsEmbeddedPtrNil(ptr) {
			continue
		}
		if isNotFirst {
			stream.WriteMore()
		}
		stream.WriteObjectField(field.toName)
		field.encoder.Encode(ptr, stream)
		isNotFirst = true
	}
	stream.WriteObjectEnd()
	if stream.Error != nil && stream.Error != io.EOF {
		stream.Error = fmt.Errorf("%v.%s", encoder.typ, stream.Error.Error())
	}
}

func (encoder *structEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

type emptyStructEncoder struct {
}

func (encoder *emptyStructEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.WriteEmptyObject()
}

func (encoder *emptyStructEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

type stringModeNumberEncoder struct {
	elemEncoder ValEncoder
}

func (encoder *stringModeNumberEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	stream.writeByte('"')
	encoder.elemEncoder.Encode(ptr, stream)
	stream.writeByte('"')
}

func (encoder *stringModeNumberEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.elemEncoder.IsEmpty(ptr)
}

type stringModeStringEncoder struct {
	elemEncoder ValEncoder
	cfg         *frozenConfig
}

func (encoder *stringModeStringEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	tempStream := encoder.cfg.BorrowStream(nil)
	tempStream.Attachment = stream.Attachment
	defer encoder.cfg.ReturnStream(tempStream)
	encoder.elemEncoder.Encode(ptr, tempStream)
	stream.WriteString(string(tempStream.Buffer()))
}

func (encoder *stringModeStringEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.elemEncoder.IsEmpty(ptr)
}
