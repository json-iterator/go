package jsoniter

import (
	"fmt"
	"io"
	"reflect"
	"unsafe"

	"github.com/modern-go/reflect2"
)

// checkIsZero is an interface for checking if a value is zero (for omitzero support)
type checkIsZero interface {
	IsZero(ptr unsafe.Pointer) bool
}

// isZeroer is the interface implemented by types that can report whether their value is zero.
// This is compatible with Go 1.24's encoding/json omitzero feature.
type isZeroer interface {
	IsZero() bool
}

func encoderOfStruct(ctx *ctx, typ reflect2.Type) ValEncoder {
	type bindingTo struct {
		binding *Binding
		toName  string
		ignored bool
	}
	orderedBindings := []*bindingTo{}
	structDescriptor := describeStruct(ctx, typ)
	for _, binding := range structDescriptor.Fields {
		for _, toName := range binding.ToNames {
			new := &bindingTo{
				binding: binding,
				toName:  toName,
			}
			for _, old := range orderedBindings {
				if old.toName != toName {
					continue
				}
				old.ignored, new.ignored = resolveConflictBinding(ctx.frozenConfig, old.binding, new.binding)
			}
			orderedBindings = append(orderedBindings, new)
		}
	}
	if len(orderedBindings) == 0 {
		return &emptyStructEncoder{}
	}
	finalOrderedFields := []structFieldTo{}
	for _, bindingTo := range orderedBindings {
		if !bindingTo.ignored {
			finalOrderedFields = append(finalOrderedFields, structFieldTo{
				encoder: bindingTo.binding.Encoder.(*structFieldEncoder),
				toName:  bindingTo.toName,
			})
		}
	}
	return &structEncoder{typ, finalOrderedFields}
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

// isZeroerType is the reflect2.Type for the isZeroer interface
var isZeroerType = reflect2.TypeOfPtr((*isZeroer)(nil)).Elem()

// createCheckIsZero creates a zero checker for the given type.
func createCheckIsZero(ctx *ctx, typ reflect2.Type) checkIsZero {
	// Check if the type implements isZeroer interface
	if typ.Implements(isZeroerType) {
		return &directIsZeroerChecker{valType: typ}
	}
	// Check if pointer to type implements isZeroer interface
	ptrType := reflect2.PtrTo(typ)
	if ptrType.Implements(isZeroerType) {
		return &ptrIsZeroerChecker{valType: typ}
	}
	// Fall back to comparing with zero value
	return createZeroValueChecker(typ)
}

// directIsZeroerChecker checks IsZero() for types that directly implement the interface
type directIsZeroerChecker struct {
	valType reflect2.Type
}

func (checker *directIsZeroerChecker) IsZero(ptr unsafe.Pointer) bool {
	obj := checker.valType.UnsafeIndirect(ptr)
	if checker.valType.IsNullable() && reflect2.IsNil(obj) {
		return true
	}
	return obj.(isZeroer).IsZero()
}

// ptrIsZeroerChecker checks IsZero() for types where *T implements the interface
type ptrIsZeroerChecker struct {
	valType reflect2.Type
}

func (checker *ptrIsZeroerChecker) IsZero(ptr unsafe.Pointer) bool {
	obj := checker.valType.PackEFace(ptr)
	return obj.(isZeroer).IsZero()
}

// createZeroValueChecker creates a checker that compares against zero value
func createZeroValueChecker(typ reflect2.Type) checkIsZero {
	kind := typ.Kind()
	switch kind {
	case reflect.String:
		return &stringZeroChecker{}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &intZeroChecker{typ: typ}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return &uintZeroChecker{typ: typ}
	case reflect.Float32:
		return &float32ZeroChecker{}
	case reflect.Float64:
		return &float64ZeroChecker{}
	case reflect.Bool:
		return &boolZeroChecker{}
	case reflect.Interface, reflect.Ptr:
		return &nilZeroChecker{}
	case reflect.Slice, reflect.Map:
		return &nilZeroChecker{}
	case reflect.Array:
		return createArrayZeroChecker(typ)
	case reflect.Struct:
		return createStructZeroChecker(typ)
	default:
		return nil
	}
}

// stringZeroChecker checks if a string is empty
type stringZeroChecker struct{}

func (checker *stringZeroChecker) IsZero(ptr unsafe.Pointer) bool {
	return *((*string)(ptr)) == ""
}

// intZeroChecker checks if an int-like value is zero
type intZeroChecker struct {
	typ reflect2.Type
}

func (checker *intZeroChecker) IsZero(ptr unsafe.Pointer) bool {
	switch checker.typ.Kind() {
	case reflect.Int:
		return *((*int)(ptr)) == 0
	case reflect.Int8:
		return *((*int8)(ptr)) == 0
	case reflect.Int16:
		return *((*int16)(ptr)) == 0
	case reflect.Int32:
		return *((*int32)(ptr)) == 0
	case reflect.Int64:
		return *((*int64)(ptr)) == 0
	}
	return false
}

// uintZeroChecker checks if a uint-like value is zero
type uintZeroChecker struct {
	typ reflect2.Type
}

func (checker *uintZeroChecker) IsZero(ptr unsafe.Pointer) bool {
	switch checker.typ.Kind() {
	case reflect.Uint:
		return *((*uint)(ptr)) == 0
	case reflect.Uint8:
		return *((*uint8)(ptr)) == 0
	case reflect.Uint16:
		return *((*uint16)(ptr)) == 0
	case reflect.Uint32:
		return *((*uint32)(ptr)) == 0
	case reflect.Uint64:
		return *((*uint64)(ptr)) == 0
	case reflect.Uintptr:
		return *((*uintptr)(ptr)) == 0
	}
	return false
}

// float32ZeroChecker checks if a float32 is zero
type float32ZeroChecker struct{}

func (checker *float32ZeroChecker) IsZero(ptr unsafe.Pointer) bool {
	return *((*float32)(ptr)) == 0
}

// float64ZeroChecker checks if a float64 is zero
type float64ZeroChecker struct{}

func (checker *float64ZeroChecker) IsZero(ptr unsafe.Pointer) bool {
	return *((*float64)(ptr)) == 0
}

// boolZeroChecker checks if a bool is false
type boolZeroChecker struct{}

func (checker *boolZeroChecker) IsZero(ptr unsafe.Pointer) bool {
	return !*((*bool)(ptr))
}

// nilZeroChecker checks if a pointer/interface/slice/map is nil
type nilZeroChecker struct{}

func (checker *nilZeroChecker) IsZero(ptr unsafe.Pointer) bool {
	return *((*unsafe.Pointer)(ptr)) == nil
}

// arrayBytesZeroChecker checks if all bytes in the array are zero
type arrayBytesZeroChecker struct {
	size uintptr
}

func (checker *arrayBytesZeroChecker) IsZero(ptr unsafe.Pointer) bool {
	bytes := (*[1 << 30]byte)(ptr)[:checker.size]
	for _, b := range bytes {
		if b != 0 {
			return false
		}
	}
	return true
}

func createArrayZeroChecker(typ reflect2.Type) checkIsZero {
	arrayType := typ.(*reflect2.UnsafeArrayType)
	elemType := arrayType.Elem()
	// For simple types, we can use a simpler check
	if isSimpleZeroType(elemType) {
		return &arrayBytesZeroChecker{size: typ.Type1().Size()}
	}
	return nil // Complex arrays fall back to reflect-based comparison
}

func isSimpleZeroType(typ reflect2.Type) bool {
	switch typ.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

// structZeroChecker uses reflect to check if a struct is zero
type structZeroChecker struct {
	typ reflect2.Type
}

func createStructZeroChecker(typ reflect2.Type) checkIsZero {
	return &structZeroChecker{typ: typ}
}

func (checker *structZeroChecker) IsZero(ptr unsafe.Pointer) bool {
	// Use reflect to compare with zero value
	val := checker.typ.UnsafeIndirect(ptr)
	rv := reflect.ValueOf(val)
	return rv.IsZero()
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
	omitzero     bool
	zeroChecker  checkIsZero
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

// IsZero checks if the field value is zero using the IsZero() method if available,
// or by comparing with the zero value of the type
func (encoder *structFieldEncoder) IsZero(ptr unsafe.Pointer) bool {
	fieldPtr := encoder.field.UnsafeGet(ptr)
	// First check if field encoder implements IsZero
	if isZeroer, ok := encoder.fieldEncoder.(checkIsZero); ok {
		return isZeroer.IsZero(fieldPtr)
	}
	// Use the pre-computed zero checker
	if encoder.zeroChecker != nil {
		return encoder.zeroChecker.IsZero(fieldPtr)
	}
	return false
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
		if field.encoder.omitzero && field.encoder.IsZero(ptr) {
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
