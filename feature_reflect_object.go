package jsoniter

import (
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

func encoderOfStruct(cfg *frozenConfig, typ reflect.Type) (ValEncoder, error) {
	type bindingTo struct {
		binding *Binding
		toName  string
		ignored bool
	}
	orderedBindings := []*bindingTo{}
	structDescriptor, err := describeStruct(cfg, typ)
	if err != nil {
		return nil, err
	}
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
				old.ignored, new.ignored = resolveConflictBinding(old.binding, new.binding)
			}
			orderedBindings = append(orderedBindings, new)
		}
	}
	if len(orderedBindings) == 0 {
		return &emptyStructEncoder{}, nil
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
	return &structEncoder{structDescriptor.onePtrEmbedded, structDescriptor.onePtrOptimization, finalOrderedFields}, nil
}

func resolveConflictBinding(old, new *Binding) (ignoreOld, ignoreNew bool) {
	newTagged := new.Field.Tag.Get("json") != ""
	oldTagged := old.Field.Tag.Get("json") != ""
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
		} else {
			if len(old.levels) > len(new.levels) {
				return true, false
			} else if len(new.levels) > len(old.levels) {
				return false, true
			} else {
				return true, true
			}
		}
	}
}

func decoderOfStruct(cfg *frozenConfig, typ reflect.Type) (ValDecoder, error) {
	bindings := map[string]*Binding{}
	structDescriptor, err := describeStruct(cfg, typ)
	if err != nil {
		return nil, err
	}
	for _, binding := range structDescriptor.Fields {
		for _, fromName := range binding.FromNames {
			old := bindings[fromName]
			if old == nil {
				bindings[fromName] = binding
				continue
			}
			ignoreOld, ignoreNew := resolveConflictBinding(old, binding)
			if ignoreOld {
				delete(bindings, fromName)
			}
			if !ignoreNew {
				bindings[fromName] = binding
			}
		}
	}
	fields := map[string]*structFieldDecoder{}
	for k, binding := range bindings {
		fields[k] = binding.Decoder.(*structFieldDecoder)
	}
	return createStructDecoder(typ, fields)
}

func createStructDecoder(typ reflect.Type, fields map[string]*structFieldDecoder) (ValDecoder, error) {
	knownHash := map[int32]struct{}{
		0: {},
	}
	switch len(fields) {
	case 0:
		return &skipDecoder{typ}, nil
	case 1:
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName)
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields}, nil
			} else {
				knownHash[fieldHash] = struct{}{}
			}
			return &oneFieldStructDecoder{typ, fieldHash, fieldDecoder}, nil
		}
	case 2:
		var fieldHash1 int32
		var fieldHash2 int32
		var fieldDecoder1 *structFieldDecoder
		var fieldDecoder2 *structFieldDecoder
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName)
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields}, nil
			} else {
				knownHash[fieldHash] = struct{}{}
			}
			if fieldHash1 == 0 {
				fieldHash1 = fieldHash
				fieldDecoder1 = fieldDecoder
			} else {
				fieldHash2 = fieldHash
				fieldDecoder2 = fieldDecoder
			}
		}
		return &twoFieldsStructDecoder{typ, fieldHash1, fieldDecoder1, fieldHash2, fieldDecoder2}, nil
	case 3:
		var fieldName1 int32
		var fieldName2 int32
		var fieldName3 int32
		var fieldDecoder1 *structFieldDecoder
		var fieldDecoder2 *structFieldDecoder
		var fieldDecoder3 *structFieldDecoder
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName)
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields}, nil
			} else {
				knownHash[fieldHash] = struct{}{}
			}
			if fieldName1 == 0 {
				fieldName1 = fieldHash
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == 0 {
				fieldName2 = fieldHash
				fieldDecoder2 = fieldDecoder
			} else {
				fieldName3 = fieldHash
				fieldDecoder3 = fieldDecoder
			}
		}
		return &threeFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1, fieldName2, fieldDecoder2, fieldName3, fieldDecoder3}, nil
	case 4:
		var fieldName1 int32
		var fieldName2 int32
		var fieldName3 int32
		var fieldName4 int32
		var fieldDecoder1 *structFieldDecoder
		var fieldDecoder2 *structFieldDecoder
		var fieldDecoder3 *structFieldDecoder
		var fieldDecoder4 *structFieldDecoder
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName)
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields}, nil
			} else {
				knownHash[fieldHash] = struct{}{}
			}
			if fieldName1 == 0 {
				fieldName1 = fieldHash
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == 0 {
				fieldName2 = fieldHash
				fieldDecoder2 = fieldDecoder
			} else if fieldName3 == 0 {
				fieldName3 = fieldHash
				fieldDecoder3 = fieldDecoder
			} else {
				fieldName4 = fieldHash
				fieldDecoder4 = fieldDecoder
			}
		}
		return &fourFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1, fieldName2, fieldDecoder2, fieldName3, fieldDecoder3,
			fieldName4, fieldDecoder4}, nil
	case 5:
		var fieldName1 int32
		var fieldName2 int32
		var fieldName3 int32
		var fieldName4 int32
		var fieldName5 int32
		var fieldDecoder1 *structFieldDecoder
		var fieldDecoder2 *structFieldDecoder
		var fieldDecoder3 *structFieldDecoder
		var fieldDecoder4 *structFieldDecoder
		var fieldDecoder5 *structFieldDecoder
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName)
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields}, nil
			} else {
				knownHash[fieldHash] = struct{}{}
			}
			if fieldName1 == 0 {
				fieldName1 = fieldHash
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == 0 {
				fieldName2 = fieldHash
				fieldDecoder2 = fieldDecoder
			} else if fieldName3 == 0 {
				fieldName3 = fieldHash
				fieldDecoder3 = fieldDecoder
			} else if fieldName4 == 0 {
				fieldName4 = fieldHash
				fieldDecoder4 = fieldDecoder
			} else {
				fieldName5 = fieldHash
				fieldDecoder5 = fieldDecoder
			}
		}
		return &fiveFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1, fieldName2, fieldDecoder2, fieldName3, fieldDecoder3,
			fieldName4, fieldDecoder4, fieldName5, fieldDecoder5}, nil
	case 6:
		var fieldName1 int32
		var fieldName2 int32
		var fieldName3 int32
		var fieldName4 int32
		var fieldName5 int32
		var fieldName6 int32
		var fieldDecoder1 *structFieldDecoder
		var fieldDecoder2 *structFieldDecoder
		var fieldDecoder3 *structFieldDecoder
		var fieldDecoder4 *structFieldDecoder
		var fieldDecoder5 *structFieldDecoder
		var fieldDecoder6 *structFieldDecoder
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName)
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields}, nil
			} else {
				knownHash[fieldHash] = struct{}{}
			}
			if fieldName1 == 0 {
				fieldName1 = fieldHash
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == 0 {
				fieldName2 = fieldHash
				fieldDecoder2 = fieldDecoder
			} else if fieldName3 == 0 {
				fieldName3 = fieldHash
				fieldDecoder3 = fieldDecoder
			} else if fieldName4 == 0 {
				fieldName4 = fieldHash
				fieldDecoder4 = fieldDecoder
			} else if fieldName5 == 0 {
				fieldName5 = fieldHash
				fieldDecoder5 = fieldDecoder
			} else {
				fieldName6 = fieldHash
				fieldDecoder6 = fieldDecoder
			}
		}
		return &sixFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1, fieldName2, fieldDecoder2, fieldName3, fieldDecoder3,
			fieldName4, fieldDecoder4, fieldName5, fieldDecoder5, fieldName6, fieldDecoder6}, nil
	case 7:
		var fieldName1 int32
		var fieldName2 int32
		var fieldName3 int32
		var fieldName4 int32
		var fieldName5 int32
		var fieldName6 int32
		var fieldName7 int32
		var fieldDecoder1 *structFieldDecoder
		var fieldDecoder2 *structFieldDecoder
		var fieldDecoder3 *structFieldDecoder
		var fieldDecoder4 *structFieldDecoder
		var fieldDecoder5 *structFieldDecoder
		var fieldDecoder6 *structFieldDecoder
		var fieldDecoder7 *structFieldDecoder
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName)
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields}, nil
			} else {
				knownHash[fieldHash] = struct{}{}
			}
			if fieldName1 == 0 {
				fieldName1 = fieldHash
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == 0 {
				fieldName2 = fieldHash
				fieldDecoder2 = fieldDecoder
			} else if fieldName3 == 0 {
				fieldName3 = fieldHash
				fieldDecoder3 = fieldDecoder
			} else if fieldName4 == 0 {
				fieldName4 = fieldHash
				fieldDecoder4 = fieldDecoder
			} else if fieldName5 == 0 {
				fieldName5 = fieldHash
				fieldDecoder5 = fieldDecoder
			} else if fieldName6 == 0 {
				fieldName6 = fieldHash
				fieldDecoder6 = fieldDecoder
			} else {
				fieldName7 = fieldHash
				fieldDecoder7 = fieldDecoder
			}
		}
		return &sevenFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1, fieldName2, fieldDecoder2, fieldName3, fieldDecoder3,
			fieldName4, fieldDecoder4, fieldName5, fieldDecoder5, fieldName6, fieldDecoder6,
			fieldName7, fieldDecoder7}, nil
	case 8:
		var fieldName1 int32
		var fieldName2 int32
		var fieldName3 int32
		var fieldName4 int32
		var fieldName5 int32
		var fieldName6 int32
		var fieldName7 int32
		var fieldName8 int32
		var fieldDecoder1 *structFieldDecoder
		var fieldDecoder2 *structFieldDecoder
		var fieldDecoder3 *structFieldDecoder
		var fieldDecoder4 *structFieldDecoder
		var fieldDecoder5 *structFieldDecoder
		var fieldDecoder6 *structFieldDecoder
		var fieldDecoder7 *structFieldDecoder
		var fieldDecoder8 *structFieldDecoder
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName)
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields}, nil
			} else {
				knownHash[fieldHash] = struct{}{}
			}
			if fieldName1 == 0 {
				fieldName1 = fieldHash
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == 0 {
				fieldName2 = fieldHash
				fieldDecoder2 = fieldDecoder
			} else if fieldName3 == 0 {
				fieldName3 = fieldHash
				fieldDecoder3 = fieldDecoder
			} else if fieldName4 == 0 {
				fieldName4 = fieldHash
				fieldDecoder4 = fieldDecoder
			} else if fieldName5 == 0 {
				fieldName5 = fieldHash
				fieldDecoder5 = fieldDecoder
			} else if fieldName6 == 0 {
				fieldName6 = fieldHash
				fieldDecoder6 = fieldDecoder
			} else if fieldName7 == 0 {
				fieldName7 = fieldHash
				fieldDecoder7 = fieldDecoder
			} else {
				fieldName8 = fieldHash
				fieldDecoder8 = fieldDecoder
			}
		}
		return &eightFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1, fieldName2, fieldDecoder2, fieldName3, fieldDecoder3,
			fieldName4, fieldDecoder4, fieldName5, fieldDecoder5, fieldName6, fieldDecoder6,
			fieldName7, fieldDecoder7, fieldName8, fieldDecoder8}, nil
	case 9:
		var fieldName1 int32
		var fieldName2 int32
		var fieldName3 int32
		var fieldName4 int32
		var fieldName5 int32
		var fieldName6 int32
		var fieldName7 int32
		var fieldName8 int32
		var fieldName9 int32
		var fieldDecoder1 *structFieldDecoder
		var fieldDecoder2 *structFieldDecoder
		var fieldDecoder3 *structFieldDecoder
		var fieldDecoder4 *structFieldDecoder
		var fieldDecoder5 *structFieldDecoder
		var fieldDecoder6 *structFieldDecoder
		var fieldDecoder7 *structFieldDecoder
		var fieldDecoder8 *structFieldDecoder
		var fieldDecoder9 *structFieldDecoder
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName)
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields}, nil
			} else {
				knownHash[fieldHash] = struct{}{}
			}
			if fieldName1 == 0 {
				fieldName1 = fieldHash
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == 0 {
				fieldName2 = fieldHash
				fieldDecoder2 = fieldDecoder
			} else if fieldName3 == 0 {
				fieldName3 = fieldHash
				fieldDecoder3 = fieldDecoder
			} else if fieldName4 == 0 {
				fieldName4 = fieldHash
				fieldDecoder4 = fieldDecoder
			} else if fieldName5 == 0 {
				fieldName5 = fieldHash
				fieldDecoder5 = fieldDecoder
			} else if fieldName6 == 0 {
				fieldName6 = fieldHash
				fieldDecoder6 = fieldDecoder
			} else if fieldName7 == 0 {
				fieldName7 = fieldHash
				fieldDecoder7 = fieldDecoder
			} else if fieldName8 == 0 {
				fieldName8 = fieldHash
				fieldDecoder8 = fieldDecoder
			} else {
				fieldName9 = fieldHash
				fieldDecoder9 = fieldDecoder
			}
		}
		return &nineFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1, fieldName2, fieldDecoder2, fieldName3, fieldDecoder3,
			fieldName4, fieldDecoder4, fieldName5, fieldDecoder5, fieldName6, fieldDecoder6,
			fieldName7, fieldDecoder7, fieldName8, fieldDecoder8, fieldName9, fieldDecoder9}, nil
	case 10:
		var fieldName1 int32
		var fieldName2 int32
		var fieldName3 int32
		var fieldName4 int32
		var fieldName5 int32
		var fieldName6 int32
		var fieldName7 int32
		var fieldName8 int32
		var fieldName9 int32
		var fieldName10 int32
		var fieldDecoder1 *structFieldDecoder
		var fieldDecoder2 *structFieldDecoder
		var fieldDecoder3 *structFieldDecoder
		var fieldDecoder4 *structFieldDecoder
		var fieldDecoder5 *structFieldDecoder
		var fieldDecoder6 *structFieldDecoder
		var fieldDecoder7 *structFieldDecoder
		var fieldDecoder8 *structFieldDecoder
		var fieldDecoder9 *structFieldDecoder
		var fieldDecoder10 *structFieldDecoder
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName)
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields}, nil
			} else {
				knownHash[fieldHash] = struct{}{}
			}
			if fieldName1 == 0 {
				fieldName1 = fieldHash
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == 0 {
				fieldName2 = fieldHash
				fieldDecoder2 = fieldDecoder
			} else if fieldName3 == 0 {
				fieldName3 = fieldHash
				fieldDecoder3 = fieldDecoder
			} else if fieldName4 == 0 {
				fieldName4 = fieldHash
				fieldDecoder4 = fieldDecoder
			} else if fieldName5 == 0 {
				fieldName5 = fieldHash
				fieldDecoder5 = fieldDecoder
			} else if fieldName6 == 0 {
				fieldName6 = fieldHash
				fieldDecoder6 = fieldDecoder
			} else if fieldName7 == 0 {
				fieldName7 = fieldHash
				fieldDecoder7 = fieldDecoder
			} else if fieldName8 == 0 {
				fieldName8 = fieldHash
				fieldDecoder8 = fieldDecoder
			} else if fieldName9 == 0 {
				fieldName9 = fieldHash
				fieldDecoder9 = fieldDecoder
			} else {
				fieldName10 = fieldHash
				fieldDecoder10 = fieldDecoder
			}
		}
		return &tenFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1, fieldName2, fieldDecoder2, fieldName3, fieldDecoder3,
			fieldName4, fieldDecoder4, fieldName5, fieldDecoder5, fieldName6, fieldDecoder6,
			fieldName7, fieldDecoder7, fieldName8, fieldDecoder8, fieldName9, fieldDecoder9,
			fieldName10, fieldDecoder10}, nil
	}
	return &generalStructDecoder{typ, fields}, nil
}

type generalStructDecoder struct {
	typ    reflect.Type
	fields map[string]*structFieldDecoder
}

func (decoder *generalStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	fieldBytes := iter.readObjectFieldAsBytes()
	field := *(*string)(unsafe.Pointer(&fieldBytes))
	fieldDecoder := decoder.fields[field]
	if fieldDecoder == nil {
		iter.Skip()
	} else {
		fieldDecoder.Decode(ptr, iter)
	}
	for iter.nextToken() == ',' {
		fieldBytes = iter.readObjectFieldAsBytes()
		field = *(*string)(unsafe.Pointer(&fieldBytes))
		fieldDecoder = decoder.fields[field]
		if fieldDecoder == nil {
			iter.Skip()
		} else {
			fieldDecoder.Decode(ptr, iter)
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type skipDecoder struct {
	typ reflect.Type
}

func (decoder *skipDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	iter.Skip()
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type oneFieldStructDecoder struct {
	typ          reflect.Type
	fieldHash    int32
	fieldDecoder *structFieldDecoder
}

func (decoder *oneFieldStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	for {
		if iter.readFieldHash() == decoder.fieldHash {
			decoder.fieldDecoder.Decode(ptr, iter)
		} else {
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type twoFieldsStructDecoder struct {
	typ           reflect.Type
	fieldHash1    int32
	fieldDecoder1 *structFieldDecoder
	fieldHash2    int32
	fieldDecoder2 *structFieldDecoder
}

func (decoder *twoFieldsStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	for {
		switch iter.readFieldHash() {
		case decoder.fieldHash1:
			decoder.fieldDecoder1.Decode(ptr, iter)
		case decoder.fieldHash2:
			decoder.fieldDecoder2.Decode(ptr, iter)
		default:
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type threeFieldsStructDecoder struct {
	typ           reflect.Type
	fieldHash1    int32
	fieldDecoder1 *structFieldDecoder
	fieldHash2    int32
	fieldDecoder2 *structFieldDecoder
	fieldHash3    int32
	fieldDecoder3 *structFieldDecoder
}

func (decoder *threeFieldsStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	for {
		switch iter.readFieldHash() {
		case decoder.fieldHash1:
			decoder.fieldDecoder1.Decode(ptr, iter)
		case decoder.fieldHash2:
			decoder.fieldDecoder2.Decode(ptr, iter)
		case decoder.fieldHash3:
			decoder.fieldDecoder3.Decode(ptr, iter)
		default:
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type fourFieldsStructDecoder struct {
	typ           reflect.Type
	fieldHash1    int32
	fieldDecoder1 *structFieldDecoder
	fieldHash2    int32
	fieldDecoder2 *structFieldDecoder
	fieldHash3    int32
	fieldDecoder3 *structFieldDecoder
	fieldHash4    int32
	fieldDecoder4 *structFieldDecoder
}

func (decoder *fourFieldsStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	for {
		switch iter.readFieldHash() {
		case decoder.fieldHash1:
			decoder.fieldDecoder1.Decode(ptr, iter)
		case decoder.fieldHash2:
			decoder.fieldDecoder2.Decode(ptr, iter)
		case decoder.fieldHash3:
			decoder.fieldDecoder3.Decode(ptr, iter)
		case decoder.fieldHash4:
			decoder.fieldDecoder4.Decode(ptr, iter)
		default:
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type fiveFieldsStructDecoder struct {
	typ           reflect.Type
	fieldHash1    int32
	fieldDecoder1 *structFieldDecoder
	fieldHash2    int32
	fieldDecoder2 *structFieldDecoder
	fieldHash3    int32
	fieldDecoder3 *structFieldDecoder
	fieldHash4    int32
	fieldDecoder4 *structFieldDecoder
	fieldHash5    int32
	fieldDecoder5 *structFieldDecoder
}

func (decoder *fiveFieldsStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	for {
		switch iter.readFieldHash() {
		case decoder.fieldHash1:
			decoder.fieldDecoder1.Decode(ptr, iter)
		case decoder.fieldHash2:
			decoder.fieldDecoder2.Decode(ptr, iter)
		case decoder.fieldHash3:
			decoder.fieldDecoder3.Decode(ptr, iter)
		case decoder.fieldHash4:
			decoder.fieldDecoder4.Decode(ptr, iter)
		case decoder.fieldHash5:
			decoder.fieldDecoder5.Decode(ptr, iter)
		default:
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type sixFieldsStructDecoder struct {
	typ           reflect.Type
	fieldHash1    int32
	fieldDecoder1 *structFieldDecoder
	fieldHash2    int32
	fieldDecoder2 *structFieldDecoder
	fieldHash3    int32
	fieldDecoder3 *structFieldDecoder
	fieldHash4    int32
	fieldDecoder4 *structFieldDecoder
	fieldHash5    int32
	fieldDecoder5 *structFieldDecoder
	fieldHash6    int32
	fieldDecoder6 *structFieldDecoder
}

func (decoder *sixFieldsStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	for {
		switch iter.readFieldHash() {
		case decoder.fieldHash1:
			decoder.fieldDecoder1.Decode(ptr, iter)
		case decoder.fieldHash2:
			decoder.fieldDecoder2.Decode(ptr, iter)
		case decoder.fieldHash3:
			decoder.fieldDecoder3.Decode(ptr, iter)
		case decoder.fieldHash4:
			decoder.fieldDecoder4.Decode(ptr, iter)
		case decoder.fieldHash5:
			decoder.fieldDecoder5.Decode(ptr, iter)
		case decoder.fieldHash6:
			decoder.fieldDecoder6.Decode(ptr, iter)
		default:
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type sevenFieldsStructDecoder struct {
	typ           reflect.Type
	fieldHash1    int32
	fieldDecoder1 *structFieldDecoder
	fieldHash2    int32
	fieldDecoder2 *structFieldDecoder
	fieldHash3    int32
	fieldDecoder3 *structFieldDecoder
	fieldHash4    int32
	fieldDecoder4 *structFieldDecoder
	fieldHash5    int32
	fieldDecoder5 *structFieldDecoder
	fieldHash6    int32
	fieldDecoder6 *structFieldDecoder
	fieldHash7    int32
	fieldDecoder7 *structFieldDecoder
}

func (decoder *sevenFieldsStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	for {
		switch iter.readFieldHash() {
		case decoder.fieldHash1:
			decoder.fieldDecoder1.Decode(ptr, iter)
		case decoder.fieldHash2:
			decoder.fieldDecoder2.Decode(ptr, iter)
		case decoder.fieldHash3:
			decoder.fieldDecoder3.Decode(ptr, iter)
		case decoder.fieldHash4:
			decoder.fieldDecoder4.Decode(ptr, iter)
		case decoder.fieldHash5:
			decoder.fieldDecoder5.Decode(ptr, iter)
		case decoder.fieldHash6:
			decoder.fieldDecoder6.Decode(ptr, iter)
		case decoder.fieldHash7:
			decoder.fieldDecoder7.Decode(ptr, iter)
		default:
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type eightFieldsStructDecoder struct {
	typ           reflect.Type
	fieldHash1    int32
	fieldDecoder1 *structFieldDecoder
	fieldHash2    int32
	fieldDecoder2 *structFieldDecoder
	fieldHash3    int32
	fieldDecoder3 *structFieldDecoder
	fieldHash4    int32
	fieldDecoder4 *structFieldDecoder
	fieldHash5    int32
	fieldDecoder5 *structFieldDecoder
	fieldHash6    int32
	fieldDecoder6 *structFieldDecoder
	fieldHash7    int32
	fieldDecoder7 *structFieldDecoder
	fieldHash8    int32
	fieldDecoder8 *structFieldDecoder
}

func (decoder *eightFieldsStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	for {
		switch iter.readFieldHash() {
		case decoder.fieldHash1:
			decoder.fieldDecoder1.Decode(ptr, iter)
		case decoder.fieldHash2:
			decoder.fieldDecoder2.Decode(ptr, iter)
		case decoder.fieldHash3:
			decoder.fieldDecoder3.Decode(ptr, iter)
		case decoder.fieldHash4:
			decoder.fieldDecoder4.Decode(ptr, iter)
		case decoder.fieldHash5:
			decoder.fieldDecoder5.Decode(ptr, iter)
		case decoder.fieldHash6:
			decoder.fieldDecoder6.Decode(ptr, iter)
		case decoder.fieldHash7:
			decoder.fieldDecoder7.Decode(ptr, iter)
		case decoder.fieldHash8:
			decoder.fieldDecoder8.Decode(ptr, iter)
		default:
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type nineFieldsStructDecoder struct {
	typ           reflect.Type
	fieldHash1    int32
	fieldDecoder1 *structFieldDecoder
	fieldHash2    int32
	fieldDecoder2 *structFieldDecoder
	fieldHash3    int32
	fieldDecoder3 *structFieldDecoder
	fieldHash4    int32
	fieldDecoder4 *structFieldDecoder
	fieldHash5    int32
	fieldDecoder5 *structFieldDecoder
	fieldHash6    int32
	fieldDecoder6 *structFieldDecoder
	fieldHash7    int32
	fieldDecoder7 *structFieldDecoder
	fieldHash8    int32
	fieldDecoder8 *structFieldDecoder
	fieldHash9    int32
	fieldDecoder9 *structFieldDecoder
}

func (decoder *nineFieldsStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	for {
		switch iter.readFieldHash() {
		case decoder.fieldHash1:
			decoder.fieldDecoder1.Decode(ptr, iter)
		case decoder.fieldHash2:
			decoder.fieldDecoder2.Decode(ptr, iter)
		case decoder.fieldHash3:
			decoder.fieldDecoder3.Decode(ptr, iter)
		case decoder.fieldHash4:
			decoder.fieldDecoder4.Decode(ptr, iter)
		case decoder.fieldHash5:
			decoder.fieldDecoder5.Decode(ptr, iter)
		case decoder.fieldHash6:
			decoder.fieldDecoder6.Decode(ptr, iter)
		case decoder.fieldHash7:
			decoder.fieldDecoder7.Decode(ptr, iter)
		case decoder.fieldHash8:
			decoder.fieldDecoder8.Decode(ptr, iter)
		case decoder.fieldHash9:
			decoder.fieldDecoder9.Decode(ptr, iter)
		default:
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type tenFieldsStructDecoder struct {
	typ            reflect.Type
	fieldHash1     int32
	fieldDecoder1  *structFieldDecoder
	fieldHash2     int32
	fieldDecoder2  *structFieldDecoder
	fieldHash3     int32
	fieldDecoder3  *structFieldDecoder
	fieldHash4     int32
	fieldDecoder4  *structFieldDecoder
	fieldHash5     int32
	fieldDecoder5  *structFieldDecoder
	fieldHash6     int32
	fieldDecoder6  *structFieldDecoder
	fieldHash7     int32
	fieldDecoder7  *structFieldDecoder
	fieldHash8     int32
	fieldDecoder8  *structFieldDecoder
	fieldHash9     int32
	fieldDecoder9  *structFieldDecoder
	fieldHash10    int32
	fieldDecoder10 *structFieldDecoder
}

func (decoder *tenFieldsStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	for {
		switch iter.readFieldHash() {
		case decoder.fieldHash1:
			decoder.fieldDecoder1.Decode(ptr, iter)
		case decoder.fieldHash2:
			decoder.fieldDecoder2.Decode(ptr, iter)
		case decoder.fieldHash3:
			decoder.fieldDecoder3.Decode(ptr, iter)
		case decoder.fieldHash4:
			decoder.fieldDecoder4.Decode(ptr, iter)
		case decoder.fieldHash5:
			decoder.fieldDecoder5.Decode(ptr, iter)
		case decoder.fieldHash6:
			decoder.fieldDecoder6.Decode(ptr, iter)
		case decoder.fieldHash7:
			decoder.fieldDecoder7.Decode(ptr, iter)
		case decoder.fieldHash8:
			decoder.fieldDecoder8.Decode(ptr, iter)
		case decoder.fieldHash9:
			decoder.fieldDecoder9.Decode(ptr, iter)
		case decoder.fieldHash10:
			decoder.fieldDecoder10.Decode(ptr, iter)
		default:
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.typ, iter.Error.Error())
	}
}

type structFieldDecoder struct {
	field        *reflect.StructField
	fieldDecoder ValDecoder
}

func (decoder *structFieldDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	fieldPtr := uintptr(ptr) + decoder.field.Offset
	decoder.fieldDecoder.Decode(unsafe.Pointer(fieldPtr), iter)
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%s: %s", decoder.field.Name, iter.Error.Error())
	}
}

type structFieldEncoder struct {
	field        *reflect.StructField
	fieldEncoder ValEncoder
	omitempty    bool
}

func (encoder *structFieldEncoder) Encode(ptr unsafe.Pointer, stream *Stream) {
	fieldPtr := uintptr(ptr) + encoder.field.Offset
	encoder.fieldEncoder.Encode(unsafe.Pointer(fieldPtr), stream)
	if stream.Error != nil && stream.Error != io.EOF {
		stream.Error = fmt.Errorf("%s: %s", encoder.field.Name, stream.Error.Error())
	}
}

func (encoder *structFieldEncoder) EncodeInterface(val interface{}, stream *Stream) {
	WriteToStream(val, stream, encoder)
}

func (encoder *structFieldEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	fieldPtr := uintptr(ptr) + encoder.field.Offset
	return encoder.fieldEncoder.IsEmpty(unsafe.Pointer(fieldPtr))
}

type structEncoder struct {
	onePtrEmbedded     bool
	onePtrOptimization bool
	fields             []structFieldTo
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
		if isNotFirst {
			stream.WriteMore()
		}
		stream.WriteObjectField(field.toName)
		field.encoder.Encode(ptr, stream)
		isNotFirst = true
	}
	stream.WriteObjectEnd()
}

func (encoder *structEncoder) EncodeInterface(val interface{}, stream *Stream) {
	e := (*emptyInterface)(unsafe.Pointer(&val))
	if encoder.onePtrOptimization {
		if e.word == nil && encoder.onePtrEmbedded {
			stream.WriteObjectStart()
			stream.WriteObjectEnd()
			return
		}
		ptr := uintptr(e.word)
		e.word = unsafe.Pointer(&ptr)
	}
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		encoder.Encode(unsafe.Pointer(&e.word), stream)
	} else {
		encoder.Encode(e.word, stream)
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

func (encoder *emptyStructEncoder) EncodeInterface(val interface{}, stream *Stream) {
	WriteToStream(val, stream, encoder)
}

func (encoder *emptyStructEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}
